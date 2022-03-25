// 放一些零碎的小工具
// 不依赖数据层的增删改查
package utils

import (
	"code/Hahachitchat/definition"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	//设置随机种子
	rand.Seed(time.Now().UnixNano())
}

//根据时间戳加上随机数生成唯一图片id 或者生成唯一session
func TimeRandId() string {
	kaishi := time.Now().UnixNano()
	timeid := strconv.FormatInt(kaishi, 10)
	randid := strconv.FormatInt(rand.Int63(), 10)
	return timeid + randid
}

//输入id生成session
func CreateSession(id uint64) *definition.Session {
	return &definition.Session{
		Id:     strconv.FormatUint(id, 10), //真实id
		Randid: TimeRandId(),               //随机生成
		Expire: 3600 * 48,                  //默认两天,
	}
}

func GetSession(r *http.Request) *string { // 获取 session 的 randId
	var res string
	for _, cookienow := range r.Cookies() { //遍历所有cookie
		if cookienow.Name == "randid" { //找到的cookie("name"为"randid")
			res = cookienow.Value
		}
	}
	return &res
}

func GetFileContentType(out *os.File) (string, error) { // 获取文件类型，用来判断是否是图片
	// 只需要前 512 个字节就可以了
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func StrToZone(zoneStr string) (definition.ZoneType, error) {
	zoneInt, err := strconv.Atoi(zoneStr)
	if err != nil { //参数不能转为int
		return 0, err
	}
	zone := definition.ZoneType(zoneInt)

	if zone < definition.SmallTalk || zone > definition.Market {
		return 0, errors.New("溢出")
	}

	return zone, nil
}

// 数据库的数组以string存储 格式为 1 2 3
func ArrToString(array []uint64) string {
	if len(array) == 0 {
		return ""
	}
	str := fmt.Sprint(array)
	return str[1 : len(str)-1]
}

// 把数据库里以string存储的数组转换取出 格式为 1 2 3
func StringToArr(str string) ([]uint64, error) {
	if str == "" {
		return nil, nil
	}
	strArr := strings.Split(str, ` `)
	var arr []uint64
	for _, str := range strArr {
		n, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		arr = append(arr, n)
	}
	return arr, nil
}

// 数据库的map[uint64]string以string存储 格式为 {"0":"xixi","1":"haha"}  空则为{}
func MapToString(i map[uint64]string) string {
	if i == nil || len(i) == 0 {
		return "{}"
	}
	b, _ := json.Marshal(i)
	return string(b)
}

// 把数据库里以string存储的map[uint64]string转换取出 格式为 {"0":"xixi","1":"haha"}  空则为{}
func StringToMap(str string) (map[uint64]string, error) {
	if str == "" {
		return nil, nil
	}

	var res map[uint64]string
	if err := json.Unmarshal([]byte(str), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func PackageCommentMessage(comments []definition.Comment, unreadMessages []definition.UnreadMessage) []definition.CommentMessage {
	var res []definition.CommentMessage
	unreadMessageMap := make(map[uint64]bool, len(unreadMessages)) // 转为 hash 集合，优化算法
	for _, message := range unreadMessages {
		unreadMessageMap[message.MessageId] = message.IsIgnore
	}
	for _, comment := range comments {
		commentMessage := definition.CommentMessage{
			CommentId:   comment.CommentId,
			CommentTxt:  comment.CommentTxt,
			CommentUId:  comment.UId,
			CommentTime: comment.CommentTime,
			PostId:      comment.PostId,
			IsUnread:    false,
		}
		if IsIgnore, exist := unreadMessageMap[comment.CommentId]; exist { //是未读的评论
			if IsIgnore { // 用户忽略了消息,不加进返回
				continue
			}
			commentMessage.IsUnread = true
		}
		res = append(res, commentMessage)
	}
	return res
}

func PackageReplyMessage(replies []definition.Reply, unreadMessages []definition.UnreadMessage) []definition.ReplyMessage {
	var res []definition.ReplyMessage
	unreadMessageMap := make(map[uint64]bool, len(unreadMessages)) // 转为 hash 集合，优化算法
	for _, message := range unreadMessages {
		unreadMessageMap[message.MessageId] = message.IsIgnore
	}
	for _, reply := range replies {
		replyMessage := definition.ReplyMessage{
			ReplyId:   reply.ReplyId,
			PostId:    reply.PostId,
			CommentId: reply.CommentId,
			ReplyUId:  reply.UId,
			ReplyTxt:  reply.ReplyTxt,
			ReplyTime: reply.ReplyTime,
			IsUnread:  false,
		}
		if IsIgnore, exist := unreadMessageMap[reply.ReplyId]; exist { //是未读的回复
			if IsIgnore { // 用户忽略了消息,不加进返回
				continue
			}
			replyMessage.IsUnread = true
		}
		res = append(res, replyMessage)
	}
	return res
}

func PackageAtMessage(ats []definition.At, unreadMessages []definition.UnreadMessage) []definition.AtMessage {
	var res []definition.AtMessage
	unreadMessageMap := make(map[uint64]bool, len(unreadMessages)) // 转为 hash 集合，优化算法
	for _, message := range unreadMessages {
		unreadMessageMap[message.MessageId] = message.IsIgnore
	}
	for _, at := range ats {
		isUnread := false
		if IsIgnore, exist := unreadMessageMap[at.Id]; exist { //是未读的@
			if IsIgnore { // 用户忽略了消息,不加进返回
				continue
			}
			isUnread = true
		}
		atMessage := definition.AtMessage{
			AtId:       at.Id,
			UId:        at.UId,
			PostID:     at.PostID,
			MessageTxt: at.MessageTxt,
			Place:      at.Place,
			IsUnread:   isUnread,
		}
		res = append(res, atMessage)
	}
	return res
}

func PackageChatInfos(myUId uint64, chats []definition.Chat, unreadMessages []definition.UnreadMessage) map[uint64][]definition.ChatInfo {
	chatInfos := make(map[uint64][]definition.ChatInfo)
	unreadMessageMap := make(map[uint64]bool, len(unreadMessages)) // 转为 hash 集合，优化算法
	for _, message := range unreadMessages {
		unreadMessageMap[message.MessageId] = message.IsIgnore
	}
	for _, chat := range chats {
		var uId uint64     // 聊天对象的id
		var amISender bool // 我是否是发送人
		isUnread := false  // 是否已读

		if IsIgnore, exist := unreadMessageMap[chat.ChatId]; exist { //是未读的私聊
			if IsIgnore { // 用户忽略了消息,不加进返回
				continue
			}
			isUnread = true
		}

		if chat.SenderId == myUId {
			uId = chat.AddresseeId
			amISender = true
		} else {
			uId = chat.SenderId
			amISender = false
		}
		// 拼装聊天记录
		chatInfos[uId] = append(chatInfos[uId], definition.ChatInfo{
			AmISender: amISender,
			ChatTxt:   chat.ChatTxt,
			ImgId:     chat.ImgId,
			ChatTime:  chat.ChatTime,
			IsUnread:  isUnread,
		})
	}

	return chatInfos
}

func GetNewPrivacySetting(PrivacySetting byte, PostIsPrivate *bool, CommentAndReplyIsPrivate *bool, SavedPostIsPrivate *bool, SubscribedIsPrivate *bool) byte {
	if PostIsPrivate != nil {
		if *PostIsPrivate {
			PrivacySetting = PrivacySetting | 1
		} else {
			PrivacySetting = PrivacySetting & (255 - 1)
		}
	}

	if CommentAndReplyIsPrivate != nil {
		if *CommentAndReplyIsPrivate {
			PrivacySetting = PrivacySetting | 2
		} else {
			PrivacySetting = PrivacySetting & (255 - 2)
		}
	}

	if SavedPostIsPrivate != nil {
		if *SavedPostIsPrivate {
			PrivacySetting = PrivacySetting | 4
		} else {
			PrivacySetting = PrivacySetting & (255 - 4)
		}
	}

	if SubscribedIsPrivate != nil {
		if *SubscribedIsPrivate {
			PrivacySetting = PrivacySetting | 8
		} else {
			PrivacySetting = PrivacySetting & (255 - 8)
		}
	}

	return PrivacySetting
}

//PostIsPrivate 隐私设置判断 发帖记录1
func PostIsPrivate(PrivacySetting byte) bool {
	if PrivacySetting&1 > 0 {
		return true
	} else {
		return false
	}
}

//CommentAndReplyIsPrivate 隐私设置判断 评论和回复记录2
func CommentAndReplyIsPrivate(PrivacySetting byte) bool {
	if PrivacySetting&2 > 0 {
		return true
	} else {
		return false
	}
}

//SavedPostIsPrivate 隐私设置判断 收藏记录4
func SavedPostIsPrivate(PrivacySetting byte) bool {
	if PrivacySetting&4 > 0 {
		return true
	} else {
		return false
	}
}

//SubscribedIsPrivate 隐私设置判断 关注的人8
func SubscribedIsPrivate(PrivacySetting byte) bool {
	if PrivacySetting&8 > 0 {
		return true
	} else {
		return false
	}
}

//md5加密
func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
