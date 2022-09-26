package database

import (
	"strings"
)

// 启动时初始化
var cmdTable = make(map[string]*command)

type command struct {
	executor ExecFunc
	arity    int
}

func RegisterCommand(name string, executor ExecFunc, arity int) {

	lowerName := strings.ToLower(name)

	cmdTable[lowerName] = &command{
		executor: executor,
		arity:    arity,
	}

}
