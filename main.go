package main

import (
	"fmt"
	"gopro/config"
	"gopro/lib/logger"
	"gopro/resp/handler"
	"gopro/tcp"
	"os"
)

const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExist(filename string) bool {

	stat, err := os.Stat(filename)

	return err == nil && !stat.IsDir()

}

func main() {

	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "gopro",
		Ext:        "log",
		TimeFormat: "2022-10-01",
	})

	if fileExist(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	},
		//tcp.MakeHandler(),
		handler.MakeHandler(),
	)
	if err != nil {
		logger.Error(err)
	}
}
