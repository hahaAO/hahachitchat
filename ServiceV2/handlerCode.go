package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"code/Hahachitchat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"path"
	"strconv"
)

func DefaultTest(c *gin.Context) {
	c.JSON(http.StatusOK, definition.CommonResponse{
		State:        definition.Success,
		StateMessage: "来吧展示",
	})
}

func AllPostId(c *gin.Context) {
	acode, aposts := dataLayer.AllSelectPost(nil)
	switch acode {
	case definition.DB_EXIST:
		var postIds []uint64
		along := len(aposts)
		for i := 0; i < along; i++ {
			apost := aposts[i]
			postIds = append(postIds, apost.PostId)
		}
		c.JSON(http.StatusOK, definition.AllPostIdResponse{
			State:        definition.Success,
			StateMessage: "查询成功",
			PostIds:      postIds,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.AllPostIdResponse{
			State:        definition.Success,
			StateMessage: "无帖子",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetTopPost(c *gin.Context) {
	code, topPosts := dataLayer.SelectTopPost(nil)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.GetTopPostResponse{
			State:        definition.Success,
			StateMessage: "查询成功",
			TopPosts:     topPosts,
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func AllPostByZone(c *gin.Context) {
	zoneStr := c.Param("zone")
	zone, err := utils.StrToZone(zoneStr)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}
	acode, aposts := dataLayer.AllPostByZone(nil, zone)
	switch acode {
	case definition.DB_EXIST:
		var postIds []uint64
		along := len(aposts)
		for i := 0; i < along; i++ {
			apost := aposts[i]
			postIds = append(postIds, apost.PostId)
		}
		c.JSON(http.StatusOK, definition.AllPostIdByZoneResponse{
			State:        definition.Success,
			StateMessage: "查询成功",
			PostIds:      postIds,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.AllPostIdByZoneResponse{
			State:        definition.Success,
			StateMessage: "该分区无帖子",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetUserById(c *gin.Context) {
	userIdStr := c.Param("u_id")
	uId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}
	scode, suser := dataLayer.SelectUserById(nil, uId)
	switch scode {
	case definition.DB_EXIST:
		c.Header("Cache-Control", "max-age=100") // 缓存到本地100秒
		c.JSON(http.StatusOK, definition.GetUserByIdResponse{
			State:        definition.Success,
			StateMessage: "查询用户成功",
			UNickname:    suser.UNickname,
			UTime:        suser.UTime,
			ImgId:        suser.ImgId,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetUserByIdResponse{
			State:        definition.BadRequest,
			StateMessage: "用户不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetUidByUserNickname(c *gin.Context) {
	uName := c.Param("u_nickname")
	if uName == "" {
		SetParamErrorResponse(c)
		return
	}
	code, user := dataLayer.SelectUserByNickname(nil, uName)
	switch code {
	case definition.DB_EXIST: // 已注册
		c.JSON(http.StatusOK, definition.GetUidByUnameResponse{
			State:        definition.Success,
			StateMessage: "查询用户成功",
			UId:          user.UId,
		})
		return
	case definition.DB_NOEXIST: // 未注册
		c.JSON(http.StatusOK, definition.GetUidByUnameResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
			UId:          0, // 不存在则传0
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetReplyById(c *gin.Context) {
	replyIdStr := c.Param("reply_id")
	replyId, err := strconv.ParseUint(replyIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}
	scode, sReply := dataLayer.SelectReplyById(nil, replyId)
	switch scode {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.GetReplyByIdResponse{
			State:        definition.Success,
			StateMessage: "查询回复成功",
			ReplyId:      sReply.ReplyId,
			UId:          sReply.UId,
			PostId:       sReply.PostId,
			CommentId:    sReply.CommentId,
			Target:       sReply.Target,
			TargetUid:    sReply.TargetUid,
			ReplyTxt:     sReply.ReplyTxt,
			ReplyTime:    sReply.ReplyTime,
			SomeoneBeAt:  sReply.SomeoneBeAt,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetReplyByIdResponse{
			State:        definition.BadRequest,
			StateMessage: "回复不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetPostById(c *gin.Context) {
	postIdStr := c.Param("post_id")
	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}
	scode, spost := dataLayer.SelectPostById(nil, postId)
	switch scode {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.GetPostByIdResponse{
			State:        definition.Success,
			StateMessage: "查询帖子成功",
			UId:          spost.UId,
			PostName:     spost.PostName,
			PostTxt:      spost.PostTxt,
			PostTime:     spost.PostTime,
			PostTxtHtml:  spost.PostTxtHtml,
			ImgId:        spost.ImgId,
			SomeoneBeAt:  spost.SomeoneBeAt,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetPostByIdResponse{
			State:        definition.BadRequest,
			StateMessage: "查询的帖子不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func GetCommentById(c *gin.Context) {
	commentIdStr := c.Param("comment_id")
	commentId, err := strconv.ParseUint(commentIdStr, 10, 64)
	if err != nil { // 参数不能转为int
		SetParamErrorResponse(c)
		return
	}
	scode, scomment := dataLayer.SelectCommentById(nil, commentId)
	switch scode {
	case definition.DB_EXIST:
		c.Header("Cache-Control", "max-age=100") // 缓存到本地100秒
		c.JSON(http.StatusOK, definition.GetCommentByIdResponse{
			State:        definition.Success,
			StateMessage: "查询评论成功",
			UId:          scomment.UId,
			PostId:       scomment.PostId,
			CommentTxt:   scomment.CommentTxt,
			CommentTime:  scomment.CommentTime,
			ImgId:        scomment.ImgId,
			SomeoneBeAt:  scomment.SomeoneBeAt,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetCommentByIdResponse{
			State:        definition.Success,
			StateMessage: "评论不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func GetCommentByIdV2(c *gin.Context) {
	commentIdStr := c.Param("comment_id")
	commentId, err := strconv.ParseUint(commentIdStr, 10, 64)
	if err != nil { // 参数不能转为int
		SetParamErrorResponse(c)
		return
	}
	scode, scomment := dataLayer.SelectCommentById(nil, commentId)
	switch scode {
	case definition.DB_EXIST:
		scode, sReplies := dataLayer.SelectRepliesByCommentId(nil, commentId)
		if scode != definition.DB_SUCCESS {
			SetDBErrorResponse(c)
			return
		}
		c.JSON(http.StatusOK, definition.GetCommentByIdV2Response{
			State:        definition.Success,
			StateMessage: "查询评论成功",
			UId:          scomment.UId,
			PostId:       scomment.PostId,
			CommentTxt:   scomment.CommentTxt,
			CommentTime:  scomment.CommentTime,
			ImgId:        scomment.ImgId,
			SomeoneBeAt:  scomment.SomeoneBeAt,
			Replies:      sReplies,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetCommentByIdV2Response{
			State:        definition.Success,
			StateMessage: "评论不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func BatchQueryPost(c *gin.Context) {
	var req definition.BatchQueryPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	code, posts := dataLayer.SelectPostsById(nil, req.PostIds)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.BatchQueryPostResponse{
			State:        definition.Success,
			StateMessage: "批量查询帖子成功",
			Posts:        posts,
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func AllCommentIdByPostId(c *gin.Context) {
	postIdStr := c.Param("post_id")
	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}
	acode, acommentids := dataLayer.AllCommentIdByPostId(nil, postId)
	switch acode {
	case definition.DB_EXIST: // 成功
		c.JSON(http.StatusOK, definition.AllCommentIdByPostIdResponse{
			State:        definition.Success,
			StateMessage: "查询帖子评论ID成功",
			CommentIds:   acommentids,
		})
	case definition.DB_NOEXIST: // 没有评论
		c.JSON(http.StatusOK, definition.AllCommentIdByPostIdResponse{
			State:        definition.Success,
			StateMessage: "帖子没有评论",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetUserAllPostId(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		myUid = 0 // 没有登录态
	}

	userIdStr := c.Param("u_id")
	uId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}

	scode, spostids := dataLayer.SelectAllPostIdByUid(myUid, uId)
	switch scode {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.GetUserAllPostIdResponse{
			State:        definition.Success,
			StateMessage: "查询帖子评论ID成功",
			PostIds:      spostids,
		})
	case definition.DB_NOT_THE_OWNER:
		c.JSON(http.StatusOK, definition.GetUserAllPostIdResponse{
			State:        definition.NoPermission,
			StateMessage: "该用户对发帖记录设置了仅自己可见",
		})
	case definition.DB_NOEXIST_USER:
		c.JSON(http.StatusOK, definition.GetUserAllPostIdResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_NOEXIST_POST:
		c.JSON(http.StatusOK, definition.GetUserAllPostIdResponse{
			State:        definition.Success,
			StateMessage: "该用户没发过帖子",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetUserAllCommentId(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		myUid = 0 // 没有登录态
	}

	userIdStr := c.Param("u_id")
	uId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}

	code, commentIds := dataLayer.SelectAllCommentIdByUid(myUid, uId)
	switch code {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.GetUserAllCommentIdResponse{
			State:        definition.Success,
			StateMessage: "查询用户评论记录成功",
			CommentIds:   commentIds,
		})
	case definition.DB_NOT_THE_OWNER:
		c.JSON(http.StatusOK, definition.GetUserAllCommentIdResponse{
			State:        definition.NoPermission,
			StateMessage: "该用户对评论记录设置了仅自己可见",
		})
	case definition.DB_NOEXIST_USER:
		c.JSON(http.StatusOK, definition.GetUserAllCommentIdResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_NOEXIST_COMMENT:
		c.JSON(http.StatusOK, definition.GetUserAllCommentIdResponse{
			State:        definition.Success,
			StateMessage: "该用户没发过评论",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetUserAllReplyId(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		myUid = 0 // 没有登录态
	}

	userIdStr := c.Param("u_id")
	uId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}

	code, replyIds := dataLayer.SelectAllReplyIdByUid(myUid, uId)
	switch code {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.GetUserAllReplyIdResponse{
			State:        definition.Success,
			StateMessage: "查询用户回复记录成功",
			ReplyIds:     replyIds,
		})
	case definition.DB_NOT_THE_OWNER:
		c.JSON(http.StatusOK, definition.GetUserAllReplyIdResponse{
			State:        definition.NoPermission,
			StateMessage: "该用户对回复记录设置了仅自己可见",
		})
	case definition.DB_NOEXIST_USER:
		c.JSON(http.StatusOK, definition.GetUserAllReplyIdResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_NOEXIST_REPLY:
		c.JSON(http.StatusOK, definition.GetUserAllReplyIdResponse{
			State:        definition.Success,
			StateMessage: "该用户没发过回复",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func AllPostHot(c *gin.Context) {
	hotDesc, err := dataLayer.Allposthot()
	if err != nil {
		SetDBErrorResponse(c)
	} else {
		plong := len(hotDesc)
		for i := 0; i < plong-1; i++ {
			for j := 0; j < plong-i-1; j++ {
				if hotDesc[j].Post_hot < hotDesc[j+1].Post_hot {
					hotDesc[j], hotDesc[j+1] = hotDesc[j+1], hotDesc[j]
				}
			}
		}
		c.JSON(http.StatusOK, definition.AllPostHotResponse{
			State:        definition.Success,
			StateMessage: "查询热度成功",
			HotDesc:      hotDesc,
		})
	}
}

func GetImg(c *gin.Context) {
	imgId := c.Param("img_id")
	if imgId == "" {
		SetParamErrorResponse(c)
		return
	}
	imgF := path.Join(definition.ImgDocPath, imgId)
	c.Header("Content-Type", "image/*")
	c.File(imgF)
}

func Register(c *gin.Context) {
	var req definition.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	ccode, cuser := dataLayer.CreateUser(req.UName, req.UPassword, req.UNickname, req.PasswordQuestion, req.PasswordAnswer)
	switch ccode {
	case definition.DB_ERROR_UNAME_UNIQUE: //已注册
		c.JSON(http.StatusOK, definition.RegisterResponse{
			State:        definition.BadRequest,
			StateMessage: "账号已注册",
		})
	case definition.DB_ERROR_NICKNAME_UNIQUE: //已注册
		c.JSON(http.StatusOK, definition.RegisterResponse{
			State:        definition.BadRequest,
			StateMessage: "昵称已注册",
		})
	case definition.DB_SUCCESS: //未注册
		c.JSON(http.StatusOK, definition.RegisterResponse{
			State:        definition.Success,
			StateMessage: "账号注册成功",
		})
		dataLayer.Serverlog.Println(cuser.UId, "注册成功")
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func Login(c *gin.Context) {
	var req definition.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode, suser := dataLayer.SelectUserByname(nil, req.UName)
	switch scode {
	case definition.DB_EXIST: // 已注册
		if utils.Md5(req.UPassword) == suser.UPassword { // 密码正确
			//设置cookie与session
			session := utils.CreateSession(suser.UId) //先初始化sesion
			c.SetCookie("randid", session.Randid, session.Expire,
				"/", "", false, true) // 把cookie写入响应头 设置cookie
			rcode := dataLayer.Redis_CreateSession(*session) //把session存入Redis
			if rcode != definition.DB_SUCCESS {              //设置session失败
				c.JSON(http.StatusOK, definition.LoginResponse{
					State:        definition.ServerError,
					StateMessage: "缓存会话信息出错",
				})
				return
			}
			c.JSON(http.StatusOK, definition.LoginResponse{
				State:        definition.Success,
				StateMessage: "登录成功",
				UNickname:    suser.UNickname,
				UId:          suser.UId,
			})
			return
		}
		// 密码错误
		c.JSON(http.StatusOK, definition.LoginResponse{
			State:        definition.BadRequest,
			StateMessage: "账号或密码有误",
		})
	case definition.DB_NOEXIST: // 未注册
		c.JSON(http.StatusOK, definition.LoginResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户未注册,无法登录",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func PasswordQuestion(c *gin.Context) {
	uName := c.Param("u_name")
	if uName == "" {
		SetParamErrorResponse(c)
		return
	}

	code, user := dataLayer.SelectUserByname(nil, uName)
	switch code {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.PasswordQuestionResponse{
			State:            definition.Success,
			StateMessage:     "密保问题如下",
			PasswordQuestion: user.PasswordQuestion,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.PasswordQuestionResponse{
			State:        definition.BadRequest,
			StateMessage: "没有该用户",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func ResetPassword(c *gin.Context) {
	var req definition.ResetPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.UpdatePassword(req.UName, req.NewPassword, req.PasswordAnswer)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.ResetPasswordResponse{
			State:        definition.Success,
			StateMessage: "修改密码成功",
		})
	case definition.DB_ERROR_PARAM:
		c.JSON(http.StatusOK, definition.ResetPasswordResponse{
			State:        definition.BadRequest,
			StateMessage: "密保答案不对",
		})
	case definition.DB_NOEXIST_USER:
		c.JSON(http.StatusOK, definition.ResetPasswordResponse{
			State:        definition.BadRequest,
			StateMessage: "没有该用户",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func SignOut(c *gin.Context) {
	session := utils.GetSession(c.Request)
	code := dataLayer.Redis_DeleteSession(*session)
	if code != definition.DB_SUCCESS {
		c.JSON(http.StatusOK, definition.CommonResponse{
			State:        definition.ServerError,
			StateMessage: "服务器出错",
		})
	}

	c.SetCookie("randid", *session, -1, // 设置为马上过期
		"/", "", false, true) // 把cookie写入响应头 设置cookie
	return
}

//func CreatePost(c *gin.Context) {
//	userId, exists := c.Get("u_id")
//	uId, ok := userId.(uint64)
//	if !exists || !ok {
//		SetGetUidErrorResponse(c)
//		return
//	}
//
//	var req definition.CreatePostRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		SetParamErrorResponse(c)
//		return
//	}
//
//	ccode, cpostId := dataLayer.CreatePost(uId, req.PostName, req.PostTxt, req.Zone, req.PostTxtHtml)
//	switch ccode {
//	case definition.DB_SUCCESS:
//		c.JSON(http.StatusOK, definition.CreatePostResponse{
//			State:        definition.Success,
//			StateMessage: "创建帖子成功",
//			PostId:       cpostId,
//		})
//	case definition.DB_NOEXIST: // 用户不存在
//		c.JSON(http.StatusOK, definition.CreatePostResponse{
//			State:        definition.BadRequest,
//			StateMessage: "用户不存在,无法创建帖子",
//		})
//	case definition.DB_ERROR: // 其他问题
//		SetDBErrorResponse(c)
//	default:
//		SetServerErrorResponse(c)
//	}
//}

func CreatePost(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.CreatePostV2Request
	if err := c.ShouldBind(&req); err != nil {
		dataLayer.Serverlog.Println("[CreatePostV2] err: ", err)
		SetParamErrorResponse(c)
		return
	}

	imgId := ""                   // 默认不带图片
	if req.ImgFileHeader != nil { // 带图片
		imgId = utils.TimeRandId() //图片唯一id
		filepath := path.Join(definition.ImgDocPath, imgId)
		if err := c.SaveUploadedFile(req.ImgFileHeader, filepath); err != nil {
			c.JSON(http.StatusOK, definition.CreateChatResponse{
				State:        definition.ServerError,
				StateMessage: "服务端出错,保存图片失败",
			})
			return
		}
	}

	ccode, cpostId := dataLayer.CreatePostV2(uId, req.PostName, req.PostTxt, req.Zone, req.PostTxtHtml, imgId, req.SomeoneBeAt, false)
	switch ccode {
	case definition.DB_SUCCESS_APPROVAL:
		go dataLayer.CreatePostStatistic(nil, req.Zone, imgId != "")
		c.JSON(http.StatusAccepted, definition.CreatePostV2Response{
			State:        definition.SuccessAccepted,
			StateMessage: "帖子需要审核",
			PostId:       cpostId, //审批帖子的id
		})
	case definition.DB_SUCCESS:
		go dataLayer.CreatePostStatistic(nil, req.Zone, imgId != "")
		c.JSON(http.StatusOK, definition.CreatePostV2Response{
			State:        definition.Success,
			StateMessage: "创建帖子成功",
			PostId:       cpostId,
		})
	case definition.DB_NOEXIST: // 用户不存在
		c.JSON(http.StatusOK, definition.CreatePostV2Response{
			State:        definition.BadRequest,
			StateMessage: "用户不存在,无法创建帖子",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

//func CreateComment(c *gin.Context) {
//	userId, exists := c.Get("u_id")
//	uId, ok := userId.(uint64)
//	if !exists || !ok {
//		SetGetUidErrorResponse(c)
//		return
//	}
//
//	var req definition.CreateCommentRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		SetParamErrorResponse(c)
//		return
//	}
//
//	ccode, ccomid := dataLayer.CreateComment(req.PostId, uId, req.CommentTxt)
//	switch ccode {
//	case definition.DB_SUCCESS: // 成功
//		c.JSON(http.StatusOK, definition.CreateCommentResponse{
//			State:        definition.Success,
//			StateMessage: "创建评论成功",
//			CommentId:    ccomid,
//		})
//	case definition.DB_NOEXIST_USER: // 无此人id
//		c.JSON(http.StatusOK, definition.CreateCommentResponse{
//			State:        definition.BadRequest,
//			StateMessage: "用户不存在,创建评论失败",
//		})
//	case definition.DB_NOEXIST_POST: // 无此帖子id
//		c.JSON(http.StatusOK, definition.CreateCommentResponse{
//			State:        definition.BadRequest,
//			StateMessage: "帖子不存在,创建评论失败",
//		})
//	case definition.DB_ERROR: // 失败
//		SetDBErrorResponse(c)
//	default:
//		SetServerErrorResponse(c)
//	}
//}

func CreateComment(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.CreateCommentV2Request
	if err := c.ShouldBind(&req); err != nil {
		dataLayer.Serverlog.Println("[CreateCommentV2] err: ", err)
		SetParamErrorResponse(c)
		return
	}

	imgId := ""                   // 默认不带图片
	if req.ImgFileHeader != nil { // 带图片
		imgId = utils.TimeRandId() //图片唯一id
		filepath := path.Join(definition.ImgDocPath, imgId)
		if err := c.SaveUploadedFile(req.ImgFileHeader, filepath); err != nil {
			c.JSON(http.StatusOK, definition.CreateChatResponse{
				State:        definition.ServerError,
				StateMessage: "服务端出错,保存图片失败",
			})
			return
		}
	}

	ccode, ccomid := dataLayer.CreateCommentV2(req.PostId, uId, req.CommentTxt, imgId, req.SomeoneBeAt)
	switch ccode {
	case definition.DB_SUCCESS: // 成功
		c.JSON(http.StatusOK, definition.CreateCommentV2Response{
			State:        definition.Success,
			StateMessage: "创建评论成功",
			CommentId:    ccomid,
		})
	case definition.DB_NOEXIST_USER: // 无此人id
		c.JSON(http.StatusOK, definition.CreateCommentV2Response{
			State:        definition.BadRequest,
			StateMessage: "用户不存在,创建评论失败",
		})
	case definition.DB_NOEXIST_POST: // 无此帖子id
		c.JSON(http.StatusOK, definition.CreateCommentV2Response{
			State:        definition.BadRequest,
			StateMessage: "帖子不存在,创建评论失败",
		})
	case definition.DB_ERROR: // 失败
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func CreateReply(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	ccode, cReplyId := dataLayer.CreateReply(req.CommentId, uId, req.ReplyTxt, *req.Target, req.SomeoneBeAt)
	switch ccode {
	case definition.DB_SUCCESS: // 成功
		c.JSON(http.StatusOK, definition.CreateReplyResponse{
			State:        definition.Success,
			StateMessage: "创建回复成功",
			ReplyId:      cReplyId,
		})
	case definition.DB_NOEXIST_USER: // 无此人id
		c.JSON(http.StatusOK, definition.CreateReplyResponse{
			State:        definition.BadRequest,
			StateMessage: "用户不存在,创建回复失败",
		})
	case definition.DB_NOEXIST_TARGET:
		c.JSON(http.StatusOK, definition.CreateReplyResponse{
			State:        definition.BadRequest,
			StateMessage: "回复目标不存在,创建回复失败",
		})
	case definition.DB_ERROR: // 失败
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func CreateChat(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.CreateChatRequest
	if err := c.ShouldBind(&req); err != nil {
		dataLayer.Serverlog.Println("[CreateChat] err: ", err)
		SetParamErrorResponse(c)
		return
	}

	imgId := ""                   // 默认不带图片
	if req.ImgFileHeader != nil { // 带图片
		imgId = utils.TimeRandId() //图片唯一id
		filepath := path.Join(definition.ImgDocPath, imgId)
		if err := c.SaveUploadedFile(req.ImgFileHeader, filepath); err != nil {
			c.JSON(http.StatusOK, definition.CreateChatResponse{
				State:        definition.ServerError,
				StateMessage: "服务端出错,保存图片失败",
			})
			return
		}
	}

	cCode, cChatId := dataLayer.CreateChat(uId, req.AddresseeId, req.ChatTxt, imgId)
	switch cCode {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.CreateChatResponse{
			State:        definition.Success,
			StateMessage: "发送私聊成功",
			ChatId:       cChatId,
		})
	case definition.DB_NOEXIST_USER:
		c.JSON(http.StatusOK, definition.CreateChatResponse{
			State:        definition.BadRequest,
			StateMessage: "用户不存在,发送私聊失败",
		})
	case definition.DB_NOEXIST_ADDRESSEE:
		c.JSON(http.StatusOK, definition.CreateChatResponse{
			State:        definition.BadRequest,
			StateMessage: "收信人不存在,发送私聊失败",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func DeletePostById(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.DeletePostByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}
	scode, spost := dataLayer.SelectPostById(nil, req.PostId)
	switch scode {
	case definition.DB_EXIST: // 帖子存在
		if spost.UId == uId { // 是拥有者才有权限删除
			dcode := dataLayer.DeletePostOnId(req.PostId)
			if dcode == definition.DB_SUCCESS {
				c.JSON(http.StatusOK, definition.DeletePostByIdResponse{
					State:        definition.Success,
					StateMessage: "删除成功",
				})
			} else {
				c.JSON(http.StatusOK, definition.DeletePostByIdResponse{
					State:        definition.ServerError,
					StateMessage: "删除失败",
				})
			}
		} else { // 无权删除
			c.JSON(http.StatusOK, definition.DeletePostByIdResponse{
				State:        definition.NoPermission,
				StateMessage: "无权删除",
			})
		}
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.DeletePostByIdResponse{
			State:        definition.BadRequest,
			StateMessage: "删除的帖子不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)

	}
}

func AdminDeletePostById(c *gin.Context) {
	var req definition.DeletePostByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}
	scode, _ := dataLayer.SelectPostById(nil, req.PostId)
	switch scode {
	case definition.DB_EXIST: // 帖子存在
		dcode := dataLayer.DeletePostOnId(req.PostId)
		if dcode == definition.DB_SUCCESS {
			c.JSON(http.StatusOK, definition.DeletePostByIdResponse{
				State:        definition.Success,
				StateMessage: "删除成功",
			})
		} else {
			c.JSON(http.StatusOK, definition.DeletePostByIdResponse{
				State:        definition.ServerError,
				StateMessage: "删除失败",
			})
		}
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.DeletePostByIdResponse{
			State:        definition.BadRequest,
			StateMessage: "删除的帖子不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)

	}
}

func DeleteCommentById(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.DeleteCommentByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode, scomment := dataLayer.SelectCommentById(nil, req.CommentId)
	switch scode {
	case definition.DB_EXIST: // 评论存在
		if scomment.UId == uId { // 是拥有者才有权限删除
			dcode := dataLayer.DeleteCommentById(req.CommentId)
			if dcode == definition.DB_SUCCESS {
				c.JSON(http.StatusOK, definition.DeleteCommentByIdResponse{
					State:        definition.Success,
					StateMessage: "删除成功",
				})
			} else {
				c.JSON(http.StatusOK, definition.DeleteCommentByIdResponse{
					State:        definition.ServerError,
					StateMessage: "删除失败",
				})
			}
		} else { // 无权删除
			c.JSON(http.StatusOK, definition.DeleteCommentByIdResponse{
				State:        definition.NoPermission,
				StateMessage: "无权删除",
			})
		}
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.DeleteCommentByIdResponse{
			State:        definition.BadRequest,
			StateMessage: "删除的评论不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func AdminDeleteCommentById(c *gin.Context) {
	var req definition.DeleteCommentByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode, _ := dataLayer.SelectCommentById(nil, req.CommentId)
	switch scode {
	case definition.DB_EXIST: // 评论存在
		dcode := dataLayer.DeleteCommentById(req.CommentId)
		if dcode == definition.DB_SUCCESS {
			c.JSON(http.StatusOK, definition.DeleteCommentByIdResponse{
				State:        definition.Success,
				StateMessage: "删除成功",
			})
		} else {
			c.JSON(http.StatusOK, definition.DeleteCommentByIdResponse{
				State:        definition.ServerError,
				StateMessage: "删除失败",
			})
		}
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.DeleteCommentByIdResponse{
			State:        definition.BadRequest,
			StateMessage: "删除的评论不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func DeleteReplyById(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.DeleteReplyByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode, sReply := dataLayer.SelectReplyById(nil, req.ReplyId)
	switch scode {
	case definition.DB_EXIST: // 回复存在
		if sReply.UId == uId { // 是拥有者才有权限删除
			dcode := dataLayer.DeleteReplyById(nil, req.ReplyId)
			if dcode == definition.DB_SUCCESS {
				c.JSON(http.StatusOK, definition.DeleteReplyByIdResponse{
					State:        definition.Success,
					StateMessage: "删除成功",
				})
			} else {
				c.JSON(http.StatusOK, definition.DeleteReplyByIdResponse{
					State:        definition.ServerError,
					StateMessage: "删除失败",
				})
			}
		} else { // 无权删除
			c.JSON(http.StatusOK, definition.DeleteReplyByIdResponse{
				State:        definition.NoPermission,
				StateMessage: "无权删除",
			})
		}
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.DeleteReplyByIdResponse{
			State:        definition.BadRequest,
			StateMessage: "删除的回复不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func AdminDeleteReplyById(c *gin.Context) {
	var req definition.DeleteReplyByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode, _ := dataLayer.SelectReplyById(nil, req.ReplyId)
	switch scode {
	case definition.DB_EXIST: // 回复存在
		dcode := dataLayer.DeleteReplyById(nil, req.ReplyId)
		if dcode == definition.DB_SUCCESS {
			c.JSON(http.StatusOK, definition.DeleteReplyByIdResponse{
				State:        definition.Success,
				StateMessage: "删除成功",
			})
		} else {
			c.JSON(http.StatusOK, definition.DeleteReplyByIdResponse{
				State:        definition.ServerError,
				StateMessage: "删除失败",
			})
		}
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.DeleteReplyByIdResponse{
			State:        definition.BadRequest,
			StateMessage: "删除的回复不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func ReadMessage(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.DeleteUnreadMessagedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode := dataLayer.DeleteUnreadMessage(nil, uId, req.MessageType, req.MessageId)
	switch scode {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.DeleteUnreadMessageResponse{
			State:        definition.Success,
			StateMessage: "删除成功",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func UploadImg(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.UploadImgRequest
	if err := c.ShouldBind(&req); err != nil {
		dataLayer.Serverlog.Println("[UploadImgV2] err: ", err)
		SetParamErrorResponse(c)
		return
	}

	imgId := utils.TimeRandId() //图片唯一id
	filepath := path.Join(definition.ImgDocPath, imgId)
	if err := c.SaveUploadedFile(req.ImgFileHeader, filepath); err != nil {
		c.JSON(http.StatusOK, definition.UploadImgResponse{
			State:        definition.ServerError,
			StateMessage: "服务端出错,保存图片失败",
		})
		return
	}

	sCode := dataLayer.UpdateObjectImgId(uId, req.Object, req.ObjectId, imgId)
	switch sCode {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.UploadImgResponse{
			State:        definition.Success,
			StateMessage: "上传图片成功",
			ImgId:        imgId,
		})
	case definition.DB_ERROR_PARAM:
		dataLayer.DeleteImgProduce(imgId)
		c.JSON(http.StatusOK, definition.UploadImgResponse{
			State:        definition.BadRequest,
			StateMessage: "object不正确",
		})
	case definition.DB_UNMATCH:
		dataLayer.DeleteImgProduce(imgId)
		c.JSON(http.StatusOK, definition.UploadImgResponse{
			State:        definition.BadRequest,
			StateMessage: "无权更新不属于你的头像/评论/帖子的图片",
		})
	case definition.DB_ERROR: // 其他问题
		dataLayer.DeleteImgProduce(imgId)
		SetDBErrorResponse(c)
	default:
		dataLayer.DeleteImgProduce(imgId)
		SetServerErrorResponse(c)
	}

}

func SavePost(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.SavePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode, suser := dataLayer.SelectUserById(nil, uId)
	switch scode {
	case definition.DB_EXIST:
		savedPost, err := utils.StringToArr(suser.SavedPost)
		if err != nil {
			dataLayer.Serverlog.Println("[SavePost] err: ", err)
			SetDBParamErrorResponse(c)
			return
		}

		for _, postId := range savedPost {
			if postId == req.PostId {
				c.JSON(http.StatusOK, definition.SavePostResponse{
					State:        definition.BadRequest,
					StateMessage: "帖子已经在收藏列表,无需重复提交",
				})
				return
			}
		}
		savedPost = append(savedPost, req.PostId)
		ucode := dataLayer.UpdateSavedPostByUid(nil, savedPost, uId)
		if ucode == definition.DB_SUCCESS {
			c.JSON(http.StatusOK, definition.SavePostResponse{
				State:        definition.Success,
				StateMessage: "帖子收藏成功",
			})
		} else {
			c.JSON(http.StatusOK, definition.SavePostResponse{
				State:        definition.ServerError,
				StateMessage: "帖子收藏更新失败",
			})
		}
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.SavePostResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func CancelSavePost(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.CancelSavePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode, suser := dataLayer.SelectUserById(nil, uId)
	switch scode {
	case definition.DB_EXIST:
		savedPost, err := utils.StringToArr(suser.SavedPost)
		if err != nil {
			dataLayer.Serverlog.Println("[CancelSavePost] err: ", err)
			SetDBParamErrorResponse(c)
			return
		}

		for i, postId := range savedPost {
			if postId == req.PostId {
				savedPost = append(savedPost[0:i], savedPost[i+1:]...)
				ucode := dataLayer.UpdateSavedPostByUid(nil, savedPost, uId)
				if ucode == definition.DB_SUCCESS {
					c.JSON(http.StatusOK, definition.CancelSavePostResponse{
						State:        definition.Success,
						StateMessage: "帖子取消收藏成功",
					})
				} else {
					c.JSON(http.StatusOK, definition.CancelSavePostResponse{
						State:        definition.ServerError,
						StateMessage: "帖子收藏更新失败",
					})
				}
				return
			}
		}
		c.JSON(http.StatusOK, definition.CancelSavePostResponse{
			State:        definition.BadRequest,
			StateMessage: "帖子不在收藏列表",
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.CancelSavePostResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func Subscribe(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode, suser := dataLayer.SelectUserById(nil, uId)
	switch scode {
	case definition.DB_EXIST:
		subscribed, err := utils.StringToArr(suser.Subscribed)
		if err != nil {
			dataLayer.Serverlog.Println("[Subscribe] err: ", err)
			SetDBParamErrorResponse(c)
			return
		}

		for _, subscribedUserId := range subscribed {
			if subscribedUserId == req.UserId {
				c.JSON(http.StatusOK, definition.SubscribeResponse{
					State:        definition.BadRequest,
					StateMessage: "该用户已关注,无需重复提交",
				})
				return
			}
		}
		subscribed = append(subscribed, req.UserId)
		ucode := dataLayer.UpdateSubscribedByUid(nil, subscribed, uId)
		if ucode == definition.DB_SUCCESS {
			c.JSON(http.StatusOK, definition.SubscribeResponse{
				State:        definition.Success,
				StateMessage: "关注用户成功",
			})
		} else {
			c.JSON(http.StatusOK, definition.SubscribeResponse{
				State:        definition.ServerError,
				StateMessage: "用户关注更新失败",
			})
		}
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.SubscribeResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func CancelSubscribe(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.CancelSubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode, suser := dataLayer.SelectUserById(nil, uId)
	switch scode {
	case definition.DB_EXIST:
		subscribed, err := utils.StringToArr(suser.Subscribed)
		if err != nil {
			dataLayer.Serverlog.Println("[Subscribe] err: ", err)
			SetDBParamErrorResponse(c)
			return
		}

		for i, subscribedUserId := range subscribed {
			if subscribedUserId == req.UserId {
				subscribed = append(subscribed[0:i], subscribed[i+1:]...)
				ucode := dataLayer.UpdateSubscribedByUid(nil, subscribed, uId)
				if ucode == definition.DB_SUCCESS {
					c.JSON(http.StatusOK, definition.CancelSubscribeResponse{
						State:        definition.Success,
						StateMessage: "取消关注用户成功",
					})
				} else {
					c.JSON(http.StatusOK, definition.CancelSubscribeResponse{
						State:        definition.ServerError,
						StateMessage: "用户关注更新失败",
					})
				}
				return
			}
		}
		c.JSON(http.StatusOK, definition.CancelSubscribeResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户已不在关注列表",
		})

	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.CancelSubscribeResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetPrivacySetting(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	scode, suser := dataLayer.SelectUserById(nil, uId)
	switch scode {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.GetPrivacySettingResponse{
			State:                    definition.Success,
			StateMessage:             "查询隐私设置成功",
			PostIsPrivate:            utils.PostIsPrivate(suser.PrivacySetting),
			CommentAndReplyIsPrivate: utils.CommentAndReplyIsPrivate(suser.PrivacySetting),
			SavedPostIsPrivate:       utils.SavedPostIsPrivate(suser.PrivacySetting),
			SubscribedIsPrivate:      utils.SubscribedIsPrivate(suser.PrivacySetting),
		})

	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetPrivacySettingResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func PostPrivacySetting(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.PostPrivacySettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	scode := dataLayer.UpdatePrivacySettingByUid(uId, req.PostIsPrivate, req.CommentAndReplyIsPrivate, req.SavedPostIsPrivate, req.SubscribedIsPrivate)
	switch scode {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.PostPrivacySettingResponse{
			State:        definition.Success,
			StateMessage: "更新隐私设置成功",
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.PostPrivacySettingResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func GetUserSavedPost(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		myUid = 0 // 没有登录态
	}

	userIdStr := c.Param("u_id")
	uId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}

	scode, suser := dataLayer.SelectUserById(nil, uId)
	switch scode {
	case definition.DB_EXIST:
		if utils.SavedPostIsPrivate(suser.PrivacySetting) && uId != myUid {
			c.JSON(http.StatusOK, definition.GetUserSavedPostResponse{
				State:        definition.NoPermission,
				StateMessage: "该用户对收藏夹设置了仅自己可见",
			})
			return
		}

		savedPost, err := utils.StringToArr(suser.SavedPost)
		if err != nil {
			dataLayer.Serverlog.Println("[GetUserSavedPost] err: ", err)
			SetDBParamErrorResponse(c)
			return
		}

		c.JSON(http.StatusOK, definition.GetUserSavedPostResponse{
			State:        definition.Success,
			StateMessage: "查询收藏夹成功",
			PostIds:      savedPost,
		})

	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetUserSavedPostResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetSubscriptions(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		myUid = 0 // 没有登录态
	}

	userIdStr := c.Param("u_id")
	uId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}

	scode, suser := dataLayer.SelectUserById(nil, uId)
	switch scode {
	case definition.DB_EXIST:
		if utils.SubscribedIsPrivate(suser.PrivacySetting) && uId != myUid {
			c.JSON(http.StatusOK, definition.GetUserSubscribedUserResponse{
				State:        definition.NoPermission,
				StateMessage: "该用户对关注的人设置了仅自己可见",
			})
			return
		}

		subscribed, err := utils.StringToArr(suser.Subscribed)
		if err != nil {
			dataLayer.Serverlog.Println("[GetUserSavedPost] err: ", err)
			SetDBParamErrorResponse(c)
			return
		}

		c.JSON(http.StatusOK, definition.GetUserSubscribedUserResponse{
			State:        definition.Success,
			StateMessage: "查询关注的人成功",
			UserIds:      subscribed,
		})

	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetUserSubscribedUserResponse{
			State:        definition.BadRequest,
			StateMessage: "该用户不存在",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetAllChat(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	code, chatInfos := dataLayer.GetAllChatInfosByUid(myUid)
	switch code {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.GetAllChatResponse{
			State:        definition.Success,
			StateMessage: "查询聊天记录成功",
			ChatInfos:    chatInfos,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetAllChatResponse{
			State:        definition.Success,
			StateMessage: "没有聊天记录,快去交个朋友吧",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetChatByUserId(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	uIdStr := c.Param("u_id")
	uId, err := strconv.ParseUint(uIdStr, 10, 64)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	code, chatInfo := dataLayer.GetChatInfosByUid(myUid, uId)
	switch code {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.GetChatInfoResponse{
			State:        definition.Success,
			StateMessage: "查询聊天记录成功",
			ChatInfo:     chatInfo,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.GetChatInfoResponse{
			State:        definition.Success,
			StateMessage: "没有聊天记录,快去交个朋友吧",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetUserState(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}
	code1, user := dataLayer.SelectUserById(nil, myUid)
	if code1 != definition.DB_EXIST {
		SetServerErrorResponse(c)
		return
	}
	code, unreadCommentNumber, unreadReplyNumber, unreadChatNumber, unreadAtNumber := dataLayer.CountUnreadMessageByUid(nil, myUid)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.GetUserStateResponse{
			State:               definition.Success,
			StateMessage:        "查询用户状态成功",
			MyUserId:            myUid,
			DisableSendMsgTime:  user.DisableSendMsgTime,
			UnreadMessageNumber: unreadCommentNumber + unreadReplyNumber + unreadChatNumber + unreadAtNumber,
			UnreadCommentNumber: unreadCommentNumber,
			UnreadReplyNumber:   unreadReplyNumber,
			UnreadChatNumber:    unreadChatNumber,
			UnreadAtNumber:      unreadAtNumber,
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetAllCommentMessage(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	code, commentMessages := dataLayer.GetAllCommentMessage(myUid)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.GetAllCommentMessageResponse{
			State:           definition.Success,
			StateMessage:    "查询评论信息成功",
			CommentMessages: commentMessages,
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetAllReplyMessage(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	code, replyMessages := dataLayer.GetAllReplyMessage(myUid)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.GetAllReplyMessageResponse{
			State:         definition.Success,
			StateMessage:  "查询回复信息成功",
			ReplyMessages: replyMessages,
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func GetAllAtMessage(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	code, atMessages := dataLayer.GetAllAtMessage(myUid)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.GetAllAtMessageResponse{
			State:        definition.Success,
			StateMessage: "查询被@信息成功",
			AtMessages:   atMessages,
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func IgnoreMessages(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.IgnoreMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.UpdateMessageIsIgnore(myUid, req.MessageType, req.MessageIds)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.IgnoreMessagesResponse{
			State:        definition.Success,
			StateMessage: "忽略成功",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func WebSocketConnect(c *gin.Context) {
	myUserId, exists := c.Get("u_id")
	myUid, ok := myUserId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		dataLayer.Serverlog.Println("[WebSocketConnect] err: ", err)
		return
	}

	dataLayer.RegisterNotificationHub(myUid, ws)
}

func GetPostVote(c *gin.Context) {
	postIdStr := c.Param("post_id")
	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}
	code, postVotes := dataLayer.SelectPostVoteById(nil, postId)
	switch code {
	case definition.DB_SUCCESS:
		res := make(map[uint64]bool, len(postVotes))
		for _, vote := range postVotes {
			if vote.Vote > 0 {
				res[vote.UId] = true
			} else if vote.Vote < 0 {
				res[vote.UId] = false
			}
		}
		c.JSON(http.StatusOK, definition.GetPostVoteResponse{
			State:        definition.Success,
			StateMessage: "查询投票情况成功",
			VoteMessage:  res,
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func VotePost(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.VotePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	} else if req.Vote == nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.UpdatePostVote(uId, req.PostId, *req.Vote)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.VotePostResponse{
			State:        definition.Success,
			StateMessage: "更新投票情况成功",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func GetCommentVote(c *gin.Context) {
	commentIdStr := c.Param("comment_id")
	commentId, err := strconv.ParseUint(commentIdStr, 10, 64)
	if err != nil { //参数不能转为int
		SetParamErrorResponse(c)
		return
	}
	code, commmentVotes := dataLayer.SelectCommentVoteById(nil, commentId)
	switch code {
	case definition.DB_SUCCESS:
		res := make(map[uint64]bool, len(commmentVotes))
		for _, vote := range commmentVotes {
			if vote.Vote > 0 {
				res[vote.UId] = true
			} else if vote.Vote < 0 {
				res[vote.UId] = false
			}
		}
		c.JSON(http.StatusOK, definition.GetCommentVoteResponse{
			State:        definition.Success,
			StateMessage: "查询投票情况成功",
			VoteMessage:  res,
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func VoteComment(c *gin.Context) {
	userId, exists := c.Get("u_id")
	uId, ok := userId.(uint64)
	if !exists || !ok {
		SetGetUidErrorResponse(c)
		return
	}

	var req definition.VoteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SetParamErrorResponse(c)
		return
	} else if req.Vote == nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.UpdateCommentVote(uId, req.CommentId, *req.Vote)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.VoteCommentResponse{
			State:        definition.Success,
			StateMessage: "更新投票情况成功",
		})
	case definition.DB_ERROR: // 其他问题
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}
