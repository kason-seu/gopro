package reply

// PongReply redis客户端输入一个PING Server端返回一个PONG
type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (p *PongReply) ToBytes() []byte {

	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

type OkReply struct {
}

var okBytes = []byte("+OK\r\n")

func (o *OkReply) ToBytes() []byte {
	return okBytes
}

// 这样无需每次make的时候都去创建OkReply对象
var okReply = &OkReply{}

func MakeOkReply() *OkReply {
	return okReply
}

// NullBulkReply 服务端的空回复
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// EmptyMultiBulkReply 服务端空数组的回复
type EmptyMultiBulkReply struct {
}

var emptyMultiBulkReplyBytes = []byte("*0\r\n")

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkReplyBytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

type NoReply struct {
}

var noReplyBytes = []byte("")

func (n *NoReply) ToBytes() []byte {
	return noReplyBytes
}
func MakeNoReply() *NoReply {
	return &NoReply{}
}
