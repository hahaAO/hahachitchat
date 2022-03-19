package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"github.com/gin-gonic/gin"
	"os"
)

func imgServiceInit() {
	os.Mkdir(definition.ImgDocPath, os.ModePerm)    //创建图片文件夹
	definition.DeleteImg_ch = make(chan string, 10) //初始化创建待删除图片消息队列
	go dataLayer.DeleteImg_consum()                 //启动一个协程去订阅id删除图片
}

func StartService(port string) {
	imgServiceInit()

	r := gin.Default()
	r.Use(HearsetMiddleWare())

	r.GET("/", DefaultTest)
	r.GET("/allpostid", AllPostId)
	r.GET("/allcommentid/:post_id", AllCommentIdByPostId)
	r.GET("/user/:u_id", GetUserById)
	r.GET("/post/:post_id", GetPostById)
	r.GET("/comment/:comment_id", GetCommentById)
	r.GET("/reply/:reply_id", GetReplyById)
	r.GET("/allposthot", AllPostHot)
	r.GET("/getimg/:img_id", GetImg)

	r.POST("/register", Register)
	r.POST("/login", Login)

	//可能需要登录态的操作 个人资料(根据用户隐私设置判断是否展示)
	profileRoute := r.Group("/profile", SetSessionMiddleWare())
	profileRoute.GET("/user-saved/:u_id", GetUserSavedPost)
	profileRoute.GET("/subscribed-user/:u_id", GetUserSubscribedUser)
	profileRoute.GET("/allpostid-by-uid/:u_id", GetUserAllPostId)
	profileRoute.GET("/allcommentid-by-uid/:u_id", GetUserAllCommentId)
	profileRoute.GET("/allreplyid-by-uid/:u_id", GetUserAllReplyId)

	// 需要登录态的操作
	needSessionRoute := r.Group("", AuthMiddleWare())

	needSessionRoute.GET("/allchat/", GetAllChat)
	needSessionRoute.GET("/user_state", GetUserState)

	needSessionRoute.POST("/create-post", CreatePost)
	needSessionRoute.POST("/create-comment", CreateComment)
	needSessionRoute.POST("/create-reply", CreateReply)
	needSessionRoute.POST("/delete-post", DeletePostById)
	needSessionRoute.POST("/delete-comment", DeleteCommentById)
	needSessionRoute.POST("/delete-reply", DeleteReplyById)
	needSessionRoute.POST("/delete-unread-message", DeleteUnreadMessage)
	needSessionRoute.POST("/uploadimg", UploadImg)
	needSessionRoute.POST("/save-post", SavePost)
	needSessionRoute.POST("/cancel-save", CancelSavePost)
	needSessionRoute.POST("/subscribe", Subscribe)
	needSessionRoute.POST("/cancel-subscribe", CancelSubscribe)
	needSessionRoute.GET("/PrivacySetting", GetPrivacySetting)
	needSessionRoute.POST("/PrivacySetting", PostPrivacySetting)

	needSessionRoute.POST("/create-chat", CreateChat)

	// ------------V2--------------
	routeV2 := r.Group("/v2")
	routeV2.GET("/zone/:zone", AllPostByZone)
	routeV2.GET("/comment/:comment_id", GetCommentByIdV2)

	routeV2.POST("/create-post", CreatePostV2)
	routeV2.POST("/create-comment", CreateCommentV2)

	r.Run(port)
}
