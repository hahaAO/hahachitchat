package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	go dataLayer.RunNotificationHub() // 在线消息通知中心

	go dataLayer.LoadForbiddenConfig() // 加载封禁名单
}

func StartService(port string) {
	ServiceInit()

	r := gin.Default()
	// use prometheus metrics exporter middleware.
	r.Use(ginprom.PromMiddleware(ginprom.NewDefaultOpts()))

	// register the `/metrices` route.

	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	r.Use(HearsetMiddleWare(), ForbiddenMiddleWare())

	r.GET("/", DefaultTest)
	r.GET("/allpostid", AllPostId)
	r.GET("/top-post", GetTopPost)
	r.GET("/allcommentid/:post_id", AllCommentIdByPostId)
	r.GET("/user/:u_id", GetUserById)
	r.GET("/post/:post_id", GetPostById)
	r.GET("/comment/:comment_id", GetCommentById)
	r.GET("/reply/:reply_id", GetReplyById)
	r.GET("/allposthot", AllPostHot)
	r.GET("/getimg/:img_id", GetImg)
	r.GET("/id-by-name/:u_nickname", GetUidByUserNickname)

	r.POST("/register", Register)
	r.POST("/login", Login)

	//可能需要登录态的操作 个人资料(根据用户隐私设置判断是否展示)
	profileRoute := r.Group("/profile", SetSessionMiddleWare())
	profileRoute.GET("/user-saved/:u_id", GetUserSavedPost)
	profileRoute.GET("/subscriptions/:u_id", GetSubscriptions)
	profileRoute.GET("/allpostid-by-uid/:u_id", GetUserAllPostId)
	profileRoute.GET("/allcommentid-by-uid/:u_id", GetUserAllCommentId)
	profileRoute.GET("/allreplyid-by-uid/:u_id", GetUserAllReplyId)

	// 需要登录态的操作
	needSessionRoute := r.Group("", AuthMiddleWare())

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
	routeV2 := r.Group("/v2")
	routeV2.GET("/zone/:zone", AllPostByZone)
	routeV2.GET("/comment/:comment_id", GetCommentByIdV2)
	routeV2.POST("/posts", BatchQueryPost)

	// 点赞功能
	voteRoute := r.Group("/vote")
	voteRoute.GET("/post/:post_id", GetPostVote)
	voteRoute.GET("/comment/:comment_id", GetCommentVote)
	voteRoute.POST("/post", AuthMiddleWare(), VotePost)
	voteRoute.POST("/comment", AuthMiddleWare(), VoteComment)

	// ------------管理页面--------------
	adminRoute := r.Group("/admin")
	adminRoute.GET("/users", GetAllUser)
	adminRoute.GET("/ban-users", GetBanUser)
	adminRoute.POST("/add-ban-users", AddBanUser)
	adminRoute.POST("/cancel-ban-users",CancelBanUser)
	adminRoute.GET("/ban-ips", GetBanIPs)
	adminRoute.POST("/add-ban-ips", AddBanIP)
	adminRoute.POST("/cancel-ban-ip", CancelBanIp)
	adminRoute.POST("/silence-user", SilenceUser)
	adminRoute.POST("/delete-post", AdminDeletePostById)
	adminRoute.POST("/delete-comment", AdminDeleteCommentById)
	adminRoute.POST("/delete-reply", AdminDeleteReplyById)
	adminRoute.POST("/set-top-post", SetTopPost)
	postStatisticsRoute := adminRoute.Group("/post-statistics")
	postStatisticsRoute.GET("/line-chart", PostStatisticsLineChart)
	postStatisticsRoute.POST("/pie-chart", PostStatisticsPieChart)
	postStatisticsRoute.POST("/bar-chart", PostStatisticsBarChart)

	r.Run(port)
}
