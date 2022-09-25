package database

import (
	"gopro/interface/resp"
	"gopro/resp/reply"
)

type EchoDatabse struct {
}

func NewEchodatabse() *EchoDatabse {

	return &EchoDatabse{}

}

func (e *EchoDatabse) Exec(client resp.Connection, args [][]byte) resp.Reply {

	return reply.MakeMultiBulkReply(args)
}

func (e *EchoDatabse) Close() {

}

func (e *EchoDatabse) AfterClientClose(c resp.Connection) {

}
