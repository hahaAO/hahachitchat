//对应数据库中表的结构
//不使用外键
package main

import "time"

type User struct {
	U_id       int       `json:"u_id"`       //用户id,唯一
	U_name     string    `json:"u_name"`     //用户名,唯一
	U_password string    `json:"u_password"` //用户密码
	U_time     time.Time `json:"u_time"`     //用户注册时间
	U_nickname string    `json:"u_nickname"` //用户称昵
	Img_id     string    `json:"img_id"`     //图片唯一id用作用户头像
}

type Post struct {
	Post_id      int       `json:"post_id"`      ////帖子id，唯一
	U_id         int       `json:"u_id"`         //用户id
	Post_name    string    `json:"post_name"`    //帖子主题
	Post_txt     string    `json:"post_txt"`     //帖子内容
	Post_time    time.Time `json:"post_time"`    //帖子发布时间
	Post_txthtml string    `json:"post_txthtml"` //帖子内容的html
	Img_id       string    `json:"img_id"`       //图片唯一id用作镇楼图
}

type Comment struct {
	Comment_id   int       `json:"comment_id"`   //评论id，唯一
	Post_id      int       `json:"post_id"`      //帖子id
	U_id         int       `json:"u_id"`         //用户id
	Comment_txt  string    `json:"comment_txt"`  //评论内容
	Comment_time time.Time `json:"comment_time"` //评论时间
	Img_id       string    `json:"img_id"`       //图片唯一id用作评论图
}
