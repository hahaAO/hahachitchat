//redis做缓存
//使用连接池
//目前做了comment表的缓存
package main

import (
	"fmt"
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
	redis_conn.Do("FLUSHALL") //初始化redis
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
func Redis_SelectCommentOnid(comment_id int) (int, Comment) {
	redis_conn := redisClient.Get()
	defer redis_conn.Close()
	var comment Comment
	args, err := redis.Values((redis_conn.Do(
		"HVALS", fmt.Sprintf("comment::%d", comment_id))))
	if err == redis.ErrNil || len(args) == 0 { //无此id0
		return 0, comment
	} else if err != nil { //其他情况2失败
		Redislog.Println("Redis_SelectCommentOnid err:", err)
		return 2, comment
	}
	comment.Comment_id, _ = strconv.Atoi(string(args[0].([]byte)))
	comment.Post_id, _ = strconv.Atoi(string(args[1].([]byte)))
	comment.U_id, _ = strconv.Atoi(string(args[2].([]byte)))
	comment.Comment_txt = string(args[3].([]byte))
	commentunix, _ := strconv.ParseInt(string(args[4].([]byte)), 10, 64)
	comment.Comment_time = time.Unix(0, commentunix) //精确到纳秒的时间戳
	comment.Img_id = string(args[5].([]byte))
	return 1, comment //查到有此id1成功
}

//把数据库的comment写入缓存 (int型，0失败，1则成功)
func Redis_CreateComment(comment Comment) int {
	redis_conn := redisClient.Get()
	defer redis_conn.Close()
	_, err := redis.String(
		redis_conn.Do(
			"HMSET", fmt.Sprintf("comment::%d", comment.Comment_id),
			"comment_id", comment.Comment_id,
			"post_id", comment.Post_id,
			"u_id", comment.U_id,
			"comment_txt", comment.Comment_txt,
			"comment_time", comment.Comment_time.UnixNano(), //精确到纳秒的时间戳
			"img_id", comment.Img_id))
	if err != nil { //其他情况3
		Redislog.Println("Redis_CreateComment err:", err)
		return 0
	}
	return 1 //插入成功
}

//根据comment_id删除comment (int型，0则失败，1则成功)
func Redis_DeleteCommentOnid(comment_id int) int {
	redis_conn := redisClient.Get()
	defer redis_conn.Close()
	_, err := redis_conn.Do(
		"DEL", fmt.Sprintf("comment::%d", comment_id))
	if err != nil { //其他情况0失败
		Redislog.Println("Redis_DeleteCommentOnid:", err)
		return 0
	}
	return 1 //删除成功
}
