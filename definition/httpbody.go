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
	ImgFileHeader *multipart.FileHeader `form:"image" binding:"required"`
	Object        string                `form:"object" binding:"required"`
	ObjectId      uint64                `form:"object_id" binding:"required"`
}
type UploadImgResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	ImgId        string `json:"img_id"`
}

//type CreatePostRequest struct {
//	PostName    string   `json:"post_name" binding:"required"`
//	PostTxt     string   `json:"post_txt" binding:"required"`
//	Zone        ZoneType `json:"zone"`
//	PostTxtHtml string   `json:"post_txt_html" binding:"required"` //帖子内容的html
//}
//type CreatePostResponse struct {
//	State        int    `json:"state"`
//	StateMessage string `json:"state_message"`
//	PostId       uint64 `json:"post_id"`
//}

//type CreateCommentRequest struct {
//	PostId     uint64 `json:"post_id" binding:"required"`
//	CommentTxt string `json:"comment_txt" binding:"required"`
//}
//type CreateCommentResponse struct {
//	State        int    `json:"state"`
//	StateMessage string `json:"state_message"`
//	CommentId    uint64 `json:"comment_id"`
//}

type CreateReplyRequest struct {
	CommentId   uint64            `json:"comment_id" binding:"required"`
	ReplyTxt    string            `json:"reply_txt" binding:"required"`
	Target      *uint64           `json:"target" binding:"required"` // 用指针目的:binding不为空，但可以传零值
	SomeoneBeAt map[uint64]string `json:"someone_be_at"`             //被@的人
}
type CreateReplyResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	ReplyId      uint64 `json:"reply_id"`
}

type CreateChatRequest struct {
	ImgFileHeader *multipart.FileHeader `form:"image"`                           // 图片 image
	AddresseeId   uint64                `form:"addressee_id" binding:"required"` // 收件人 addressee_id
	ChatTxt       string                `form:"chat_txt" binding:"required"`     // 聊天内容 chat_txt
}
type CreateChatResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	ChatId       uint64 `json:"chat_id"`
}

type CreatePostV2Request struct {
	ImgFileHeader *multipart.FileHeader `form:"image"`
	PostName      string                `form:"post_name" binding:"required"`
	PostTxt       string                `form:"post_txt" binding:"required"`
	Zone          ZoneType              `form:"zone"`
	PostTxtHtml   string                `form:"post_txt_html" binding:"required"` //帖子内容的html
	SomeoneBeAt   map[uint64]string     `form:"someone_be_at"`                    //被@的人
}
type CreatePostV2Response struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	PostId       uint64 `json:"post_id"`
}

