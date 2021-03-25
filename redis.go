//redis做缓存
//目前做了comment表的缓存
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var redis_conn redis.Conn

func Redis_open() {
	redis_conn, _ = redis.Dial("tcp", ":6379",
		// redis.DialKeepAlive(300*time.Second),
		// redis.DialConnectTimeout(30*time.Second),
		redis.DialReadTimeout(20*time.Second),
		redis.DialWriteTimeout(20*time.Second))
	if err != nil {
		return
	}
	redis_conn.Do("FLUSHALL") //初始化redis
	ee, err := redis.String(redis_conn.Do("PING", "nihao"))
	if err != nil {
		Redislog.Println("Redis_open error:", err)
		return
	}
	Redislog.Println("Redis_open strat OK:", ee)
}

func Redis_close() {
	redis_conn.Close()
}

//根据comment_id获取comment (int型，0无此id，1则成功,2则失败)（comment）
func Redis_SelectCommentOnid(comment_id int) (int, Comment) {
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
	_, err := redis_conn.Do(
		"DEL", fmt.Sprintf("comment::%d", comment_id))
	if err != nil { //其他情况0失败
		Redislog.Println("Redis_DeleteCommentOnid:", err)
		return 0
	}
	return 1 //删除成功
}
