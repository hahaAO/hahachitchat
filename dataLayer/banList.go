package dataLayer

import (
	"code/Hahachitchat/definition"
	"time"
)

func LoadForbiddenConfig() {
	for {
		newConf :=definition.Forbidden{
			ForbiddenIP:   make(map[string]struct{}),
			ForbiddenUser: make(map[uint64]struct{}),
		}


		code,forbiddenIp:= SelectForbiddenIp(nil)
		if code==	definition.DB_SUCCESS{
			for _, ip := range forbiddenIp {
				newConf.ForbiddenIP[ip.Ip]= struct{}{}
			}
		}else {
			DBlog.Println("[LoadForbiddenConfig] err")
		}

		code,forbiddenUser:= SelectForbiddenUser(nil)
		if code==	definition.DB_SUCCESS{
			for _, user := range forbiddenUser {
				newConf.ForbiddenUser[user.UserId]= struct{}{}
			}
		}else {
			DBlog.Println("[LoadForbiddenConfig] err")
		}
		definition.ForbiddenConfig=newConf

		time.Sleep(5*time.Second)
	}
}
