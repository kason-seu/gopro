package tcp

import (
	"bufio"
	"context"
	"gopro/lib/logger"
	"gopro/lib/sync/atomic"
	"gopro/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	// 等10s在关闭连接
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

type EchoHandler struct {
	// map 当set用，用于存储client
	activeConn sync.Map
	// handler的状态
	closing atomic.Boolean
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {

	if handler.closing.Get() {
		_ = conn.Close()
		return
	}
	logger.Info("create EchoClient")
	client := &EchoClient{Conn: conn}
	// 存储客户端
	handler.activeConn.Store(client, struct{}{})
	// 创建一个缓存buffer，读取内容
	reader := bufio.NewReader(conn)

	// 循环处理客户端的请求
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				// 代表读取到了终点了
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		// 该客户端增加一个干活的请求
		client.Waiting.Add(1)
		// 处理连接
		bytes := []byte(msg)
		_, _ = client.Conn.Write(bytes)
		// 该客户端完成一个干活的请求
		client.Waiting.Done()
	}

}

func (handler *EchoHandler) Close() error {
	// 设置关闭状态
	handler.closing.Set(true)

	// 关闭handler中的连接资源
	handler.activeConn.Range(func(key, value interface{}) bool {

		// 因为key是空接口，需要强转类型
		client := key.(*EchoClient)
		_ = client.Close()
		// 代表继续处理
		return true
	})

	return nil
}
