//主函数,启动http服务，连接数据库
package main

import (
	"code/Hahachitchat/ServiceV2"
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
)

func main() {
	dataLayer.DB_conn()

	dataLayer.Redis_open()
	defer dataLayer.Redis_close()

	ServiceV2.StartService(definition.Socket)
}
