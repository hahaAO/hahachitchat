//放一些零碎的小工具
package utils

import (
	"code/Hahachitchat/db"
	"code/Hahachitchat/definition"
	"math/rand"
	"net/http"
	"os"
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

//验证cookie 正确返回 (对应session) 错误返回 (nil)
func VerifyCookie(r *http.Request) *definition.Session {
	var session *definition.Session
	for _, cookienow := range r.Cookies() { //遍历所有cookie
		if cookienow.Name == "randid" { //找到的cookie("name"为"randid")
			session = ParseToSession(*cookienow)    //初始化对应session 设置session的randid
			sint := db.Redis_SelectSession(session) //验证session 设置session的id
			if sint == 1 {                          //验证成功
				return session
			} else { //验证失败
				return nil
			}
		}
	}
	//没有该cookie("name"为"randid")
	return nil
}

func DeleteImg(id string) error {
	err := os.Remove("./imgdoc/" + id) //转化为路径并删除
	if err != nil {
		Imglog.Println("deleteImg_consumer Remove err:", err) //没有删除成功有两种情况：操作出错，图片不存在
		return nil                                            //默认为图片不存在,不用再返回消息队列
	}
	Imglog.Println("delete OK Img:", id) //删除成功
	return nil
}
