package dataLayer

import (
	"code/Hahachitchat/definition"
	"code/Hahachitchat/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//根据name password Unickname插入user
func CreateUser(uNname string, uPassword string, uNickname string) (definition.DBcode, *definition.User) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		if code, _ := SelectUserByname(tx, uNname); code == definition.DB_EXIST {
			return definition.DB_ERROR_UNAME_UNIQUE, nil
		}
		if err := tx.Model(&definition.User{}).Where("u_nickname = ?", uNickname).First(&definition.User{}).
			Error; err != gorm.ErrRecordNotFound {
			return definition.DB_ERROR_NICKNAME_UNIQUE, nil
		}

		user := definition.User{
			UName:     uNname,
			UPassword: utils.Md5(uPassword), // 密码md5加密后存储
			UNickname: uNickname,
		}
		err := tx.Model(&definition.User{}).Create(&user).Error
		if err != nil {
			DBlog.Println("[CreateUser] err1:", err)
			return definition.DB_ERROR, nil //其他问题,注册失败
		}
		if code, user := SelectUserByname(tx, uNname); code == definition.DB_EXIST {
			return definition.DB_SUCCESS, user //注册成功
		}
		return definition.DB_ERROR, nil //其他问题,注册失败
	})
	if code == definition.DB_SUCCESS {
		return code, content.(*definition.User)
	} else {
		return code, nil
	}
}

//根据uid post_name post_txt post_txthtml插入post
func CreatePost(u_id uint64, post_name string, post_txt string, zone definition.ZoneType, post_txthtml string) (definition.DBcode, uint64) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		code, _ := SelectUserById(tx, u_id)
		switch code {
		case definition.DB_NOEXIST: // 无此人id
			return definition.DB_NOEXIST, 0
		case definition.DB_EXIST: // 有此人id
			post := definition.Post{
				UId:         u_id,
				PostName:    post_name,
				PostTxt:     post_txt,
				Zone:        zone,
				PostTxtHtml: post_txthtml,
			}
			err := tx.Model(&definition.Post{}).Create(&post).Error
			if err != nil {
				DBlog.Println("CreatePost err1:", err)
				return definition.DB_ERROR, 0 //其他问题,插入失败
			}
			return definition.DB_SUCCESS, post.PostId //1则成功
		case definition.DB_ERROR: // 其他问题,查询失败
			return definition.DB_ERROR, 0
		default:
			return definition.DB_ERROR_UNEXPECTED, 0
		}
	})
	if code == definition.DB_SUCCESS {
		return code, content.(uint64)
	} else {
		return code, 0
	}
}

//根据uid post_name post_txt post_txthtml imgId插入post
func CreatePostV2(u_id uint64, post_name string, post_txt string, zone definition.ZoneType, post_txthtml string, imgId string) (definition.DBcode, uint64) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		code, _ := SelectUserById(tx, u_id)
		switch code {
		case definition.DB_NOEXIST: // 无此人id
			return definition.DB_NOEXIST, 0
		case definition.DB_EXIST: // 有此人id
			post := definition.Post{
				UId:         u_id,
				PostName:    post_name,
				PostTxt:     post_txt,
				Zone:        zone,
				PostTxtHtml: post_txthtml,
				ImgId:       imgId,
			}
			err := tx.Model(&definition.Post{}).Create(&post).Error
			if err != nil {
				DBlog.Println("CreatePost err1:", err)
				return definition.DB_ERROR, 0 //其他问题,插入失败
			}
			return definition.DB_SUCCESS, post.PostId //1则成功
		case definition.DB_ERROR: // 其他问题,查询失败
			return definition.DB_ERROR, 0
		default:
			return definition.DB_ERROR_UNEXPECTED, 0
		}
	})
	if code == definition.DB_SUCCESS {
		return code, content.(uint64)
	} else {
		return code, 0
	}
}

