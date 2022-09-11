package reply

import (
	"bytes"
	"gopro/interface/resp"
	"strconv"
)

// Reply用于客户端的请求和服务端的回复，通用的。
var (
	nullBulkReply = "$-1"
	CRLF          = "\r\n"
)

// BulkReply  自定义字符串的Reply
type BulkReply struct {
	// 字符串的值，比如hello, 那么reply 就是$5hello\r\n，相当于自动转换这个过程
	Arg []byte
}

func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return []byte(nullBulkReply + CRLF)
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

// MakeBulkReply 外面传入一个字符串，自动拼装返回值
func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

// MultiBulkReply 多字符串的自定义封装
type MultiBulkReply struct {
	Args [][]byte
}

func (m *MultiBulkReply) ToBytes() []byte {
	argNum := len(m.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argNum) + CRLF)

	for _, arg := range m.Args {
		if arg == nil {
			buf.WriteString(nullBulkReply + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {

	return &MultiBulkReply{
		Args: args,
	}
}

// StatusReply 正常答复
type StatusReply struct {
	Status string
}

func (s *StatusReply) ToBytes() []byte {
	return []byte("+" + s.Status + CRLF)
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

// IntReply 整数答复
type IntReply struct {
	Code int64
}

func (i *IntReply) ToBytes() []byte {

	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)

}

func MakeIntReply(code int64) *IntReply {

	return &IntReply{
		Code: code,
	}

}

// ErrorReply 定义redis错误的回复接口
type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

// StandardErrReply 标准错误答复
type StandardErrReply struct {
	Status string
}

func (s *StandardErrReply) Error() string {
	return s.Status
}

func (s *StandardErrReply) ToBytes() []byte {

	return []byte("-" + s.Status + CRLF)
}

func MakeStandardErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

// IsErrReply 判断是否是错误回复
func IsErrReply(r resp.Reply) bool {
	return r.ToBytes()[0] == '-'
}