type CreateCommentV2Request struct {
	ImgFileHeader *multipart.FileHeader `form:"image"`
	PostId        uint64                `form:"post_id" binding:"required"`
	CommentTxt    string                `form:"comment_txt" binding:"required"`
	SomeoneBeAt   map[uint64]string     `form:"someone_be_at"` //被@的人
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

type IgnoreMessagesRequest struct {
	MessageIds  []uint64    `json:"message_ids" binding:"required"`
	MessageType MessageType `json:"message_type" binding:"required"`
}
type IgnoreMessagesResponse struct {
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
	SomeoneBeAt  string    `json:"someone_be_at"` //被@的人
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
	SomeoneBeAt  string    `json:"someone_be_at"` //被@的人
}

type GetCommentByIdV2Response struct {
	State        int       `json:"state"`
	StateMessage string    `json:"state_message"`
	UId          uint64    `json:"u_id"`
	PostId       uint64    `json:"post_id"`
	CommentTxt   string    `json:"comment_txt"`
	CommentTime  time.Time `json:"comment_time"`
	ImgId        string    `json:"img_id"`
	SomeoneBeAt  string    `json:"someone_be_at"` //被@的人
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
	PostId       uint64    `json:"post_id"`
	CommentId    uint64    `json:"comment_id"`
	Target       uint64    `json:"target"`
	TargetUid    uint64    `json:"target_uid"` //回应对象的用户id
	ReplyTxt     string    `json:"reply_txt"`
	ReplyTime    time.Time `json:"reply_time"`
	SomeoneBeAt  string    `json:"someone_be_at"` //被@的人
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
	ChatId    uint64    `json:"chat_id"`
	AmISender bool      `json:"am_i_sender"`
	ChatTxt   string    `json:"chat_txt"`
	ImgId     string    `json:"img_id"`
	ChatTime  time.Time `json:"chat_time"`
	IsUnread  bool      `json:"is_unread"`
}
type GetAllChatResponse struct {
	State        int                   `json:"state"`
	StateMessage string                `json:"state_message"`
	ChatInfos    map[uint64][]ChatInfo `json:"chat_infos"` // 根据uid获取私聊消息
}

type GetChatInfoResponse struct {
	State        int        `json:"state"`
	StateMessage string     `json:"state_message"`
	ChatInfo     []ChatInfo `json:"chat_info"`
}

type GetUserStateResponse struct {
	State               int    `json:"state"`
	StateMessage        string `json:"state_message"`
	MyUserId            uint64 `json:"my_user_id"`
	DisableSendMsgTime  string `json:"disable_send_msg_time"`
	UnreadMessageNumber uint64 `json:"unread_message_number"`
	UnreadCommentNumber uint64 `json:"unread_comment_number"`
	UnreadReplyNumber   uint64 `json:"unread_reply_number"`
	UnreadChatNumber    uint64 `json:"unread_chat_number"`
	UnreadAtNumber      uint64 `json:"unread_at_number"`
}

type CommentMessage struct {
	CommentId   uint64    `json:"comment_id"`
	CommentTxt  string    `json:"comment_txt"`
	CommentUId  uint64    `json:"comment_u_id"`
	CommentTime time.Time `json:"comment_time"`
	PostId      uint64    `json:"post_id"`
	IsUnread    bool      `json:"is_unread"`
}
type GetAllCommentMessageResponse struct {
	State           int              `json:"state"`
	StateMessage    string           `json:"state_message"`
	CommentMessages []CommentMessage `json:"comment_messages"`
}

type ReplyMessage struct {
	ReplyId   uint64    `json:"reply_id"`
	PostId    uint64    `json:"post_id"`
	CommentId uint64    `json:"comment_id"`
	ReplyUId  uint64    `json:"reply_u_id"`
	ReplyTxt  string    `json:"reply_txt"`
	ReplyTime time.Time `json:"reply_time"`
	IsUnread  bool      `json:"is_unread"`
}
type GetAllReplyMessageResponse struct {
	State         int            `json:"state"`
	StateMessage  string         `json:"state_message"`
	ReplyMessages []ReplyMessage `json:"reply_messages"`
}

type AtMessage struct {
	AtId       uint64 `json:"at_id"`
	UId        uint64 `json:"u_id"`
	CallerUId  uint64 `json:"caller_u_id"`
	PostID     uint64 `json:"post_id"`
	MessageTxt string `json:"message_txt"`
	Place      string `json:"place"`
	IsUnread   bool   `json:"is_unread"`
}
type GetAllAtMessageResponse struct {
	State        int         `json:"state"`
	StateMessage string      `json:"state_message"`
	AtMessages   []AtMessage `json:"at_messages"`
}

type GetPrivacySettingResponse struct {
	State                    int    `json:"state"`
	StateMessage             string `json:"state_message"`
	PostIsPrivate            bool   `json:"post_is_private"`
	CommentAndReplyIsPrivate bool   `json:"comment_and_reply_is_private"`
	SavedPostIsPrivate       bool   `json:"saved_post_is_private"`
	SubscribedIsPrivate      bool   `json:"subscribed_is_private"`
}

type GetUidByUnameResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	UId          uint64 `json:"u_id"`
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

type BatchQueryPostRequest struct {
	PostIds []uint64 `json:"post_ids" binding:"required"`
}
type BatchQueryPostResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	Posts        []Post `json:"posts"`
}

type AllUserResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
	Users        []User `json:"users"`
}

