package resp

// Connection 代表redis的连接
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
