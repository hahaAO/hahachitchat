//日志文件路径和记录格式的设置
package dataLayer

import (
	"log"
	"os"
)

var (
	DBlog     *log.Logger
	Serverlog *log.Logger
	Redislog  *log.Logger
	Mqlog     *log.Logger
)

func init() {
	logfile, err := os.OpenFile("hahalog.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("Failed to open hahalog!")
	}
	logfile2, err := os.OpenFile("hahaMqlog.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("Failed to open hahaImglog!")
	}
	DBlog = log.New(logfile, "DBLOG", log.Ldate|log.Ltime)
	Serverlog = log.New(logfile, "ServerLOG", log.Ldate|log.Ltime)
	Redislog = log.New(logfile, "RedisLOG", log.Ldate|log.Ltime)
	Mqlog = log.New(logfile2, "MqLOG", log.Ldate|log.Ltime)
}
