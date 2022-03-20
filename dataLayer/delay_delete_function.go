package dataLayer

import (
	"code/Hahachitchat/definition"
	"os"
)

//把要删除的图片id放进通道
func DeleteImgProduce(id string) {
	if id == "" { //空则不用发送 发送空的东西到消息队列会引发错误
		return
	}
	definition.DeleteImgChan <- id
}

//把要删除的评论id放进通道
func DeleteCommentProduce(id uint64) {
	if id == 0 {
		return
	}
	definition.DeleteCommentChan <- id
}

//把要删除的回复id放进通道
func DeleteReplyProduce(id uint64) {
	if id == 0 {
		return
	}
	definition.DeleteReplyChan <- id
}

//把要删除的消息id放进通道
func DeleteMessageProduce(message definition.Message) {
	definition.DeleteMessageChan <- message
}

//获取要删除的图片id并删除
func DeleteImgConsumer() {
	for id := range definition.DeleteImgChan {
		if err := DeleteImg(id); err != nil {
			Mqlog.Fatalln("deleteImg_consumer Remove err:", err)
		}
	}
}

func DeleteCommentConsumer() {
	for id := range definition.DeleteCommentChan {
		if code := DeleteCommentById(id); code != definition.DB_SUCCESS {
			Mqlog.Fatalln("[DeleteCommentConsum] Remove fail id: ", id)
		}
	}
}

func DeleteReplyConsumer() {
	for id := range definition.DeleteReplyChan {
		if code := DeleteReplyById(nil, id); code != definition.DB_SUCCESS {
			Mqlog.Fatalln("[DeleteReplyConsumer] Remove fail id: ", id)
		}
	}
}

func DeleteMessageConsumer() {
	for message := range definition.DeleteMessageChan {
		if code := DeleteUnreadMessage(nil, message.UId, message.MessageType, message.MessageId); code != definition.DB_SUCCESS {
			Mqlog.Fatalln("[DeleteMessageConsumer] Remove fail id: ", message)
		}
	}
}

func DeleteImg(id string) error {
	err := os.Remove("./imgdoc/" + id) //转化为路径并删除
	if err != nil {
		Mqlog.Println("deleteImg_consumer Remove err:", err) //没有删除成功有两种情况：操作出错，图片不存在
		return nil                                           //默认为图片不存在,不用再返回消息队列
	}
	Mqlog.Println("delete OK Img:", id) //删除成功
	return nil
}
