package reply

// UnknownErrReply 未知错误
type UnknownErrReply struct {
}

var unKnownErrBytes = []byte("-Err unknown\r\n")

func (u *UnknownErrReply) Error() string {
	return "Err unknown"
}

func (u *UnknownErrReply) ToBytes() []byte {
	return unKnownErrBytes
}

func MakeUnkownErrReply() *UnknownErrReply {
	return &UnknownErrReply{}
}

// ArgNumErrReply 参数个数异常，比如SET KEY VALUE 需要三个参数，但是只传2个，则参数个数有问题
type ArgNumErrReply struct {
	Cmd string // 记录那个客户端的指令
}

func (a *ArgNumErrReply) Error() string {
	return "Err wrong number of arguments for '" + a.Cmd + "' command\r\n"
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-Err wrong number of arguments for '" + a.Cmd + "' command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

// SyntaxErrReply 语法错误
type SyntaxErrReply struct {
}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

func (s *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

var syntaxErrReply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return syntaxErrReply
}

// WrongTypeErrReply 数据类型错误
type WrongTypeErrReply struct {
}

var wrongTypeErrBytes = []byte("-Err wrong type operation against a key holding the wrong kind of value\r\n")

func (w *WrongTypeErrReply) Error() string {
	return "wrong type operation against a key holding the wrong kind of value"
}

func (w *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

var wrongTypeErrRepley = new(WrongTypeErrReply)

func MakeWrongTypeErrorReply() *WrongTypeErrReply {
	return wrongTypeErrRepley
}

// ProtocolErrReply 协议错误
type ProtocolErrReply struct {
	Msg string
}

func (p *ProtocolErrReply) Error() string {
	return "Err protocol error: " + p.Msg
}

func (p *ProtocolErrReply) ToBytes() []byte {
	return []byte("-Err protocol error: '" + p.Msg + "'\r\n")
}

func MakeProtocolErrReply(msg string) *ProtocolErrReply {

	return &ProtocolErrReply{
		Msg: msg,
	}

}
