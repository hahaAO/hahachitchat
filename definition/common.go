package definition

const ImgDocPath = "./imgdoc"
const Socket = ":15656"

type ZoneType = uint8 // 类型别名

const ( // 帖子分区
	SmallTalk  ZoneType = 1 //闲聊
	StudyShare ZoneType = 2 //学习交流
	Market     ZoneType = 3 //交易区
)

type MessageType = uint8 //

const ( // 消息类型
	MessageTypeComment MessageType = 1 //评论
	MessageTypeReply   MessageType = 2 //回复
	MessageTypeChat    MessageType = 3 //私聊
	MessageTypeAt      MessageType = 4 //@
)

// 延迟删除消息队列
var DeleteImgChan chan string
var DeleteUnreadMessageChan chan UnreadMessage
var DeleteAtChan chan At          // 只传uid和place
var DeleteRepliesChan chan uint64 // 注意：传commentId而不是replyId
var DeleteComentsChan chan uint64 // 注意：传postId而不是commentId

type Session struct {
	Id     string //用户id
	Randid string //随机的唯一id
	Expire int    //存活时间单位为秒
}

//总热度
type Post_idandhot struct {
	Post_id  uint64 `json:"post_id"`
	Post_hot int64  `json:"post_hot"`
}
