//日志文件路径和记录格式的设置
package main

import (
	"log"
	"os"
)

var (
	DBlog     *log.Logger
	Serverlog *log.Logger
	Redislog  *log.Logger
)

func init() {
	logfile, err := os.OpenFile("hahalog.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("Failed to open hahalog!")
	}

	DBlog = log.New(logfile, "DBLOG", log.Ldate|log.Ltime)
	Serverlog = log.New(logfile, "ServerLOG", log.Ldate|log.Ltime)
	Redislog = log.New(logfile, "RedisLOG", log.Ldate|log.Ltime)
}
