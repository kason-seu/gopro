package handler

import (
	"context"
	"gopro/database"
	databaseface "gopro/interface/database"
	"gopro/lib/logger"
	"gopro/lib/sync/atomic"
	"gopro/resp/conn"
	"gopro/resp/parse"
	"gopro/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

type RespHandler struct {
	activeConn sync.Map
	db         databaseface.Database
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	db := database.NewEchodatabse()
	return &RespHandler{
		db: db,
	}
}

// 关闭一个客户端
func (r *RespHandler) closeClient(conn *conn.Connection) {

	_ = conn.Close()
	// 善后
	r.db.AfterClientClose(conn)
	r.activeConn.Delete(conn)

}
func (r *RespHandler) Handle(ctx context.Context, connection net.Conn) {
	if r.closing.Get() {
		_ = r.Close()
		return
	}
	newConn := conn.NewConn(connection)
	r.activeConn.Store(newConn, struct{}{})

	// 解析器里面创建协程 for循环不断处理该客户端的请求，并返回给channel
	ch := parse.ParseStream(connection)

	for payloadChannel := range ch {
		err := payloadChannel.Err
		// 处理错误
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF ||
				strings.Contains(err.Error(), "use of closed network connection") {
				r.closeClient(newConn)
				logger.Error("client close connection" + newConn.RemoteAddr().String())
				return
			} else { // 处理协议错误
				errReply := reply.MakeStandardErrReply(err.Error())
				err := newConn.Write(errReply.ToBytes())

				if err != nil {
					logger.Error("client close connection" + newConn.RemoteAddr().String())
					return
				}
				continue
			}
		} else { // 没有错误，处理Server返回的Payload
			if payloadChannel.Data == nil {
				continue
			}
			bulkReply, ok := payloadChannel.Data.(*reply.MultiBulkReply)

			if !ok {
				logger.Info("need bulk reply")
				continue
			}

			result := r.db.Exec(newConn, bulkReply.Args)
			if result == nil {
				_ = newConn.Write(reply.MakeStandardErrReply("Err unknown").ToBytes())
				continue
			} else {
				_ = newConn.Write(result.ToBytes())
			}
		}

	}

}

// 关闭整个handler,代表redis关闭
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(func(key interface{}, value interface{}) bool {
		client := key.(*conn.Connection)

		_ = client.Close()

		return true

	})

	r.db.Close()

	return nil
}
