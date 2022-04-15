//对应数据库中表的结构
//不使用外键
package definition

import "time"

type User struct {
	UId        uint64    `gorm:"column:u_id; primaryKey" json:"u_id"`                        //用户id,唯一主键
	UName      string    `gorm:"column:u_name; uniqueIndex; not null" json:"u_name"`         //用户名,唯一索引
	UPassword  string    `gorm:"column:u_password; not null" json:"u_password"`              //用户密码,非空
	UTime      time.Time `gorm:"column:u_time; autoCreateTime" json:"u_time"`                //用户注册时间
	UNickname  string    `gorm:"column:u_nickname; uniqueIndex; not null" json:"u_nickname"` //用户称昵,唯一索引，非空
	ImgId      string    `gorm:"column:img_id; default:defaultAvatar" json:"img_id"`         //图片唯一id用作用户头像
	SavedPost  string    `gorm:"column:saved_post" json:"saved_post"`                        //用户收藏帖子,数组格式为:"1 2 3"
	Subscribed string    `gorm:"column:subscribed" json:"subscribed"`                        //用户关注的人,数组格式为:"1 2 3"

	//PrivacySetting 00000000 0位允许 1为禁止 用位运算的&判断
	//PrivacySetting 128 64 32 16 8关注的人 4收藏帖子 2评论和回复记录 1发帖记录
	PrivacySetting     byte   `gorm:"column:privacy_setting; default:0" json:"privacy_setting"`              //用户隐私设置，为8位byte
	DisableSendMsgTime string `gorm:"column:disable_send_msg_time; default:''" json:"disable_send_msg_time"` // 用户禁言到什么时候
}

func (User) TableName() string {
	return "user"
}

type Post struct {
	PostId      uint64    `gorm:"column:post_id; primaryKey" json:"post_id"`                                         //帖子id,唯一主键
	UId         uint64    `gorm:"column:u_id; not null" json:"u_id"`                                                 //用户id,非空
	Zone        ZoneType  `gorm:"column:zone; index; not null; default:1; check:max_checker,(zone < 4)" json:"zone"` //帖子分区
	PostName    string    `gorm:"column:post_name; not null" json:"post_name"`                                       //帖子主题
	PostTxt     string    `gorm:"column:post_txt; not null" json:"post_txt"`                                         //帖子内容
	PostTime    time.Time `gorm:"column:post_time; autoCreateTime" json:"post_time"`                                 //帖子发布时间
	PostTxtHtml string    `gorm:"column:post_txt_html" json:"post_txt_html"`                                         //帖子内容的html
	ImgId       string    `gorm:"column:img_id" json:"img_id"`                                                       //图片唯一id用作镇楼图
	SomeoneBeAt string    `gorm:"column:someone_be_at;  default:'{}'" json:"someone_be_at"`                          //被@的人的 uid 和 uNickname 以 map[uint64]string的json格式存储
}

func (Post) TableName() string {
	return "post"
}

type Comment struct {
	CommentId   uint64    `gorm:"column:comment_id; primaryKey"`       //评论id,唯一主键
	PostId      uint64    `gorm:"column:post_id; index; not null"`     //帖子id
	PostUid     uint64    `gorm:"column:post_u_id; not null"`          //帖子主人id
	UId         uint64    `gorm:"column:u_id; not null"`               //用户id
	CommentTxt  string    `gorm:"column:comment_txt; not null"`        //评论内容
	CommentTime time.Time `gorm:"column:comment_time; autoCreateTime"` //评论时间
	ImgId       string    `gorm:"column:img_id"`                       //图片唯一id用作评论图
	SomeoneBeAt string    `gorm:"column:someone_be_at;  default:'{}'"` //被@的人的 uid 和 uNickname 以 map[uint64]string的json格式存储
}

func (Comment) TableName() string {
	return "comment"
}

type Reply struct {
	ReplyId     uint64    `gorm:"column:reply_id; primaryKey" json:"reply_id"`              //回复id,唯一主键
	PostId      uint64    `gorm:"column:post_id; not null"`                                 //所属帖子id
	CommentId   uint64    `gorm:"column:comment_id; index; not null" json:"comment_id"`     //所属评论id
	UId         uint64    `gorm:"column:u_id; not null" json:"u_id"`                        //所属用户id
	Target      uint64    `gorm:"column:target; not null; default:0" json:"target"`         //回应对象ID（评论或回复的id），0为评论
	TargetUid   uint64    `gorm:"column:target_uid; not null" json:"target_uid"`            //回应对象的用户id
	ReplyTxt    string    `gorm:"column:reply_txt; not null" json:"reply_txt"`              //回复内容
	ReplyTime   time.Time `gorm:"column:reply_time; autoCreateTime" json:"reply_time"`      //回复时间
	SomeoneBeAt string    `gorm:"column:someone_be_at;  default:'{}'" json:"someone_be_at"` //被@的人的 uid 和 uNickname 以 map[uint64]string的json格式存储
}

