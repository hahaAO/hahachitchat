//操作数据库的函数
package dataLayer

import (
	"code/Hahachitchat/definition"
	"code/Hahachitchat/utils"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "vgdvgd111"
	dbname   = "hahadb"
)

var gormDB *gorm.DB

//连接一个数据库，并测试连接
func DB_conn() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 同步数据库模式
	gormDB.AutoMigrate(&definition.User{}, &definition.Post{}, &definition.Comment{},
		&definition.Reply{}, &definition.Chat{}, &definition.Message{})
	gormDB.Migrator().CreateConstraint(&definition.Post{}, "max_checker")

	DBlog.Printf("Successfully connect to postgres %s!\n", dbname)
}

// 无连接时获取连接。有事务连接时用事务连接
func getDB(db **gorm.DB) {
	if *db == nil {
		*db = gormDB
	}
}

// 启动事务
func runTX(a func(tx *gorm.DB) (definition.DBcode, interface{})) (definition.DBcode, interface{}) {
	tx := gormDB.Begin()
	defer func() {
		r := recover()
		if r != nil {
			DBlog.Fatalf("[runTX] r: ", r)
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return definition.DB_ERROR_TX, nil
	}

	code, content := a(tx)
	if code != definition.DB_SUCCESS {
		tx.Rollback()
		return code, content
	}
	if err := tx.Commit().Error; err != nil {
		return definition.DB_ERROR_TX, nil
	}
	return code, content
}

//根据uid返回user
func SelectUserById(db *gorm.DB, a uint64) (definition.DBcode, *definition.User) {
	getDB(&db)
	var user definition.User
	err := db.Model(&definition.User{}).
		Where("u_id = ?", a).
		First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return definition.DB_NOEXIST, nil //未注册
	} else if err != nil {
		DBlog.Println("SelectUserById:", err)
		return definition.DB_ERROR, nil //其他问题
	}
	return definition.DB_EXIST, &user //已注册
}

// SelectReplyById 根据replyid返回reply
func SelectReplyById(db *gorm.DB, replyId uint64) (definition.DBcode, *definition.Reply) {
	getDB(&db)
	var reply definition.Reply
	err := db.Model(&definition.Reply{}).
		Where("reply_id = ?", replyId).
		First(&reply).Error
	if err == gorm.ErrRecordNotFound {
		return definition.DB_NOEXIST, nil //没有该回复
	} else if err != nil {
		DBlog.Println("SelectReplyById:", err)
		return definition.DB_ERROR, nil //其他问题
	}
	return definition.DB_EXIST, &reply
}

// 根据userid返回Message
func SelectMessageByUid(db *gorm.DB, uId uint64) (definition.DBcode, []definition.Message) {
	getDB(&db)
	var messages []definition.Message
	err := db.Model(&definition.Message{}).
		Where("u_id = ?", uId).
		Find(&messages).Error
	if err == gorm.ErrRecordNotFound {
		return definition.DB_NOEXIST, nil //没有未读消息
	} else if err != nil {
		DBlog.Println("[SelectMessageByUid]:", err)
		return definition.DB_ERROR, nil //其他问题
	}
	return definition.DB_EXIST, messages
}

//根据commendid返回reply
func SelectRepliesByCommentId(db *gorm.DB, commentId uint64) (definition.DBcode, []definition.Reply) {
	getDB(&db)
	var reply []definition.Reply
	err := db.Model(&definition.Reply{}).
		Where("comment_id = ?", commentId).
		Find(&reply).Error
	if err != nil {
		DBlog.Println("SelectRepliesByCommentId:", err)
		return definition.DB_ERROR, nil //其他问题
	}
	return definition.DB_SUCCESS, reply
}

//根据name获取user （未注册0 已注册1 其他情况3）（User）
func SelectUserByname(db *gorm.DB, name string) (definition.DBcode, *definition.User) {
	getDB(&db)
	var user definition.User
	err := db.Model(&definition.User{}).
		Where("u_name = ?", name).
		First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return definition.DB_NOEXIST, nil //未注册
	} else if err != nil {
		DBlog.Println("SelectUserOnname:", err)
		return definition.DB_ERROR, nil //其他问题
	}
	return definition.DB_EXIST, &user //已注册
}

