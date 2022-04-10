//操作数据库的函数
package dataLayer

import (
	"code/Hahachitchat/definition"
	"code/Hahachitchat/utils"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
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
		&definition.Reply{}, &definition.Chat{}, &definition.UnreadMessage{}, &definition.At{},
		&definition.PostVote{}, &definition.CommentVote{}, &definition.PostStatistic{})
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

func AllUserMessage(db *gorm.DB) (definition.DBcode, []definition.User) {
	getDB(&db)
	var users []definition.User
	err := db.Model(&definition.User{}).Find(&users).Error
	if err != nil {
		return definition.DB_ERROR, nil //其他问题
	}
	if len(users) <= 0 {
		return definition.DB_NOEXIST, nil //其他问题
	}
	return definition.DB_EXIST, users

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

// 根据userid返回 未读消息 的数量
func CountUnreadMessageByUid(db *gorm.DB, uId uint64) (definition.DBcode, uint64, uint64, uint64, uint64) {
	getDB(&db)
	var messages []definition.UnreadMessage
	var unreadCommentNumber uint64
	var unreadReplyNumber uint64
	var unreadChatNumber uint64
	var unreadAtNumber uint64

	err := db.Model(&definition.UnreadMessage{}).
		Where("u_id = ? AND is_ignore = ?", uId, false).
		Find(&messages).Error
	if err != nil {
		DBlog.Println("[CountUnreadMessageByUid]:", err)
		return definition.DB_ERROR, 0, 0, 0, 0 //其他问题
	}

	for _, message := range messages {
		switch message.MessageType {
		case definition.MessageTypeComment:
			unreadCommentNumber++
		case definition.MessageTypeReply:
			unreadReplyNumber++
		case definition.MessageTypeChat:
			unreadChatNumber++
		case definition.MessageTypeAt:
			unreadAtNumber++
		}
	}
	return definition.DB_SUCCESS, unreadCommentNumber, unreadReplyNumber, unreadChatNumber, unreadAtNumber

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

//根据 Nickname 获取user （未注册0 已注册1 其他情况3）（User）
func SelectUserByNickname(db *gorm.DB, nickname string) (definition.DBcode, *definition.User) {
	getDB(&db)
	var user definition.User
	err := db.Model(&definition.User{}).
		Where("u_nickname = ?", nickname).
		First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return definition.DB_NOEXIST, nil //未注册
	} else if err != nil {
		DBlog.Println("[SelectUserByNickname]:", err)
		return definition.DB_ERROR, nil //其他问题
	}
	return definition.DB_EXIST, &user //已注册
}

//根据post id获取post
func SelectPostById(db *gorm.DB, postId uint64) (definition.DBcode, *definition.Post) {
	getDB(&db)
	var post definition.Post
	err := db.Model(&definition.Post{}).
		Where("post_id = ?", postId).
		First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return definition.DB_NOEXIST, nil //无此id0
	} else if err != nil {
		DBlog.Println("SelectPostOnid err:", err)
		return definition.DB_ERROR, nil //其他情况3
	}
	return definition.DB_EXIST, &post //查到有此id1
}

//根据post id 批量获取posts
func SelectPostsById(db *gorm.DB, postIds []uint64) (definition.DBcode, []definition.Post) {
	getDB(&db)
	var posts []definition.Post
	err := db.Model(&definition.Post{}).
		Where("post_id in ?", postIds).
		Find(&posts).Error
	if err != nil {
		DBlog.Println("SelectPostsById err:", err)
		return definition.DB_ERROR, nil
	}
	return definition.DB_SUCCESS, posts
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
		var comments []definition.Comment
		err := tx.Model(&definition.Comment{}).Where("u_id = ?", uId).Find(&comments).Error
		if err == gorm.ErrRecordNotFound { // 没有帖子
			return definition.DB_NOEXIST_COMMENT, nil
		} else if err != nil { // 则有其他问题
			return definition.DB_ERROR, nil
		}
		for _, comments := range comments {
			commentIds = append(commentIds, comments.CommentId)
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

// 获取用户（楼主）所有的评论消息 标记未读和已读
func GetAllCommentMessage(postUId uint64) (definition.DBcode, []definition.CommentMessage) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		var comments []definition.Comment
		if err := gormDB.Model(&definition.Comment{}).Where(" post_u_id = ? ", postUId).
			Where(" u_id != ? ", postUId).
			Find(&comments).Error; err != nil {
			DBlog.Println("[GetAllCommentMessage] err: ", err)
			return definition.DB_ERROR, nil
		}
		var unreadMessage []definition.UnreadMessage
		if err := gormDB.Model(&definition.UnreadMessage{}).Where(" u_id = ? AND message_type = ?", postUId, definition.MessageTypeComment).
			Find(&unreadMessage).Error; err != nil {
			DBlog.Println("[GetAllCommentMessage] err2: ", err)
			return definition.DB_ERROR, nil
		}

		return definition.DB_SUCCESS, utils.PackageCommentMessage(comments, unreadMessage)
	})
	if code == definition.DB_SUCCESS {
		return code, content.([]definition.CommentMessage)
	} else {
		return code, nil
	}
}

