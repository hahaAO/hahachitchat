package definition

const ImgDocPath = "./imgdoc"
const Socket = ":15656"

type ZoneType = uint8 // 类型别名

const ( // 帖子分区
	SmallTalk  ZoneType = 1 //闲聊
	StudyShare ZoneType = 2 //学习交流
	Market     ZoneType = 3 //交易区
)

var DeleteImg_ch chan string

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
