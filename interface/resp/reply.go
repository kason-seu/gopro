package resp

// Reply 代表服务端对客户端的回复
type Reply interface {
	// ToBytes 因为tcp协议里面的读写都是针对字节的，所以这里都要转为字节
	ToBytes() []byte
}
