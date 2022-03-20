//对应数据库中表的结构
//不使用外键
package definition

import "time"

type User struct {
	UId        uint64    `gorm:"column:u_id; primaryKey"`                  //用户id,唯一主键
	UName      string    `gorm:"column:u_name; uniqueIndex; not null"`     //用户名,唯一索引
	UPassword  string    `gorm:"column:u_password; not null"`              //用户密码,非空
	UTime      time.Time `gorm:"column:u_time; autoCreateTime"`            //用户注册时间
	UNickname  string    `gorm:"column:u_nickname; uniqueIndex; not null"` //用户称昵,唯一索引，非空
	ImgId      string    `gorm:"column:img_id; default:defaultAvatar"`     //图片唯一id用作用户头像
	SavedPost  string    `gorm:"column:saved_post"`                        //用户收藏帖子,数组格式为:"1 2 3"
	Subscribed string    `gorm:"column:subscribed"`                        //用户关注的人,数组格式为:"1 2 3"

	//PrivacySetting 00000000 0位允许 1为禁止 用位运算的&判断
	//PrivacySetting 128 64 32 16 8关注的人 4收藏帖子 2评论和回复记录 1发帖记录
	PrivacySetting byte `gorm:"column:privacy_setting; default:0"` //用户隐私设置，为8位byte
}

func (User) TableName() string {
	return "user"
}

type Post struct {
	PostId      uint64    `gorm:"column:post_id; primaryKey" `                                           //帖子id,唯一主键
	UId         uint64    `gorm:"column:u_id; not null"`                                                 //用户id,非空
	Zone        ZoneType  `gorm:"column:zone; index; not null; default:1; check:max_checker,(zone < 4)"` //帖子分区
	PostName    string    `gorm:"column:post_name; not null"`                                            //帖子主题
	PostTxt     string    `gorm:"column:post_txt; not null"`                                             //帖子内容
	PostTime    time.Time `gorm:"column:post_time; autoCreateTime"`                                      //帖子发布时间
	PostTxtHtml string    `gorm:"column:post_txt_html"`                                                  //帖子内容的html
	ImgId       string    `gorm:"column:img_id"`                                                         //图片唯一id用作镇楼图
	SomeoneBeAt string    `gorm:"column:someone_be_at"`                                                  //被@的人的 uid 和 uNickname 以 map[uint64]string的json格式存储
}

func (Post) TableName() string {
	return "post"
}

type Comment struct {
	CommentId   uint64    `gorm:"column:comment_id; primaryKey"`       //评论id,唯一主键
	PostId      uint64    `gorm:"column:post_id; index; not null"`     //帖子id
	UId         uint64    `gorm:"column:u_id; not null"`               //用户id
	CommentTxt  string    `gorm:"column:comment_txt; not null"`        //评论内容
	CommentTime time.Time `gorm:"column:comment_time; autoCreateTime"` //评论时间
	ImgId       string    `gorm:"column:img_id"`                       //图片唯一id用作评论图
	SomeoneBeAt string    `gorm:"column:someone_be_at"`                //被@的人的 uid 和 uNickname 以 map[uint64]string的json格式存储
}

func (Comment) TableName() string {
	return "comment"
}

type Reply struct {
	ReplyId     uint64    `gorm:"column:reply_id; primaryKey" json:"reply_id"`          //回复id,唯一主键
	PostId      uint64    `gorm:"column:post_id; not null"`                             //所属帖子id
	CommentId   uint64    `gorm:"column:comment_id; index; not null" json:"comment_id"` //所属评论id
	UId         uint64    `gorm:"column:u_id; not null" json:"u_id"`                    //所属用户id
	Target      uint64    `gorm:"column:target; not null; default:0" json:"target"`     //回应对象ID（评论或回复的id），0为评论
	TargetUid   uint64    `gorm:"column:target_uid; not null" json:"target_uid"`        //回应对象的用户id
	ReplyTxt    string    `gorm:"column:reply_txt; not null" json:"reply_txt"`          //回复内容
	ReplyTime   time.Time `gorm:"column:reply_time; autoCreateTime" json:"reply_time"`  //回复时间
	SomeoneBeAt string    `gorm:"column:someone_be_at"`                                 //被@的人的 uid 和 uNickname 以 map[uint64]string的json格式存储
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

type Message struct {
	UId         uint64      `gorm:"column:u_id; index"`         // 用户id
	MessageType MessageType `gorm:"column:message_type; index"` // 消息类型
	MessageId   uint64      `gorm:"column:message_id"`          // 消息id
}

func (Message) TableName() string {
	return "Message"
}

type At struct {
	Id    uint64 `gorm:"column:id; primaryKey"`                       //回复id,唯一主键
	UId   uint64 `gorm:"column:u_id; uniqueIndex:idx_uid_and_place"`  // 被@的用户id
	Place string `gorm:"column:place; uniqueIndex:idx_uid_and_place"` // @的用户的地方,如 post_1 comment_2 这种格式
}

func (At) TableName() string {
	return "at"
}
