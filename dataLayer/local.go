package dataLayer

import "os"

func DeleteImg(id string) error {
	err := os.Remove("./imgdoc/" + id) //转化为路径并删除
	if err != nil {
		Mqlog.Println("deleteImg_consumer Remove err:", err) //没有删除成功有两种情况：操作出错，图片不存在
		return nil                                           //默认为图片不存在,不用再返回消息队列
	}
	Mqlog.Println("delete OK Img:", id) //删除成功
	return nil
}