// 获取用户（被回复的人）所有的回复消息 标记未读和已读
func GetAllReplyMessage(targetUid uint64) (definition.DBcode, []definition.ReplyMessage) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		var replies []definition.Reply
		if err := gormDB.Model(&definition.Reply{}).Where(" target_uid = ? ", targetUid).
			Where(" u_id != ?", targetUid).
			Find(&replies).Error; err != nil {
			DBlog.Println("[GetAllReplyMessage] err: ", err)
			return definition.DB_ERROR, nil
		}
		var unreadMessage []definition.UnreadMessage
		if err := gormDB.Model(&definition.UnreadMessage{}).Where(" u_id = ? AND message_type = ?", targetUid, definition.MessageTypeReply).
			Find(&unreadMessage).Error; err != nil {
			DBlog.Println("[GetAllReplyMessage] err2: ", err)
			return definition.DB_ERROR, nil
		}

		return definition.DB_SUCCESS, utils.PackageReplyMessage(replies, unreadMessage)
	})
	if code == definition.DB_SUCCESS {
		return code, content.([]definition.ReplyMessage)
	} else {
		return code, nil
	}
}

// 获取用户（被@的人）所有的@消息 标记未读和已读
func GetAllAtMessage(uId uint64) (definition.DBcode, []definition.AtMessage) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		var ats []definition.At
		if err := gormDB.Model(&definition.At{}).Where(" u_id = ? ", uId).Find(&ats).Error; err != nil {
			DBlog.Println("[GetAllAtMessage] err: ", err)
			return definition.DB_ERROR, nil
		}
		var unreadMessages []definition.UnreadMessage
		if err := gormDB.Model(&definition.UnreadMessage{}).Where(" u_id = ? AND message_type = ?", uId, definition.MessageTypeAt).
			Find(&unreadMessages).Error; err != nil {
			DBlog.Println("[GetAllAtMessage] err2: ", err)
			return definition.DB_ERROR, nil
		}

		return definition.DB_SUCCESS, utils.PackageAtMessage(ats, unreadMessages)
	})
	if code == definition.DB_SUCCESS {
		return code, content.([]definition.AtMessage)
	} else {
		return code, nil
	}
}

//根据用户u_id 获取属于该用户的所有私聊chat 标记未读和已读
func GetAllChatInfosByUid(uId uint64) (definition.DBcode, map[uint64][]definition.ChatInfo) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		var chats []definition.Chat
		err := tx.Model(&definition.Chat{}).Where("sender_id = ? OR addressee_id = ?", uId, uId).Find(&chats).Error
		if err == gorm.ErrRecordNotFound { //没有私聊
			return definition.DB_NOEXIST, nil
		} else if err != nil { // 则有其他问题
			return definition.DB_ERROR, nil
		}
		if len(chats) == 0 { //没有私聊
			return definition.DB_NOEXIST, nil
		}

		// 查未读的私聊
		var unreadMessage []definition.UnreadMessage
		if err := gormDB.Model(&definition.UnreadMessage{}).Where(" u_id = ? AND message_type = ?", uId, definition.MessageTypeChat).
			Find(&unreadMessage).Error; err != nil {
			DBlog.Println("[SelectChatByuid] err: ", err)
			return definition.DB_ERROR, nil
		}

		return definition.DB_EXIST, utils.PackageChatInfos(uId, chats, unreadMessage)
	})
	if code == definition.DB_EXIST {
		return code, content.(map[uint64][]definition.ChatInfo)
	} else {
		return code, nil
	}
}

