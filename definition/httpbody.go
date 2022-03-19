package definition

import (
	"mime/multipart"
	"time"
)

type CommonResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type RegisterRequest struct {
	UName     string `json:"u_name" binding:"required"`
	UPassword string `json:"u_password" binding:"required"`
	UNickname string `json:"u_nickname" binding:"required"`
}
type RegisterResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type LoginRequest struct {
	UName     string `json:"u_name" binding:"required"`
	UPassword string `json:"u_password" binding:"required"`
}
type LoginResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	UNickname    string `json:"u_nickname"`
	UId          uint64 `json:"u_id"`
}

type UploadImgRequest struct {
	Img      *multipart.FileHeader `form:"image" binding:"required"`
	Object   string                `form:"object" binding:"required"`
	ObjectId uint64                `form:"object_id" binding:"required"`
}
type UploadImgResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	ImgId        string `json:"img_id"`
}

type CreatePostRequest struct {
	PostName    string   `json:"post_name" binding:"required"`
	PostTxt     string   `json:"post_txt" binding:"required"`
	Zone        ZoneType `json:"zone"`
	PostTxtHtml string   `json:"post_txt_html" binding:"required"` //帖子内容的html
}
type CreatePostResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	PostId       uint64 `json:"post_id"`
}

type CreateCommentRequest struct {
	PostId     uint64 `json:"post_id" binding:"required"`
	CommentTxt string `json:"comment_txt" binding:"required"`
}
type CreateCommentResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	CommentId    uint64 `json:"comment_id"`
}

type CreateReplyRequest struct {
	CommentId uint64  `json:"comment_id" binding:"required"`
	ReplyTxt  string  `json:"reply_txt" binding:"required"`
	Target    *uint64 `json:"target" binding:"required"` // 用指针目的:binding不为空，但可以传零值
}
type CreateReplyResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	ReplyId      uint64 `json:"reply_id"`
}

type CreateChatResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	ChatId       uint64 `json:"chat_id"`
}

type CreatePostV2Response struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	PostId       uint64 `json:"post_id"`
}

type CreateCommentV2Response struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	CommentId    uint64 `json:"comment_id"`
}

type DeletePostByIdRequest struct {
	PostId uint64 `json:"post_id" binding:"required"`
}
type DeletePostByIdResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type DeleteCommentByIdRequest struct {
	CommentId uint64 `json:"comment_id" binding:"required"`
}
type DeleteCommentByIdResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type DeleteReplyByIdRequest struct {
	ReplyId uint64 `json:"reply_id" binding:"required"`
}
type DeleteReplyByIdResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type DeleteUnreadMessagedRequest struct {
	MessageType MessageType `json:"message_type" binding:"required"`
	MessageId   uint64      `json:"message_id" binding:"required"`
}
type DeleteUnreadMessageResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type SavePostRequest struct {
	PostId uint64 `json:"post_id" binding:"required"`
}
type SavePostResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type CancelSavePostRequest struct {
	PostId uint64 `json:"post_id" binding:"required"`
}
type CancelSavePostResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type SubscribeRequest struct {
	UserId uint64 `json:"user_id" binding:"required"`
}
type SubscribeResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type CancelSubscribeRequest struct {
	UserId uint64 `json:"user_id" binding:"required"`
}
type CancelSubscribeResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type AllPostIdResponse struct {
	State        int      `json:"state"`
	StateMessage string   `json:"state_message"`
	PostIds      []uint64 `json:"post_ids"`
}

type AllPostIdByZoneResponse struct {
	State        int      `json:"state"`
	StateMessage string   `json:"state_message"`
	PostIds      []uint64 `json:"post_ids"`
}

type GetPostByIdResponse struct {
	State        int       `json:"state"`
	StateMessage string    `json:"state_message"`
	UId          uint64    `json:"u_id"`
	PostName     string    `json:"post_name"`
	PostTxt      string    `json:"post_txt"`
	PostTime     time.Time `json:"post_time"`
	PostTxtHtml  string    `json:"post_txt_html"`
	ImgId        string    `json:"img_id"`
}

type AllCommentIdByPostIdResponse struct {
	State        int      `json:"state"`
	StateMessage string   `json:"state_message"`
	CommentIds   []uint64 `json:"comment_ids"`
}

