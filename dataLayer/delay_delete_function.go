package dataLayer

import (
	"code/Hahachitchat/definition"
	"code/Hahachitchat/utils"
	"fmt"
)

// 删除图片的生产者和消费者
func DeleteImgProduce(id string) {
	if id == "" { //空则不用发送 发送空的东西到消息队列会引发错误
		return
	}
	definition.DeleteImgChan <- id
}
func DeleteImgConsumer() {
	for id := range definition.DeleteImgChan {
		if err := DeleteImg(id); err != nil {
			Mqlog.Fatalln("deleteImg_consumer Remove err:", err)
		}
	}
}

// 删除消息的生产者和消费者
func DeleteMessageProduce(uid uint64, messageType definition.MessageType, messageId uint64) {
	definition.DeleteUnreadMessageChan <- definition.UnreadMessage{
		UId:         uid,
		MessageType: messageType,
		MessageId:   messageId,
	}
}
func DeleteMessageConsumer() {
	for message := range definition.DeleteUnreadMessageChan {
		if code := DeleteUnreadMessage(nil, message.UId, message.MessageType, message.MessageId); code != definition.DB_SUCCESS {
			Mqlog.Fatalln("[DeleteMessageConsumer] Remove fail message: ", message)
		}
	}
}

// 删除At的生产者和消费者
func DeleteAtProduce(someoneBeAtStr string, placePrefix string, place_id uint64) { // placePrefix(前缀)有三种 post、comment、reply
	someoneBeAt, err := utils.StringToMap(someoneBeAtStr)
	if err != nil {
		Mqlog.Fatalln("[DeleteAtProduce] StringToMap err: ", err)
	}
	for uId, _ := range someoneBeAt {
		definition.DeleteAtChan <- definition.At{
			UId:   uId,
			Place: fmt.Sprintf("%s_%d", placePrefix, place_id),
		}
	}
}
func DeleteAtConsumer() {
	for at := range definition.DeleteAtChan {
		if code := DeleteAt(nil, at.UId, at.Place); code != definition.DB_SUCCESS {
			Mqlog.Fatalln("[DeleteAtConsumer] Remove fail at: ", at)
		}
	}
}

// 删除回复的生产者和消费者
func DeleteRepliesProduce(commentId uint64) {
	if commentId == 0 {
		return
	}
	definition.DeleteRepliesChan <- commentId
}
func DeleteRepliesConsumer() {
	for commentId := range definition.DeleteRepliesChan {
		if code := DeleteRepliesByCommentId(nil, commentId); code != definition.DB_SUCCESS {
			Mqlog.Fatalln("[DeleteRepliesByCommentId] Remove fail commentId:", commentId)
		}
	}
}

// 删除评论的生产者和消费者
func DeleteCommentsProduce(postId uint64) {
	if postId == 0 {
		return
	}
	definition.DeleteComentsChan <- postId
}
func DeleteCommentsConsumer() {
	for postId := range definition.DeleteComentsChan {
		if code := DeleteCommentsByPostId(postId); code != definition.DB_SUCCESS {
			Mqlog.Fatalln("[DeleteCommentsByPostId] Remove fail postId:", postId)
		}
	}
}