//根据用户u_id 和 对象uid 获取属于私聊chat 标记未读和已读
func GetChatInfosByUid(myUId uint64, uId uint64) (definition.DBcode, []definition.ChatInfo) {
	code, content := runTX(func(tx *gorm.DB) (definition.DBcode, interface{}) {
		var chats []definition.Chat
		err := tx.Model(&definition.Chat{}).Where("sender_id = ? AND addressee_id = ?", myUId, uId).
			Or("sender_id = ? AND addressee_id = ?", uId, myUId).
			Find(&chats).Error
		if err == gorm.ErrRecordNotFound { //没有私聊
			return definition.DB_NOEXIST, nil
		} else if err != nil { // 则有其他问题
			return definition.DB_ERROR, nil
		}
		if len(chats) == 0 { //没有私聊
			return definition.DB_NOEXIST, nil
		}

		// 查未读的私聊
		var unreadMessage []definition.UnreadMessage
		if err := gormDB.Model(&definition.UnreadMessage{}).Where(" u_id = ? AND message_type = ?", myUId, definition.MessageTypeChat).
			Find(&unreadMessage).Error; err != nil {
			DBlog.Println("[SelectChatByuid] err: ", err)
			return definition.DB_ERROR, nil
		}

		return definition.DB_EXIST, utils.PackageChatInfo(myUId, uId, chats, unreadMessage)
	})
	if code == definition.DB_EXIST {
		return code, content.([]definition.ChatInfo)
	} else {
		return code, nil
	}
}

func SelectPostVoteById(db *gorm.DB, postId uint64) (definition.DBcode, []definition.PostVote) {
	getDB(&db)
	var postVotes []definition.PostVote
	if err := db.Model(&definition.PostVote{}).Where(" post_id = ? ", postId).Find(&postVotes).Error; err != nil {
		DBlog.Println("[SelectPostVoteById] err: ", err)
		return definition.DB_ERROR, nil
	}
	return definition.DB_SUCCESS, postVotes
}

func SelectCommentVoteById(db *gorm.DB, commmentId uint64) (definition.DBcode, []definition.CommentVote) {
	getDB(&db)
	var commentVotes []definition.CommentVote
	if err := db.Model(&definition.CommentVote{}).Where(" comment_id = ? ", commmentId).Find(&commentVotes).Error; err != nil {
		DBlog.Println("[SelectCommentVoteById] err: ", err)
		return definition.DB_ERROR, nil
	}
	return definition.DB_SUCCESS, commentVotes
}

func PostZoneCount(db *gorm.DB, startTime time.Time, endTime time.Time) (definition.DBcode, uint64, uint64, uint64) {
	getDB(&db)
	var postStatistics []definition.PostStatistic
	err := db.Model(&definition.PostStatistic{}).Where("post_time BETWEEN ? AND  ?", startTime, endTime).Select("zone").Find(&postStatistics).Error
	if err != nil {
		DBlog.Println("[PostZoneCount] err: ", err)
		return definition.DB_ERROR, 0, 0, 0
	}
	countSmallTalk := uint64(0)
	countStudyShare := uint64(0)
	countMarket := uint64(0)
	for _, p := range postStatistics {
		switch p.Zone {
		case definition.SmallTalk:
			countSmallTalk++
		case definition.StudyShare:
			countStudyShare++
		case definition.Market:
			countMarket++
		}
	}
	return definition.DB_SUCCESS, countSmallTalk, countStudyShare, countMarket
}

// map[string]uint64 的 string 是 04-11 这样。uint64 是 当日发帖的数量
func PostEverydayCount(db *gorm.DB) (definition.DBcode, map[string]int64) {
	getDB(&db)
	// 从 2022 年 4 月 1 日起，每次查 1 天
	res := make(map[string]int64)
	str := "2022-04-01"
	startTime, _ := time.Parse("2006-01-02", str)
	endTime := startTime.Add(24 * 60 * 60 * time.Second)

	for {
		if startTime.After(time.Now()) {
			break
		}

		var n int64
		err := db.Model(&definition.PostStatistic{}).Where("post_time BETWEEN ? AND  ?", startTime, endTime).Count(&n).Error
		if err != nil {
			return definition.DB_ERROR, nil
		} else {
			res[startTime.Format("2006-01-02")] = n
		}

		startTime = endTime
		endTime = endTime.Add(24 * 60 * 60 * time.Second)
	}

	return definition.DB_SUCCESS, res
}

// map[int]uint64 的 int 是当天小时。uint64 是 当日发帖的数量
func PostEveryHourCount(db *gorm.DB, startTime time.Time) (definition.DBcode, map[int]int64) {
	getDB(&db)
	res := make(map[int]int64)

	for i := 0; i < 23; i++ {
		startTime = startTime.Add(time.Duration(i) * time.Hour)
		endTime := startTime.Add(time.Hour)

		var n int64
		err := db.Model(&definition.PostStatistic{}).Where("post_time BETWEEN ? AND  ?", startTime, endTime).Count(&n).Error
		if err != nil {
			return definition.DB_ERROR, nil
		} else {
			res[i] = n
		}
	}

	return definition.DB_SUCCESS, res
}
