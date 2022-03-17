// 放一些零碎的小工具
// 不依赖数据层的增删改查
package utils

import (
	"code/Hahachitchat/definition"
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
		Expire: int(3600 * 48),             //默认两天,
	}
}

//把初始化后的session转换为cookie
func ParseToCookie(session *definition.Session) http.Cookie {
	return http.Cookie{
		Name:     "randid",
		Value:    session.Randid,
		HttpOnly: true,
		MaxAge:   session.Expire,
		Secure:   false,
	}
}

//把cookie转换为session（需要验证）
func ParseToSession(cookie http.Cookie) *definition.Session {
	return &definition.Session{
		//id未知 验证成功再设置
		Randid: cookie.Value,
		//过期时间无所谓
	}
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

//PostIsPrivate 隐私设置判断 发帖记录1
func PostIsPrivate(PrivacySetting byte) bool {
	if PrivacySetting&1 > 0 {
		return true
	} else {
		return false
	}
}

//CommentIsPrivate 隐私设置判断 回复记录2
func CommentIsPrivate(PrivacySetting byte) bool {
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