type GetBanUserIdsResponse struct {
	State              int             `json:"state"`
	StateMessage       string          `json:"state_message"`
	BanUserIdAndReason []ForbiddenUser `json:"ban_user_id_and_reason"`
}

type AddBanUserRequest struct {
	BanUserId uint64 `json:"ban_user_id" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
}
type AddBanUserResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type CancelBanUserRequest struct {
	BanUserId uint64 `json:"ban_user_id" binding:"required"`
}
type CancelBanUserResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type GetBanIPsResponse struct {
	State          int           `json:"state"`
	StateMessage   string        `json:"state_message"`
	BanIPAndReason []ForbiddenIp `json:"ban_ip_and_reason"`
}

type AddBanIPRequest struct {
	BanIP  string `json:"ban_ip" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}
type AddBanIPResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type CancelBanIpRequest struct {
	BanIP string `json:"ban_ip" binding:"required"`
}
type CancelBanIpResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type GetPostVoteResponse struct {
	State        int             `json:"state"`
	StateMessage string          `json:"state_message"`
	VoteMessage  map[uint64]bool `json:"vote_message"` // map的键为uid 值为true则点赞 false为踩
}

type GetCommentVoteResponse struct {
	State        int             `json:"state"`
	StateMessage string          `json:"state_message"`
	VoteMessage  map[uint64]bool `json:"vote_message"` // map的键为uid 值为true则点赞 false为踩
}

type VotePostRequest struct {
	PostId uint64 `json:"post_id" binding:"required"`
	Vote   *int   `json:"vote" binding:"required"`
}
type VotePostResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type VoteCommentRequest struct {
	CommentId uint64 `json:"comment_id" binding:"required"`
	Vote      *int   `json:"vote" binding:"required"`
}
type VoteCommentResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type SilenceUserRequest struct {
	UserId             uint64 `json:"user_id" binding:"required"`
	DisableSendMsgTime string `json:"disable_send_msg_time" binding:"required"`
}
type SilenceUserResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type PostStatisticsPieChartRequest struct {
	StartTimeSTP int64 `json:"start_time" binding:"required"`
	EndTimeSTP   int64 `json:"end_time" binding:"required"`
}
type PostStatisticsPieChartResponse struct {
	State           int    `json:"state"`
	StateMessage    string `json:"state_message"`
	CountSmallTalk  uint64 `json:"count_small_talk"`
	CountStudyShare uint64 `json:"count_study_share"`
	CountMarket     uint64 `json:"count_market"`
}

type SetTopPostRequest struct {
	PostId   uint64 `json:"post_id" binding:"required"`
	Describe string `json:"describe"`
	IsTop    *bool  `json:"is_top" binding:"required"`
}
type SetTopPostResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type GetTopPostResponse struct {
	State        int       `json:"state"`
	StateMessage string    `json:"state_message"`
	TopPosts     []TopPost `json:"top_posts"` // 每小时的发帖量
}

type PostStatisticsLineChartResponse struct {
	State          int              `json:"state"`
	StateMessage   string           `json:"state_message"`
	PostCountByDay map[string]int64 `json:"post_count_by_day"` // 每天的发帖量
}

type PostStatisticsBarChartRequest struct {
	Date string `json:"date" binding:"required"` // 格式为 2016-01-02
}
type PostStatisticsBarChartResponse struct {
	State           int           `json:"state"`
	StateMessage    string        `json:"state_message"`
	PostCountByHour map[int]int64 `json:"post_count_by_day"` // 每小时的发帖量
}

type GetNeedApprovalPostResponse struct {
	State         int            `json:"state"`
	StateMessage  string         `json:"state_message"`
	ApprovalPosts []ApprovalPost `json:"approval_posts"`
}

type SetApprovalUserRequest struct {
	UserId       uint64 `json:"user_id" binding:"required"`
	NeedApproval *bool  `json:"need_approval" binding:"required"`
}
type SetApprovalUserResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}

type ApprovalPostRequest struct {
	ApprovalPostId uint64 `json:"approval_post_id" binding:"required"`
}
type ApprovalPostResponse struct {
	State        int    `json:"state"`
	StateMessage string `json:"state_message"`
}
