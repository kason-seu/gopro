package database

import "gopro/interface/resp"

type CmdLine = [][]byte

// redis的业务层，业务核心
type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close()
	AfterClientClose(c resp.Connection)
}

type DataEntity struct {
	Data interface{}
}
