// 放一些零碎的小工具
// 不依赖数据层的增删改查
package utils

import (
	"code/Hahachitchat/definition"
	"math/rand"
	"net/http"
	"strconv"
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
func CreateSession(id int) *definition.Session {
	return &definition.Session{
		Id:     strconv.FormatInt(int64(id), 10), //真实id
		Randid: TimeRandId(),                     //随机生成
		Expire: int(3600 * 48),                   //默认两天,
	}
}

//把初始化后的session转换为cookie
func ParseToCookie(session *definition.Session) http.Cookie {
	return http.Cookie{
		Name:     "randid",
		Value:    session.Randid,
		HttpOnly: true,
		MaxAge:   session.Expire,
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