//根据post_id u_id comment_txt插入comment
func CreateComment(post_id uint64, u_id uint64, comment_txt string) (definition.DBcode, uint64) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		scode, _ := SelectUserById(tx, u_id)                               //查u_id
		scode2, _ := SelectPostById(tx, post_id)                           //查post_id
		if scode == definition.DB_EXIST && scode2 == definition.DB_EXIST { // 帖子和用户存在
			comment := definition.Comment{
				PostId:     post_id,
				UId:        u_id,
				CommentTxt: comment_txt,
			}
			err := tx.Model(&definition.Comment{}).Create(&comment).Error
			if err != nil {
				DBlog.Println("CreateComment err1:", err)
				return definition.DB_ERROR, 0 // 其他问题,插入失败
			}
			return definition.DB_SUCCESS, comment.CommentId
		} else if scode == definition.DB_NOEXIST {
			return definition.DB_NOEXIST_USER, 0
		} else if scode2 == definition.DB_NOEXIST {
			return definition.DB_NOEXIST_POST, 0
		} else {
			return definition.DB_ERROR, 0
		}
	})
	if code == definition.DB_SUCCESS {
		go Redis_DeleteCommentOnid(content.(uint64)) // 把之前可能存在的空值删掉
		return code, content.(uint64)
	} else {
		return code, 0
	}
}

//根据post_id u_id comment_txt imageId插入comment
func CreateCommentV2(post_id uint64, u_id uint64, comment_txt string, imgId string) (definition.DBcode, uint64) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		scode, _ := SelectUserById(tx, u_id)                               //查u_id
		scode2, _ := SelectPostById(tx, post_id)                           //查post_id
		if scode == definition.DB_EXIST && scode2 == definition.DB_EXIST { // 帖子和用户存在
			comment := definition.Comment{
				PostId:     post_id,
				UId:        u_id,
				CommentTxt: comment_txt,
				ImgId:      imgId,
			}
			err := tx.Model(&definition.Comment{}).Create(&comment).Error
			if err != nil {
				DBlog.Println("CreateComment err1:", err)
				return definition.DB_ERROR, 0 // 其他问题,插入失败
			}
			return definition.DB_SUCCESS, comment.CommentId
		} else if scode == definition.DB_NOEXIST {
			return definition.DB_NOEXIST_USER, 0
		} else if scode2 == definition.DB_NOEXIST {
			return definition.DB_NOEXIST_POST, 0
		} else {
			return definition.DB_ERROR, 0
		}
	})
	if code == definition.DB_SUCCESS {
		go Redis_DeleteCommentOnid(content.(uint64)) // 把之前可能存在的空值删掉
		return code, content.(uint64)
	} else {
		return code, 0
	}
}

// CreateReply 根据commentId uId replyTxt target插入reply
func CreateReply(commentId uint64, uId uint64, replyTxt string, target uint64) (definition.DBcode, uint64) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		scode, _ := SelectUserById(tx, uId)                 //查u_id
		scode2, comment := SelectCommentById(tx, commentId) //查comment
		scode3, targetReply := SelectReplyById(tx, target)  //查targetReply

		var targetUid uint64
		if target == 0 { // 直接回复层主
			scode3 = definition.DB_EXIST
			targetUid = comment.UId
		} else {
			targetUid = targetReply.UId
		}
		if scode == definition.DB_EXIST && scode2 == definition.DB_EXIST && scode3 == definition.DB_EXIST { // 评论和用户和回复目标都存在
			reply := definition.Reply{
				CommentId: commentId,
				UId:       uId,
				Target:    target,
				TargetUid: targetUid,
				ReplyTxt:  replyTxt,
			}
			err := tx.Model(&definition.Reply{}).Create(&reply).Error
			if err != nil {
				DBlog.Println("[CreateReply] err1:", err)
				return definition.DB_ERROR, 0 // 其他问题,插入失败
			}
			return definition.DB_SUCCESS, reply.ReplyId
		} else if scode == definition.DB_NOEXIST {
			return definition.DB_NOEXIST_USER, 0
		} else if scode2 == definition.DB_NOEXIST || scode3 == definition.DB_NOEXIST {
			return definition.DB_NOEXIST_TARGET, 0
		} else {
			return definition.DB_ERROR, 0
		}
	})
	if code == definition.DB_SUCCESS {
		return code, content.(uint64)
	} else {
		return code, 0
	}
}

