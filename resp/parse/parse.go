package parse

import (
	"bufio"
	"errors"
	"gopro/interface/resp"
	"gopro/lib/logger"
	"gopro/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

// Payload 客户端指令的数据载体
type Payload struct {
	Data resp.Reply
	Err  error
}

// 指令解析的状态记录
type readState struct {
	readingMultiLine  bool     // 这条解析的指令是单行的还是多行的
	expectedArgsCount int      // 期望的参数的个数
	msgType           byte     //客户端的这条消息指令的类型
	args              [][]byte //客户端发送的具体数据内容, 用于处理*字符串数组
	bulkLen           int64    //客户端发送的字符串数据块的大小 比如$4\r\nPING\r\n
}

// 判断是不是解析结束了
func (s *readState) finished() bool {
	// 期望读取的个数和已经读进来的数据的len相等代表是已经读取完毕了
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount

}

// ParseStream 暴露给外部的，异步的解析，并发的解析；输出是一个收东西的channel,把解析的结果通过管道异步吐出去，而不阻塞服务端
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

// 读取客户端通过tcp传给服务端的字节流
func parse0(reader io.Reader, ch chan<- *Payload) {

	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()

	bufferReder := bufio.NewReader(reader)
	var msg []byte
	var err error
	var state readState

	for true {
		var ioErr bool
		// 读取内容
		msg, ioErr, err = readLine(bufferReder, &state)
		if err != nil {
			// io 错误
			if ioErr {
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			} else {
				ch <- &Payload{
					Err: err,
				}
				state = readState{}
				continue
			}
		}
		// 解析内容, 什么时间不是多行解析？
		// 1 刚开始初始化的时候，不是多行
		// 2 像+OK\r\n等等就是单行
		// 3 什么时间会变化为多行？ 就是比如读数据解析数据的时候会根据解析规则比如parseMultiBulkHeader 变为多行数据
		if !state.readingMultiLine {
			// 如果收到的数据起shi是*，代表后面就是数组
			if msg[0] == '*' { // *3\r\n
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("Protocol Error " + string(msg)),
					}
					state = readState{}
					continue
				}

				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: &reply.EmptyMultiBulkReply{},
					}
				}
				state = readState{}
				continue
			} else if msg[0] == '$' { // $4\r\n
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("Protocol Error " + string(msg)),
					}
					state = readState{}
					continue
				}

				// $-1\r\n
				if state.bulkLen == 0 {
					ch <- &Payload{
						Data: &reply.NullBulkReply{},
					}
				}
				// 重置state
				state = readState{}
				continue
			} else if msg[0] == '+' || msg[0] == '-' || msg[0] == ':' {
				var lineReply resp.Reply
				lineReply, err = parseSingleLineReply(msg, &state)
				ch <- &Payload{
					Data: lineReply,
					Err:  err,
				}
				state = readState{}
				continue
			} else { // 以上条件都不满足则执行下面的情况
				ch <- &Payload{
					Err: errors.New("Protocol Error " + string(msg)),
				}
				state = readState{}
				continue
			}
		} else {
			// 读取消息内容（前面是读取协议头）
			err = readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: errors.New("Protocol Error " + string(msg)),
				}
				state = readState{}
				continue
			}

			if state.finished() {
				switch state.msgType {

				case '*':
					ch <- &Payload{
						Data: reply.MakeMultiBulkReply(state.args),
						Err:  err,
					}
				case '$':
					ch <- &Payload{
						Data: reply.MakeBulkReply(state.args[0]),
						Err:  err,
					}

				}

				state = readState{}
				continue
			}

		}
	}

}

/**
ReadLine
返回这一行内容，
是否有io错误
具体的error
*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
*/
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {

	var msg []byte
	var err error
	// 没有字符串长度，只需要处理\r\n。 如果是$4\r\nPING\r\n  这种单字符串情况会塞长度，这时候严格按照长度来取值
	if state.bulkLen == 0 {
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol wrong " + string(msg))
		}

	} else { // 读取到了$之后的数字,有字符串长度，那么按照字节个数严格读取. 针对单字符串 PING\r\n
		msg = make([]byte, state.bulkLen+2)
		_, err := io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol wrong " + string(msg))
		}
		// 当把这些字节读取完毕时，需要把这个指针的要读取的字节长度置0
		state.bulkLen = 0
	}
	return msg, false, nil
}

// 解析客户端发送的单行正常结构、错误结构、整数结构. 比如+OK\r\n, -err\r\n, :156\r\n...
// 这种因为不像字符串或者字符串数组那样复杂，可以直接使用Reply返回结果
func parseSingleLineReply(msg []byte, state *readState) (resp.Reply, error) {

	suffix := strings.TrimSuffix(string(msg), "\r\n")
	if len(suffix) == 0 {
		return nil, errors.New("Protocol Errror " + string(msg))
	}

	var msgType byte = msg[0]
	switch msgType {
	case '+':
		return reply.MakeStatusReply(string(msg[1:])), nil
	case '-':
		return reply.MakeStandardErrReply(string(msg[1:])), nil
	case ':':
		parseInt, err := strconv.ParseInt(string(msg[1:]), 10, 64)
		if err != nil {
			return nil, err
		}
		return reply.MakeIntReply(parseInt), nil
	}
	return nil, nil
}

// 解析单字符串 $4\r\nPING\r\n. 考虑到通用，最终也是按照多行解析逻辑来
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	var bulkLen int64
	bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("Protocol Error " + string(msg))
	}
	if bulkLen <= 0 {
		state.bulkLen = 0
		state.expectedArgsCount = 0
		return nil
	} else if bulkLen > 0 {
		state.readingMultiLine = true
		state.msgType = msg[0]
		state.bulkLen = bulkLen
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("Protocol Error " + string(msg))
	}

}

// 解析多字符串
// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
func parseMultiBulkHeader(msg []byte, state *readState) error {

	var err error
	var expectedSize int64
	expectedSize, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)

	if err != nil {
		return errors.New("Protocol Error " + string(msg))
	}
	if expectedSize == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedSize > 0 {
		state.expectedArgsCount = int(expectedSize)
		state.msgType = msg[0]
		state.args = make([][]byte, 0, expectedSize)
		state.readingMultiLine = true
		return nil
	} else {
		return errors.New("Protocol Error " + string(msg))
	}
}

// 针对数组之后的body $3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
// 针对这种body PING\r\n  SET\r\n
func readBody(msg []byte, state *readState) error {

	// 读取一行数据,先用\r\n切除一行数据出来
	line := msg[0 : len(msg)-2]
	var err error

	// 说明这一行数据是数组的一部分，才会以$开头
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol err " + string(msg))
		}
		// $0\r\n
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
