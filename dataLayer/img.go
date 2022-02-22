package dataLayer

import (
	"code/Hahachitchat/definition"
	"os"
)

//把要删除的图片id放进通道
func DeleteImg_produce(id string) {
	if id == "" { //空则不用发送 发送空的东西到消息队列会引发错误
		return
	}
	definition.DeleteImg_ch <- id
}

//获取要删除的图片id并删除
func DeleteImg_consum() {
	for id := range definition.DeleteImg_ch {
		err := DeleteImg(id)
		if err != nil {
			Imglog.Println("deleteImg_consumer Remove err:", err)
		}
	}
}

func DeleteImg(id string) error {
	err := os.Remove("./imgdoc/" + id) //转化为路径并删除
	if err != nil {
		Imglog.Println("deleteImg_consumer Remove err:", err) //没有删除成功有两种情况：操作出错，图片不存在
		return nil                                            //默认为图片不存在,不用再返回消息队列
	}
	Imglog.Println("delete OK Img:", id) //删除成功
	return nil
}