// 创建chat
func CreateChat(senderId uint64, AddresseeId uint64, ChatTxt string, ImgId string) (definition.DBcode, uint64) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		sCode, _ := SelectUserById(tx, senderId)
		switch sCode {
		case definition.DB_ERROR:
			return definition.DB_ERROR, 0
		case definition.DB_NOEXIST:
			return definition.DB_NOEXIST_USER, 0
		}
		sCode2, _ := SelectUserById(tx, AddresseeId)
		switch sCode2 {
		case definition.DB_ERROR:
			return definition.DB_ERROR, 0
		case definition.DB_NOEXIST:
			return definition.DB_NOEXIST_ADDRESSEE, 0
		}

		chat := definition.Chat{
			SenderId:    senderId,
			AddresseeId: AddresseeId,
			ChatTxt:     ChatTxt,
			ImgId:       ImgId,
		}
		err := tx.Model(&definition.Chat{}).Create(&chat).Error
		if err != nil {
			DBlog.Println("[CreateChat] err: ", err)
			return definition.DB_ERROR, 0 // 其他问题,插入失败
		} else {
			return definition.DB_SUCCESS, chat.ChatId
		}
	})
	if code == definition.DB_SUCCESS {
		return code, content.(uint64)
	} else {
		return code, 0
	}
}

// 增加未读消息
func CreateMessage(messageType definition.MessageType, messageId uint64) definition.DBcode {
	code, _ := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		var userId uint64
		switch messageType {
		case definition.MessageTypeComment:
			var comment definition.Comment
			if err := tx.Model(&definition.Comment{}).Where("comment_id = ?", messageId).First(&comment).
				Error; err != nil {
				return definition.DB_ERROR, nil
			}
			userId = comment.UId
		case definition.MessageTypeReply:
			var reply definition.Reply
			if err := tx.Model(&definition.Reply{}).Where("reply_id = ?", messageId).First(&reply).
				Error; err != nil {
				return definition.DB_ERROR, nil
			}
			userId = reply.UId
		case definition.MessageTypeChat:
			var chat definition.Chat
			if err := tx.Model(&definition.Chat{}).Where("chat_id = ?", messageId).First(&chat).
				Error; err != nil {
				return definition.DB_ERROR, nil
			}
			userId = chat.AddresseeId
		}
		unreadMessage := definition.Message{
			UId:         userId,
			MessageType: messageType,
			MessageId:   messageId,
		}
		if err := tx.Model(&definition.Message{}).Create(&unreadMessage).Error; err != nil {
			return definition.DB_ERROR, nil
		}
		return definition.DB_SUCCESS, nil
	})
	return code
}

// 删除未读消息
func DeleteUnreadMessage(db *gorm.DB, uId uint64, messageType definition.MessageType, messageId uint64) definition.DBcode {
	getDB(&db)
	var message definition.Message
	err := db.Clauses(clause.Returning{}).
		Where("u_id = ? AND message_type = ? AND message_id = ?", uId, messageType, messageId).Delete(&message).Error
	if err != nil { //有其他问题
		DBlog.Println("[DeleteUnreadMessage] err: ", err)
		return definition.DB_ERROR
	} else { //删除成功
		return definition.DB_SUCCESS
	}
}

// DeletePostOnid 根据post_id 删除帖子及帖子里的评论
func DeletePostOnid(post_id uint64) definition.DBcode {
	code, _ := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		var comments []definition.Comment
		if err := tx.Clauses(clause.Returning{}).Where("post_id = ?", post_id).Delete(&comments).Error; err != nil { //有其他问题
			for _, comment := range comments {
				go Redis_DeleteCommentOnid(comment.CommentId) //redis缓存中的也删掉
			}
			DBlog.Println("DeletePostOnid err1:", err)
			return definition.DB_ERROR, nil
		}
		for _, comment := range comments { //读出图片id
			DeleteImg_produce(comment.ImgId) //把要删除的图片id发到消息队列
		}
		var post definition.Post
		if err := tx.Clauses(clause.Returning{}).Where("post_id = ?", post_id).Delete(&post).Error; err != nil { //有其他问题
			DBlog.Println("DeletePostOnid err2:", err)
			return definition.DB_ERROR, nil
		}

		DeleteImg_produce(post.ImgId) //把要删除的图片id发到消息队列
		DBlog.Printf("DeletePostOnid:	post_id %d 删除成功\n", post_id)
		return definition.DB_SUCCESS, nil
	})
	return code
}

//redis缓存中的也删掉	根据comment_id 删除评论
func DeleteCommentById(db *gorm.DB, comment_id uint64) definition.DBcode {
	getDB(&db)
	var comment definition.Comment
	err := db.Clauses(clause.Returning{}).Where("comment_id = ?", comment_id).Delete(&comment).Error
	if err != nil { //有其他问题
		DBlog.Println("DeleteCommentById err1:", err)
		return definition.DB_ERROR
	} else { //删除成功
		go Redis_DeleteCommentOnid(comment_id) //redis缓存中的也删掉
		DeleteImg_produce(comment.ImgId)       //把要删除的图片id发到消息队列
		return definition.DB_SUCCESS
	}
}

