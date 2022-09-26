package database

import (
	"gopro/datastruct/dict"
	"gopro/interface/database"
	"gopro/interface/resp"
	"gopro/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply
type CmdLine = [][]byte

func MakeDB() *DB {

	return &DB{
		data: dict.MakeSyncDict(),
	}

}
func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	// ping set setnx

	cmdName := strings.ToLower(string(cmdLine[0]))
	command, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeStandardErrReply("Err Unknown command" + cmdName)
	}
	if !validArityNumber(command.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	executorFunc := command.executor
	// 原始命令 SET K V -> K V
	return executorFunc(db, cmdLine[1:])
}

// SET K V   arity = 3 代表定长
// Exists k1 k2 k3 arity = -3  负号代表的是变长
func validArityNumber(arity int, cmdArgs [][]byte) bool {

	argNum := len(cmdArgs)

	if arity >= 0 {
		return arity == argNum
	} else {
		return argNum >= -arity
	}
}

// DB 包了一层dict，所以还需要把dict的接口能力包一层
func (db *DB) GetEntity(key string) (val *database.DataEntity, exists bool) {
	res, exists := db.data.Get(key)
	if !exists {
		return nil, false
	}
	val = res.(*database.DataEntity)
	return val, exists
}

func (db *DB) PutEntity(key string, val *database.DataEntity) (result int) {

	result = db.data.Put(key, val)
	return result
}

func (db *DB) PutEntityIfAbsent(key string, val *database.DataEntity) (result int) {

	result = db.data.PutIfAbsent(key, val)
	return result
}

func (db *DB) PutEntityIfExists(key string, val *database.DataEntity) (result int) {

	result = db.data.PutIfExists(key, val)
	return result
}

func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

func (db *DB) flush() {
	db.data.Clear()
}
