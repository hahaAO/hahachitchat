//对应数据库中表的结构
//不使用外键
package main

import "time"

type User struct {
	U_id       int       `gorm:"primaryKey" column:"u_id"` //用户id,唯一
	U_name     string    `column:"u_name"`                 //用户名,唯一
	U_password string    `column:"u_password"`             //用户密码
	U_time     time.Time `column:"u_time"`                 //用户注册时间
	U_nickname string    `column:"u_nickname"`             //用户称昵
	Img_id     string    `column:"img_id"`                 //图片唯一id用作用户头像
}

func (User) TableName() string {
	return "user"
}

type Post struct {
	Post_id      int       `gorm:"primaryKey" column:"post_id"` //帖子id，唯一
	U_id         int       `column:"u_id"`                      //用户id
	Post_name    string    `column:"post_name"`                 //帖子主题
	Post_txt     string    `column:"post_txt"`                  //帖子内容
	Post_time    time.Time `column:"post_time"`                 //帖子发布时间
	Post_txthtml string    `column:"post_txthtml"`              //帖子内容的html
	Img_id       string    `column:"img_id"`                    //图片唯一id用作镇楼图
}

func (Post) TableName() string {
	return "post"
}

type Comment struct {
	Comment_id   int       `gorm:"primaryKey" column:"comment_id"` //评论id，唯一
	Post_id      int       `column:"post_id"`                      //帖子id
	U_id         int       `column:"u_id"`                         //用户id
	Comment_txt  string    `column:"comment_txt"`                  //评论内容
	Comment_time time.Time `column:"comment_time"`                 //评论时间
	Img_id       string    `column:"img_id"`                       //图片唯一id用作评论图
}

func (Comment) TableName() string {
	return "comment"
}

type Session struct {
	Id     string //用户id
	Randid string //随机的唯一id
	Expire int    //存活时间单位为秒
}
