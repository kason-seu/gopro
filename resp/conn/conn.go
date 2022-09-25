package conn

import (
	"gopro/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn         net.Conn
	waitingReply wait.Wait
	mu           sync.Mutex
	selectedDB   int
}

func (c *Connection) RemoteAddr() net.Addr {

	return c.conn.RemoteAddr()

}

func NewConn(conn net.Conn) *Connection {

	return &Connection{
		conn: conn,
	}

}

func (c *Connection) Close() error {

	c.waitingReply.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()

	return nil
}

func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	// 加锁，同一时刻只能有一协程 对客户端写数据。 如果有两个xiecheng都往同一个客户端写数据，解析回发就会出错
	c.mu.Lock()
	c.waitingReply.Add(1)
	defer func() {
		c.waitingReply.Done()
		c.mu.Unlock()
	}()
	_, err := c.conn.Write(bytes)
	return err
}

func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

func (c *Connection) SelectDB(dbNum int) {

	c.selectedDB = dbNum

}
