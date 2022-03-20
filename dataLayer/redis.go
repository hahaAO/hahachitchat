//redis做缓存
//使用连接池
//目前做了comment表的缓存
//加上了session缓存
package dataLayer

import (
	"code/Hahachitchat/definition"
	"code/Hahachitchat/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var redisClient *redis.Pool

func Redis_open() {
	//初始化连接池
	redisClient = &redis.Pool{
		MaxIdle:     2,
		MaxActive:   10,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", 0)
			return c, nil
		},
	}
	//测试连接
	redis_conn := redisClient.Get()
	defer redis_conn.Close()
	//redis_conn.Do("FLUSHALL") //初始化redis
	ee, err := redis.String(redis_conn.Do("PING", "nihao"))
	if err != nil {
		Redislog.Println("Redis_open error:", err)
		return
	}
	Redislog.Println("Redis_open strat OK:", ee)
}

func Redis_close() {
	//关闭连接池
	Redislog.Println("Redis_close")
	redisClient.Close()
}

//根据comment_id获取comment (int型，0无此id，1则成功,2则失败)（comment）
func Redis_SelectCommentByid(comment_id uint64) (definition.DBcode, definition.Comment) {
	redis_conn := redisClient.Get()
	defer redis_conn.Close()
	var comment definition.Comment
	args, err := redis.Values(redis_conn.Do(
		"HMGET", fmt.Sprintf("comment::%d", comment_id),
		"comment_id",
		"post_id",
		"u_id",
		"comment_txt",
		"comment_time",
		"img_id",
		"someone_be_at",
	))
	if err == redis.ErrNil || len(args) == 0 {
		return definition.DB_NOEXIST, comment
	} else if err != nil { //其他情况2失败
		Redislog.Println("Redis_SelectCommentOnid err:", err)
		return definition.DB_ERROR, comment
	}
	for _, arg := range args { //无此参数
		if arg == nil {
			return definition.DB_NOEXIST, comment
		}
	}
	comment.CommentId, _ = strconv.ParseUint(string(args[0].([]byte)), 10, 64)
	comment.PostId, _ = strconv.ParseUint(string(args[1].([]byte)), 10, 64)
	comment.UId, _ = strconv.ParseUint(string(args[2].([]byte)), 10, 64)
	comment.CommentTxt = string(args[3].([]byte))
	commentunix, _ := strconv.ParseInt(string(args[4].([]byte)), 10, 64)
	comment.CommentTime = time.Unix(0, commentunix) //精确到纳秒的时间戳
	comment.ImgId = string(args[5].([]byte))
	comment.SomeoneBeAt = string(args[6].([]byte))
	return definition.DB_EXIST, comment //查到有此id 成功
}

//把数据库的comment写入缓存
func Redis_CreateComment(comment definition.Comment) definition.DBcode {
	redis_conn := redisClient.Get()
	defer redis_conn.Close()
	key := fmt.Sprintf("comment::%d", comment.CommentId)
	_, err := redis.String(
		redis_conn.Do(
			"HMSET", key,
			"comment_id", comment.CommentId,
			"post_id", comment.PostId,
			"u_id", comment.UId,
			"comment_txt", comment.CommentTxt,
			"comment_time", comment.CommentTime.Unix(),
			"img_id", comment.ImgId,
			"someone_be_at", comment.SomeoneBeAt,
		),
	)
	if err != nil { //其他情况
		Redislog.Println("Redis_CreateComment err:", err)
		return definition.DB_ERROR
	}
	_, err = redis.Int64(
		redis_conn.Do("EXPIRE", key, 18000),
	)
	if err != nil { //其他情况
		Redislog.Println("Redis_CreateComment err:", err)
		return definition.DB_ERROR
	}
	return definition.DB_SUCCESS // 插入成功
}

//根据comment_id删除comment
func Redis_DeleteCommentOnid(comment_id uint64) definition.DBcode {
	redis_conn := redisClient.Get()
	defer redis_conn.Close()
	_, err := redis_conn.Do(
		"DEL", fmt.Sprintf("comment::%d", comment_id))
	if err != nil { //其他情况 失败
		Redislog.Println("Redis_DeleteCommentOnid:", err)
		return definition.DB_ERROR
	}
	return definition.DB_SUCCESS //删除成功
}

//把初始化后的session存入Redis
func Redis_CreateSession(session definition.Session) definition.DBcode {
	redis_conn := redisClient.Get()
	defer redis_conn.Close()
	_, err := redis.String(
		redis_conn.Do(
			"SET",
			fmt.Sprintf("session::%s", session.Randid), //随机的id作为键
			session.Id, //真实的id作为值
			"EX",
			session.Expire, //过期时间
		))
	if err != nil { // 则失败
		Redislog.Println("Redis_CreateSession err:", err)
		return definition.DB_ERROR
	}
	return definition.DB_SUCCESS // 插入成功
}

//根据客户 session 的 ranid 查 id
func Redis_SelectSessionidByRandid(Randid string) (definition.DBcode, string) {
	redis_conn := redisClient.Get()
	defer redis_conn.Close()
	id, err := redis.String( // 把真的session拿出来对比
		redis_conn.Do(
			"GET",
			fmt.Sprintf("session::%s", Randid), //随机的id作为键
		))
	if err == redis.ErrNil { // 没有这个随机id
		return definition.DB_NOEXIST, ""
	} else if err != nil { // 其他情况
		Redislog.Println("Redis_SelectSession err:", err)
		return definition.DB_ERROR, ""
	}
	//查询成功
	return definition.DB_SUCCESS, id
}

//从cookie中提取session 正确返回 (对应session) 错误返回 (nil)
func GetSession(r *http.Request) *definition.Session {
	var session *definition.Session
	for _, cookienow := range r.Cookies() { //遍历所有cookie
		if cookienow.Name == "randid" { //找到的cookie("name"为"randid")
			session = utils.ParseToSession(*cookienow)                 //初始化对应session 设置session的randid
			scode, id := Redis_SelectSessionidByRandid(session.Randid) //验证session
			if scode == definition.DB_SUCCESS {                        //验证成功
				session.Id = id // 设置session的id
				return session
			} else { //验证失败
				return nil
			}
		}
	}
	//没有该cookie("name"为"randid")
	return nil
}