//根据reply_id 删除评论
func DeleteReplyById(db *gorm.DB, reply_id uint64) definition.DBcode {
	getDB(&db)
	var reply definition.Reply
	err := db.Clauses(clause.Returning{}).Where("reply_id = ?", reply_id).Delete(&reply).Error
	if err != nil { //有其他问题
		DBlog.Println("DeleteReplyById err1:", err)
		return definition.DB_ERROR
	} else { //删除成功
		return definition.DB_SUCCESS
	}
}

//根据对象类型，对象id，图片id 设置对应对象的图片id:即头像or镇楼图or评论图
func UpdateObjectImgId(uId uint64, object string, objectId uint64, imgId string) definition.DBcode {
	code, _ := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		var err error
		switch object {
		case "user":
			if uId != objectId {
				return definition.DB_UNMATCH, nil
			}
			err = tx.Model(&definition.User{}).Where("u_id = ?", objectId).Update("img_id", imgId).Error
		case "post":
			var post definition.Post
			if err = tx.Model(&definition.Post{}).Where("post_id = ?", objectId).Find(&post).Error; err != nil {
				break
			}

			if uId != post.UId {
				return definition.DB_UNMATCH, nil
			}

			err = tx.Model(&definition.Post{}).Where("post_id = ?", objectId).Update("img_id", imgId).Error
		case "comment":
			var comment definition.Comment
			if err = tx.Model(&definition.Comment{}).Where("comment_id = ?", objectId).Find(&comment).Error; err != nil {
				break
			}

			if uId != comment.UId {
				return definition.DB_UNMATCH, nil
			}

			err = tx.Model(&definition.Comment{}).Where("comment_id = ?", objectId).Update("img_id", imgId).Error
			if err == nil {
				go Redis_DeleteCommentOnid(objectId) // 缓存的脏数据删掉
			}
		default:
			return definition.DB_ERROR_PARAM, nil // object不正确 报错
		}
		if err != nil {
			DBlog.Println("UpdateObjectImgId err", err)
			return definition.DB_ERROR, nil
		}
		return definition.DB_SUCCESS, nil
	})
	return code
}

// 根据用户 uid 更新 收藏帖子
func UpdateSavedPostByUid(db *gorm.DB, SavedPost []uint64, uId uint64) definition.DBcode {
	getDB(&db)
	savedPostStr := utils.ArrToString(SavedPost)
	err := db.Model(&definition.User{}).Where("u_id = ?", uId).Update("saved_post", savedPostStr).Error
	if err != nil {
		return definition.DB_ERROR
	} else {
		return definition.DB_SUCCESS
	}
}

// 根据用户 uid 更新 关注的人
func UpdateSubscribedByUid(db *gorm.DB, Subscribed []uint64, uId uint64) definition.DBcode {
	getDB(&db)
	SubscribedStr := utils.ArrToString(Subscribed)
	err := db.Model(&definition.User{}).Where("u_id = ?", uId).Update("subscribed", SubscribedStr).Error
	if err != nil {
		return definition.DB_ERROR
	} else {
		return definition.DB_SUCCESS
	}
}

// 根据用户 uid 更新 隐私设置
func UpdatePrivacySettingByUid(uId uint64, PostIsPrivate *bool, CommentAndReplyIsPrivate *bool, SavedPostIsPrivate *bool, SubscribedIsPrivate *bool) definition.DBcode {
	code, _ := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		code, user := SelectUserById(tx, uId)
		if code != definition.DB_EXIST {
			return code, nil
		}

		user.PrivacySetting = utils.GetNewPrivacySetting(user.PrivacySetting, PostIsPrivate, CommentAndReplyIsPrivate, SavedPostIsPrivate, SubscribedIsPrivate)

		err := tx.Model(&definition.User{}).Save(&user).Error
		if err != nil {
			DBlog.Println("[UpdatePrivacySettingByUid] err: ", err)
			return definition.DB_ERROR, nil
		} else {
			return definition.DB_SUCCESS, nil
		}
	})
	return code
}