//根据post id获取post
func SelectPostById(db *gorm.DB, post_id uint64) (definition.DBcode, *definition.Post) {
	getDB(&db)
	var post definition.Post
	err := db.Model(&definition.Post{}).
		Where("post_id = ?", post_id).
		First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return definition.DB_NOEXIST, nil //无此id0
	} else if err != nil {
		DBlog.Println("SelectPostOnid err:", err)
		return definition.DB_ERROR, nil //其他情况3
	}
	return definition.DB_EXIST, &post //查到有此id1
}

//加了读redis缓存的功能		根据comment_id获取comment
func SelectCommentById(db *gorm.DB, comment_id uint64) (definition.DBcode, *definition.Comment) {
	getDB(&db)
	scode, scomment := Redis_SelectCommentByid(comment_id) // 先读redis缓存
	if scode == definition.DB_EXIST {                      // redis中有此comment
		if scomment.PostId == 0 { // Redis中为空值
			return definition.DB_NOEXIST, nil
		} else { //Redis中存在
			return definition.DB_EXIST, &scomment
		}
	} else { // redis中无此id	或	redis出错	要到postgres中查
		var comment definition.Comment
		err := db.Model(&definition.Comment{}).
			Where("comment_id = ?", comment_id).
			First(&comment).Error
		if err == gorm.ErrRecordNotFound { //无此id0
			comment.CommentId = comment_id
			comment.PostId = 0
			go Redis_CreateComment(comment) //把数据库的comment 空值 写入redis
			return definition.DB_NOEXIST, nil
		} else if err != nil { //其他情况3
			DBlog.Println("SelectCommentOnid err:", err)
			return definition.DB_ERROR, nil
		}
		go Redis_CreateComment(comment) //把数据库的comment写入redis
		return definition.DB_EXIST, &comment
	}
}

//获取所有帖子的post
func AllSelectPost(db *gorm.DB) (definition.DBcode, []definition.Post) {
	getDB(&db)
	var posts []definition.Post
	err := db.Model(&definition.Post{}).Find(&posts).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return definition.DB_NOEXIST, nil
		}
		DBlog.Println("AllSelectPost err1:", err)
		return definition.DB_ERROR, nil
	}
	if len(posts) == 0 { //没有帖子
		return definition.DB_NOEXIST, nil
	}
	return definition.DB_EXIST, posts
}

//获取zone下所有帖子的post
func AllPostByZone(db *gorm.DB, zone definition.ZoneType) (definition.DBcode, []definition.Post) {
	getDB(&db)
	var posts []definition.Post
	err := db.Model(&definition.Post{}).Where("zone =?", zone).Find(&posts).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return definition.DB_NOEXIST, nil
		}
		DBlog.Println("AllSelectPost err1:", err)
		return definition.DB_ERROR, nil
	}
	if len(posts) == 0 { //没有帖子
		return definition.DB_NOEXIST, nil
	}
	return definition.DB_EXIST, posts
}

//根据post_id获取所有comment_id
func AllCommentIdByPostId(db *gorm.DB, post_id uint64) (definition.DBcode, []uint64) {
	getDB(&db)
	var comments []definition.Comment
	var commentids []uint64
	err := db.Model(&definition.Comment{}).Where("post_id = ?", post_id).Find(&comments).Error
	if err == gorm.ErrRecordNotFound { //没有评论
		return definition.DB_NOEXIST, nil
	} else if err != nil {
		DBlog.Println("AllCommentidOnpostid err1:", err)
		return definition.DB_ERROR, nil
	}
	for _, comment := range comments {
		commentids = append(commentids, comment.CommentId)
	}
	if len(commentids) == 0 { //没有评论
		return definition.DB_NOEXIST, nil
	}
	return definition.DB_EXIST, commentids
}

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
		Redis_DeleteCommentOnid(comment_id) //redis缓存中的也删掉
		DeleteImg_produce(comment.ImgId)    //把要删除的图片id发到消息队列
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