type GetCommentByIdResponse struct {
	State        int       `json:"state"`
	StateMessage string    `json:"state_message"`
	UId          uint64    `json:"u_id"`
	PostId       uint64    `json:"post_id"`
	CommentTxt   string    `json:"comment_txt"`
	CommentTime  time.Time `json:"comment_time"`
	ImgId        string    `json:"img_id"`
}

type GetCommentByIdV2Response struct {
	State        int       `json:"state"`
	StateMessage string    `json:"state_message"`
	UId          uint64    `json:"u_id"`
	PostId       uint64    `json:"post_id"`
	CommentTxt   string    `json:"comment_txt"`
	CommentTime  time.Time `json:"comment_time"`
	ImgId        string    `json:"img_id"`
	Replies      []Reply   `json:"replies"`
}

type GetUserByIdResponse struct {
	State        int       `json:"state"`
	StateMessage string    `json:"state_message"`
	UNickname    string    `json:"u_nickname"`
	UTime        time.Time `json:"u_time"`
	ImgId        string    `json:"img_id"`
}

type GetReplyByIdResponse struct {
	State        int       `json:"state"`
	StateMessage string    `json:"state_message"`
	ReplyId      uint64    `json:"reply_id"`
	UId          uint64    `json:"u_id"`
	CommentId    uint64    `json:"comment_id"`
	Target       uint64    `json:"target"`
	TargetUid    uint64    `json:"target_uid"` //回应对象的用户id
	ReplyTxt     string    `json:"reply_txt"`
	ReplyTime    time.Time `json:"reply_time"`
}

type AllPostHotResponse struct {
	State        int             `json:"state"`
	StateMessage string          `json:"state_message"`
	HotDesc      []Post_idandhot `json:"hot_desc"`
}

type GetUserAllPostIdResponse struct {
	State        int      `json:"state"`
	StateMessage string   `json:"state_message"`
	PostIds      []uint64 `json:"post_ids"`
}

type GetUserAllCommentIdResponse struct {
	State        int      `json:"state"`
	StateMessage string   `json:"state_message"`
	CommentIds   []uint64 `json:"comment_ids"`
}

type GetUserAllReplyIdResponse struct {
	State        int      `json:"state"`
	StateMessage string   `json:"state_message"`
	ReplyIds     []uint64 `json:"reply_ids"`
}

type GetUserSavedPostResponse struct {
	State        int      `json:"state"`
	StateMessage string   `json:"state_message"`
	PostIds      []uint64 `json:"post_ids"`
}

type GetUserSubscribedUserResponse struct {
	State        int      `json:"state"`
	StateMessage string   `json:"state_message"`
	UserIds      []uint64 `json:"user_ids"`
}

type ChatInfo struct {
	AmISender bool      `json:"am_i_sender"`
	ChatTxt   string    `json:"chat_txt"`
	ImgId     string    `json:"img_id"`
	ChatTime  time.Time `json:"chat_time"`
}
type GetAllChatResponse struct {
	State        int                   `json:"state"`
	StateMessage string                `json:"state_message"`
	ChatInfos    map[uint64][]ChatInfo `json:"chat_infos"` // 根据uid获取私聊消息
}

type GetUserStateResponse struct {
	State               int    `json:"state"`
	StateMessage        string `json:"state_message"`
	MyUserId            uint64 `json:"my_user_id"`
	UnreadMessageNumber int    `json:"unread_message_number"`
}

type GetPrivacySettingResponse struct {
	State                    int    `json:"state"`
	StateMessage             string `json:"state_message"`
	PostIsPrivate            bool   `json:"post_is_private"`
	CommentAndReplyIsPrivate bool   `json:"comment_and_reply_is_private"`
	SavedPostIsPrivate       bool   `json:"saved_post_is_private"`
	SubscribedIsPrivate      bool   `json:"subscribed_is_private"`
}

type PostPrivacySettingRequest struct {
	PostIsPrivate            *bool `json:"post_is_private"`
	CommentAndReplyIsPrivate *bool `json:"comment_and_reply_is_private"`
	SavedPostIsPrivate       *bool `json:"saved_post_is_private"`
	SubscribedIsPrivate      *bool `json:"subscribed_is_private"`
}
type PostPrivacySettingResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}