func (Reply) TableName() string {
	return "reply"
}

type Chat struct {
	ChatId      uint64    `gorm:"column:chat_id; primaryKey"`       //聊天记录id,唯一主键
	SenderId    uint64    `gorm:"column:sender_id; index"`          //发送人id
	AddresseeId uint64    `gorm:"column:addressee_id; index"`       //收信人id
	ChatTxt     string    `gorm:"column:chat_txt"`                  //回复内容
	ImgId       string    `gorm:"column:img_id"`                    //图片
	ChatTime    time.Time `gorm:"column:chat_time; autoCreateTime"` //回复时间
}

func (Chat) TableName() string {
	return "chat"
}

// 消息是从 “评论、回复、@、聊天” 这4种场景中产生的，对应4个表 comment、reply、at、chat
// 消息有3种状态： 未读、已读、忽略。这里存储的消息只有未读和忽略的
// 先从4种场景的表中查出所有状态的消息，再根据此表确认消息的状态
// 用法如下:
//	场景表没有---->没有消息
//	场景表有---->此表没有---->已读消息
//	场景表有---->此表有---->标记未删除---->未读消息
//	场景表有---->此表有---->标记已删除---->忽略的消息
type UnreadMessage struct {
	UId         uint64      `gorm:"column:u_id; index"`              // 用户id
	MessageType MessageType `gorm:"column:message_type; index"`      // 消息类型 4种 comment reply at chat
	MessageId   uint64      `gorm:"column:message_id"`               // 消息id
	IsIgnore    bool        `gorm:"column:is_ignore; default:false"` // 用户忽略了这条消息
}

func (UnreadMessage) TableName() string {
	return "unread_message"
}

type At struct {
	Id         uint64 `gorm:"column:id; primaryKey" json:"at_id"`                       //回复id,唯一主键
	UId        uint64 `gorm:"column:u_id; uniqueIndex:idx_uid_and_place" json:"u_id"`   // 被@的用户id
	Place      string `gorm:"column:place; uniqueIndex:idx_uid_and_place" json:"place"` // @的用户的地方,如 post_1 comment_2 这种格式
	CallerUId  uint64 `gorm:"column:caller_u_id; not null" json:"caller_u_id"`
	PostID     uint64 `gorm:"column:post_id; not null" json:"post_id"`
	MessageTxt string `gorm:"column:message_txt" json:"message_txt"`
}

func (At) TableName() string {
	return "at"
}

type PostVote struct {
	ID     uint64 `gorm:"column:id; primaryKey"`                                 //唯一主键
	PostId uint64 `gorm:"column:post_id; uniqueIndex:idx_pid_and_uid; not null"` //帖子id,
	UId    uint64 `gorm:"column:u_id; uniqueIndex:idx_pid_and_uid; not null"`    //用户id,非空
	Vote   int    `gorm:"column:vote"`                                           // 1赞同 -1反对 0无感
}

func (PostVote) TableName() string {
	return "post_vote"
}

type CommentVote struct {
	ID        uint64 `gorm:"column:id; primaryKey"`                                    //唯一主键
	CommentId uint64 `gorm:"column:comment_id; uniqueIndex:idx_cid_and_uid; not null"` //帖子id,唯一主键
	UId       uint64 `gorm:"column:u_id; uniqueIndex:idx_cid_and_uid; not null"`       //用户id,非空
	Vote      int    `gorm:"column:vote"`                                              // 1赞同 -1反对 0无感
}

func (CommentVote) TableName() string {
	return "comment_vote"
}

type PostStatistic struct {
	ID       uint64    `gorm:"column:id; primaryKey"` //唯一主键
	Zone     ZoneType  `gorm:"column:zone; index; not null" json:"zone"`
	PostTime time.Time `gorm:"column:post_time; index; autoCreateTime"`
	HaveImg  bool      `gorm:"column:have_img;"`
}

func (PostStatistic) TableName() string {
	return "post_statistic"
}

type TopPost struct {
	PostId   uint64 `gorm:"column:post_id; primaryKey" json:"post_id"` //唯一主键
	Describe string `gorm:"column:describe" json:"describe"`
}

func (TopPost) TableName() string {
	return "top_post"
}

type ForbiddenIp struct {
	Ip     string `gorm:"column:ip; primaryKey" json:"ip"` //唯一主键
	Reason string `gorm:"column:reason" json:"reason"`
}

func (ForbiddenIp) TableName() string {
	return "forbidden_ip"
}

type ForbiddenUser struct {
	UserId uint64 `gorm:"column:user_id; primaryKey" json:"user_id"` //唯一主键
	Reason string `gorm:"column:reason" json:"reason"`
}

func (ForbiddenUser) TableName() string {
	return "forbidden_user"
}
