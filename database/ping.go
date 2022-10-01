package database

import (
	"gopro/interface/resp"
	"gopro/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

func init() {

	RegisterCommand("ping", Ping, 1)

}