//根据用户u_id 获取属于该用户的所有帖子postids
func SelectAllPostIdByUid(myUid uint64, uId uint64) (definition.DBcode, []uint64) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		code, user := SelectUserById(nil, uId)
		switch code {
		case definition.DB_NOEXIST:
			return definition.DB_NOEXIST_USER, nil
		case definition.DB_ERROR:
			return definition.DB_ERROR, nil
		case definition.DB_EXIST:
			if utils.PostIsPrivate(user.PrivacySetting) && uId != myUid {
				return definition.DB_NOT_THE_OWNER, nil
			}
		default:
			return definition.DB_ERROR, nil
		}
		// 可以查询
		var postids []uint64
		var posts []definition.Post
		err := tx.Model(&definition.Post{}).Where("u_id = ?", uId).Find(&posts).Error
		if err == gorm.ErrRecordNotFound { // 没有帖子
			return definition.DB_NOEXIST_POST, nil
		} else if err != nil { // 则有其他问题
			return definition.DB_ERROR, nil
		}
		for _, post := range posts {
			postids = append(postids, post.PostId)
		}

		if len(postids) == 0 { //没有帖子
			return definition.DB_NOEXIST_POST, nil
		}
		// 则成功
		return definition.DB_EXIST, postids
	})

	if code == definition.DB_EXIST {
		return code, content.([]uint64)
	} else {
		return code, nil
	}
}

//根据用户u_id 获取属于该用户的所有评论 CommentId
func SelectAllCommentIdByUid(myUid uint64, uId uint64) (definition.DBcode, []uint64) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		code, user := SelectUserById(nil, uId)
		switch code {
		case definition.DB_NOEXIST:
			return definition.DB_NOEXIST_USER, nil
		case definition.DB_ERROR:
			return definition.DB_ERROR, nil
		case definition.DB_EXIST:
			if utils.CommentAndReplyIsPrivate(user.PrivacySetting) && uId != myUid {
				return definition.DB_NOT_THE_OWNER, nil
			}
		default:
			return definition.DB_ERROR, nil
		}
		// 可以查询
		var commentIds []uint64
		var comments []definition.Post
		err := tx.Model(&definition.Comment{}).Where("u_id = ?", uId).Find(&comments).Error
		if err == gorm.ErrRecordNotFound { // 没有帖子
			return definition.DB_NOEXIST_COMMENT, nil
		} else if err != nil { // 则有其他问题
			return definition.DB_ERROR, nil
		}
		for _, comments := range comments {
			commentIds = append(commentIds, comments.PostId)
		}

		if len(commentIds) == 0 { //没有帖子
			return definition.DB_NOEXIST_COMMENT, nil
		}
		// 则成功
		return definition.DB_EXIST, commentIds
	})

	if code == definition.DB_EXIST {
		return code, content.([]uint64)
	} else {
		return code, nil
	}
}

func SelectAllReplyIdByUid(myUid uint64, uId uint64) (definition.DBcode, []uint64) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		code, user := SelectUserById(nil, uId)
		switch code {
		case definition.DB_NOEXIST:
			return definition.DB_NOEXIST_USER, nil
		case definition.DB_ERROR:
			return definition.DB_ERROR, nil
		case definition.DB_EXIST:
			if utils.CommentAndReplyIsPrivate(user.PrivacySetting) && uId != myUid {
				return definition.DB_NOT_THE_OWNER, nil
			}
		default:
			return definition.DB_ERROR, nil
		}
		// 可以查询
		var replyIds []uint64
		var replies []definition.Reply
		err := tx.Model(&definition.Reply{}).Where("u_id = ?", uId).Find(&replies).Error
		if err == gorm.ErrRecordNotFound { // 没有回复
			return definition.DB_NOEXIST_REPLY, nil
		} else if err != nil { // 则有其他问题
			return definition.DB_ERROR, nil
		}
		for _, reply := range replies {
			replyIds = append(replyIds, reply.ReplyId)
		}

		if len(replyIds) == 0 { // 没有回复
			return definition.DB_NOEXIST_REPLY, nil
		}
		// 则成功
		return definition.DB_EXIST, replyIds
	})

	if code == definition.DB_EXIST {
		return code, content.([]uint64)
	} else {
		return code, nil
	}
}

//根据用户u_id 获取属于该用户的所有私聊chat
func SelectChatByuid(db *gorm.DB, uId uint64) (definition.DBcode, []definition.Chat) {
	getDB(&db)
	var chats []definition.Chat
	err := db.Model(&definition.Chat{}).Where("sender_id = ? OR addressee_id = ?", uId, uId).Find(&chats).Error
	if err == gorm.ErrRecordNotFound { //没有私聊
		return definition.DB_NOEXIST, nil
	} else if err != nil { // 则有其他问题
		return definition.DB_ERROR, nil
	}

	if len(chats) == 0 { //没有私聊
		return definition.DB_NOEXIST, nil
	}
	// 则成功
	return definition.DB_EXIST, chats
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
