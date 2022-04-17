package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"github.com/gin-gonic/gin"
	"os"
)

func ServiceInit() {
	os.Mkdir(definition.ImgDocPath, os.ModePerm) //创建图片文件夹

	definition.DeleteImgChan = make(chan string, 10) //初始化创建待删除图片消息队列
	go dataLayer.DeleteImgConsumer()                 //启动一个协程去订阅id删除图片

	definition.DeleteUnreadMessageChan = make(chan definition.UnreadMessage, 10) //初始化创建待删除未读消息通知消息队列
	go dataLayer.DeleteMessageConsumer()                                         //启动一个协程去订阅id删除通知未读消息

	definition.DeleteAtChan = make(chan definition.At, 10) //初始化创建待删除at消息队列
	go dataLayer.DeleteAtConsumer()                        //启动一个协程去订阅id删除at

	definition.DeleteRepliesChan = make(chan uint64, 10) //初始化创建待删除回复通知消息队列
	go dataLayer.DeleteRepliesConsumer()                 //启动一个协程去订阅id删除回复

	definition.DeleteComentsChan = make(chan uint64, 10) //初始化创建待删除评论通知消息队列
	go dataLayer.DeleteCommentsConsumer()                //启动一个协程去订阅id删除评论

	InitClientLog()

	go dataLayer.RunNotificationHub() // 在线消息通知中心

	go dataLayer.LoadForbiddenConfig() // 加载封禁名单
}

func StartService(port string) {
	ServiceInit()

	r := gin.Default()

	r.Use(HearsetMiddleWare())

	clientRoute := r.Group("", ClientLogMiddleWare(), ForbiddenMiddleWare())

	clientRoute.GET("/", DefaultTest)
	clientRoute.GET("/allpostid", AllPostId)
	clientRoute.GET("/top-post", GetTopPost)
	clientRoute.GET("/allcommentid/:post_id", AllCommentIdByPostId)
	clientRoute.GET("/user/:u_id", GetUserById)
	clientRoute.GET("/post/:post_id", GetPostById)
	clientRoute.GET("/comment/:comment_id", GetCommentById)
	clientRoute.GET("/reply/:reply_id", GetReplyById)
	clientRoute.GET("/allposthot", AllPostHot)
	clientRoute.GET("/getimg/:img_id", GetImg)
	clientRoute.GET("/id-by-name/:u_nickname", GetUidByUserNickname)

	clientRoute.POST("/register", Register)
	clientRoute.POST("/login", Login)

	//可能需要登录态的操作 个人资料(根据用户隐私设置判断是否展示)
	profileRoute := clientRoute.Group("/profile", SetSessionMiddleWare())
	profileRoute.GET("/user-saved/:u_id", GetUserSavedPost)
	profileRoute.GET("/subscriptions/:u_id", GetSubscriptions)
	profileRoute.GET("/allpostid-by-uid/:u_id", GetUserAllPostId)
	profileRoute.GET("/allcommentid-by-uid/:u_id", GetUserAllCommentId)
	profileRoute.GET("/allreplyid-by-uid/:u_id", GetUserAllReplyId)

	// 需要登录态的操作
	needSessionRoute := clientRoute.Group("", AuthMiddleWare())

	needSessionRoute.GET("/user_state", GetUserState)
	needSessionRoute.GET("/ws-connect", WebSocketConnect)

	needSessionRoute.POST("/sign-out", SignOut)
	needSessionRoute.POST("/create-post", CreatePost)
	needSessionRoute.POST("/create-comment", CreateComment)
	needSessionRoute.POST("/create-reply", CreateReply)
	needSessionRoute.POST("/create-chat", CreateChat)
	needSessionRoute.POST("/delete-post", DeletePostById)
	needSessionRoute.POST("/delete-comment", DeleteCommentById)
	needSessionRoute.POST("/delete-reply", DeleteReplyById)
	needSessionRoute.POST("/uploadimg", UploadImg)
	needSessionRoute.POST("/save-post", SavePost)
	needSessionRoute.POST("/cancel-save", CancelSavePost)
	needSessionRoute.POST("/subscribe", Subscribe)
	needSessionRoute.POST("/cancel-subscribe", CancelSubscribe)
	needSessionRoute.GET("/privacy-setting", GetPrivacySetting)
	needSessionRoute.POST("/privacy-setting", PostPrivacySetting)

	MessageRoute := needSessionRoute.Group("/message")
	MessageRoute.GET("/comment", GetAllCommentMessage)
	MessageRoute.GET("/reply", GetAllReplyMessage)
	MessageRoute.GET("/at", GetAllAtMessage)
	MessageRoute.GET("/allchat", GetAllChat)
	MessageRoute.GET("/chat-user/:u_id", GetChatByUserId)
	MessageRoute.POST("/read", ReadMessage)
	MessageRoute.POST("/ignore", IgnoreMessages)

	// ------------V2--------------
	routeV2 := clientRoute.Group("/v2")
	routeV2.GET("/zone/:zone", AllPostByZone)
	routeV2.GET("/comment/:comment_id", GetCommentByIdV2)
	routeV2.POST("/posts", BatchQueryPost)

	// 点赞功能
	voteRoute := clientRoute.Group("/vote")
	voteRoute.GET("/post/:post_id", GetPostVote)
	voteRoute.GET("/comment/:comment_id", GetCommentVote)
	voteRoute.POST("/post", AuthMiddleWare(), VotePost)
	voteRoute.POST("/comment", AuthMiddleWare(), VoteComment)

	// ------------管理页面--------------
	adminRoute := r.Group("/admin")
	adminRoute.GET("/users", GetAllUser)
	adminRoute.GET("/ban-users", GetBanUser)
	adminRoute.POST("/add-ban-users", AddBanUser)
	adminRoute.POST("/cancel-ban-users", CancelBanUser)
	adminRoute.GET("/ban-ips", GetBanIPs)
	adminRoute.POST("/add-ban-ips", AddBanIP)
	adminRoute.POST("/cancel-ban-ip", CancelBanIp)
	adminRoute.POST("/silence-user", SilenceUser)
	adminRoute.POST("/delete-post", AdminDeletePostById)
	adminRoute.POST("/delete-comment", AdminDeleteCommentById)
	adminRoute.POST("/delete-reply", AdminDeleteReplyById)
	adminRoute.POST("/set-top-post", SetTopPost)
	// 审批功能
	adminRoute.POST("/set-approval-user", SetApprovalUser)
	adminRoute.GET("/need-approval-post", GetNeedApprovalPost)
	adminRoute.POST("/approval-post", ApprovalPost)
	//日志功能
	adminRoute.GET("/log-ws-connect", LogWebSocketConnect)
	// 查看统计图
	postStatisticsRoute := adminRoute.Group("/post-statistics")
	postStatisticsRoute.GET("/line-chart", PostStatisticsLineChart)
	postStatisticsRoute.POST("/pie-chart", PostStatisticsPieChart)
	postStatisticsRoute.POST("/bar-chart", PostStatisticsBarChart)

	r.Run(port)
}
