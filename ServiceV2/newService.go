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
	r.GET("/allpostid-by-uid/:u_id", AllPostIdByUserId)
	r.GET("/user/:u_id", GetUserById)
	r.GET("/post/:post_id", GetPostById)
	r.GET("/comment/:comment_id", GetCommentById)
	r.GET("/reply/:reply_id", GetReplyById)
	r.GET("/allposthot", AllPostHot)
	r.GET("/getimg/:img_id", GetImg)

	r.POST("/register", Register)
	r.POST("/login", Login)

	needSessionRoute := r.Group("", AuthMiddleWare())
	// 需要登录态的操作
	needSessionRoute.GET("/subscribed-user/:u_id", GetUserSubscribedUser)
	needSessionRoute.GET("/user-saved/:u_id", GetUserSavedPost)
	needSessionRoute.GET("/allchat/", GetAllChat)

	needSessionRoute.POST("/create-post", CreatePost)
	needSessionRoute.POST("/create-comment", CreateComment)
	needSessionRoute.POST("/create-reply", CreateReply)
	needSessionRoute.POST("/delete-post", DeletePostById)
	needSessionRoute.POST("/delete-comment", DeleteCommentById)
	needSessionRoute.POST("/delete-reply", DeleteReplyById)
	needSessionRoute.POST("/uploadimg", UploadImg)
	needSessionRoute.POST("/save-post", SavePost)
	needSessionRoute.POST("/cancel-save", CancelSavePost)
	needSessionRoute.POST("/subscribe", Subscribe)
	needSessionRoute.POST("/cancel-subscribe", CancelSubscribe)

	needSessionRoute.POST("/create-chat", CreateChat)

	// ------------V2--------------
	routeV2 := r.Group("/v2")
	routeV2.GET("/zone/:zone", AllPostByZone)
	routeV2.GET("/comment/:comment_id", GetCommentByIdV2)

	r.Run(port)
}
