//操作数据库的函数
package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "vgdvgd111"
	dbname   = "hahadb"
)

var DB *sql.DB
var err error

//连接一个数据库，并测试连接
func DB_open() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: DB,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	gormDB.AutoMigrate(&User{}, &Post{}, &Comment{})

	DBlog.Printf("Successfully connect to postgres %s!\n", dbname)
}

//关闭数据库连接
func DB_close() {
	DB.Close()
	DBlog.Printf("Successfully off to postgres %s!\n", dbname)
}

//根据uid返回user （未注册0 已注册1 其他情况3）（User）
func SelectUserOnid(a int) (int, User) {
	var user User
	err := DB.QueryRow(`SELECT * FROM "user"
	WHERE "u_id"=$1`, a).Scan(
		&user.U_id,
		&user.U_name,
		&user.U_password,
		&user.U_time,
		&user.U_nickname,
		&user.Img_id)

	if err == sql.ErrNoRows {
		return 0, user //未注册
	} else if err != nil {
		DBlog.Println("SelectUserOnid:", err)
		return 3, user //其他问题
	}
	return 1, user //已注册
}

//根据name获取user （未注册0 已注册1 其他情况3）（User）
func SelectUserOnname(name string) (int, User) {
	var user User
	err := DB.QueryRow(`SELECT * FROM "user" WHERE u_name = $1`, name).Scan(
		&user.U_id,
		&user.U_name,
		&user.U_password,
		&user.U_time,
		&user.U_nickname,
		&user.Img_id)
	if err == sql.ErrNoRows {
		return 0, user //未注册
	} else if err != nil {
		DBlog.Println("SelectUserOnname err:", err)
		return 3, user //其他问题
	}
	return 1, user //已注册
}

//根据post id获取post （无此id0 查到有此id1 其他情况3）（Post）
func SelectPostOnid(post_id int) (int, Post) {
	var post Post
	err := DB.QueryRow(`SELECT * FROM "post" WHERE post_id = $1`, post_id).Scan(
		&post.Post_id,
		&post.U_id,
		&post.Post_name,
		&post.Post_txt,
		&post.Post_time,
		&post.Post_txthtml,
		&post.Img_id)
	if err == sql.ErrNoRows {
		return 0, post //无此id0
	} else if err != nil {
		DBlog.Println("SelectPostOnid err:", err)
		return 3, post //其他情况3
	}
	return 1, post //查到有此id1
}

//加了读redis缓存的功能		根据comment_id获取comment (int型，0无此id，1则成功,2则失败)（comment）
func SelectCommentOnid(comment_id int) (int, Comment) {
	sint, scomment := Redis_SelectCommentOnid(comment_id) //先读redis缓存
	if sint == 1 {                                        //redis中有此comment
		if scomment.Post_id == 0 { //Redis中为空值
			return 0, scomment
		} else { //Redis中存在
			return 1, scomment
		}
	} else { //redis中无此id	或	redis出错	要到postgres中查
		var comment Comment
		err := DB.QueryRow(`SELECT * FROM "comment" WHERE comment_id = $1`, comment_id).Scan(
			&comment.Comment_id,
			&comment.Post_id,
			&comment.U_id,
			&comment.Comment_txt,
			&comment.Comment_time,
			&comment.Img_id)
		if err == sql.ErrNoRows { //无此id0
			comment.Comment_id = comment_id
			comment.Post_id = 0
			Redis_CreateComment(comment) //把数据库的comment 空值 写入redis
			return 0, comment
		} else if err != nil { //其他情况3
			DBlog.Println("SelectCommentOnid err:", err)
			return 3, comment
		}
		//查到有此id1
		Redis_CreateComment(comment) //把数据库的comment写入redis
		return 1, comment

	}
}

//获取所有帖子的post (int，0则没有帖子，1则成功，2则有其他问题)（ []int）
func AllSelectPost() (int, []Post) {
	var posts []Post
	rows, err := DB.Query(`SELECT * FROM "post"`)
	defer rows.Close()
	if err == sql.ErrNoRows { //没有帖子
		return 0, posts
	} else if err != nil {
		return 2, posts
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(
			&post.Post_id,
			&post.U_id,
			&post.Post_name,
			&post.Post_txt,
			&post.Post_time,
			&post.Post_txthtml,
			&post.Img_id)
		if err != nil {
			DBlog.Println("AllSelectPost err1:", err)
			return 2, posts
		}
		posts = append(posts, post)
	}
	if len(posts) == 0 { //没有帖子
		return 0, posts
	}
	return 1, posts
}

//根据post_id获取所有comment_id(int型，0则没有评论，1则成功，2则有其他问题)（ []int）
func AllCommentidOnpostid(post_id int) (int, []int) {
	var commentids []int
	rows, err := DB.Query(`SELECT "comment_id" FROM "comment" WHERE "post_id"=$1`, post_id)
	defer rows.Close()
	if err == sql.ErrNoRows { //没有评论
		return 0, commentids
	} else if err != nil {
		DBlog.Println("AllCommentidOnpostid err1:", err)
		return 2, commentids
	}
	for rows.Next() {
		var commentid int
		err = rows.Scan(&commentid)
		if err != nil {
			DBlog.Println("AllCommentidOnpostid err2:", err)
			return 2, commentids
		}
		commentids = append(commentids, commentid)
	}
	if len(commentids) == 0 { //没有评论
		return 0, commentids
	}
	return 1, commentids
}

//根据name password Unickname插入user （注册失败0 注册成功1）（User）
func CreateUser(Uname string, Upassword string, Unickname string) (int, User) {
	var user User
	_, err := DB.Exec(`INSERT INTO "user" ("u_name","u_password","u_nickname")
    VALUES ($1,$2,$3)`, Uname, Upassword, Unickname)
	if err != nil {
		DBlog.Println("CreateUser err1:", err)
		return 0, user //其他问题,注册失败
	}
	sint, suser := SelectUserOnname(Uname)
	if sint == 1 {
		return 1, suser //注册成功
	}
	return 0, user //其他问题,注册失败
}

//根据uid post_name post_txt post_txthtml插入post（0则失败，1则成功，2则无此人id，3则有其他问题） （post_id）
func CreatePost(u_id int, post_name string, post_txt string, post_txthtml string) (int, int) {
	var post_id int
	sint, _ := SelectUserOnid(u_id)
	if sint == 0 { //2则无此人id
		return 2, post_id
	} else if sint == 1 { //有此人id
		err := DB.QueryRow(`INSERT INTO "post" ("u_id","post_name","post_txt","post_txthtml")
        VALUES ($1,$2,$3,$4) RETURNING post_id`, u_id, post_name, post_txt, post_txthtml).Scan(&post_id)
		if err != nil {
			DBlog.Println("CreatePost err1:", err)
			return 3, post_id //其他问题,插入失败
		}
		//插入成功
		return 1, post_id //1则成功
	} else { //其他问题,插入失败
		return 3, post_id
	}

}

//根据post_id u_id comment_txt插入comment(int型，0则失败，1则成功，2则无此人id，3则无帖子id，4则有其他问题), （comment_id）
func CreateComment(post_id int, u_id int, comment_txt string) (int, int) {
	var comment_id int
	sint, _ := SelectUserOnid(u_id) //查u_id
	if sint == 0 {                  //则无此人id
		return 2, comment_id
	} else if sint == 1 { //有此人id
		sint, _ := SelectPostOnid(post_id) //查post_id
		if sint == 0 {                     //2则无帖子id
			return 3, comment_id
		} else if sint == 1 { //1有此帖子id
			err := DB.QueryRow(`INSERT INTO "comment" ("post_id","u_id","comment_txt")
            VALUES ($1,$2,$3) RETURNING comment_id`, post_id, u_id, comment_txt).Scan(&comment_id)
			if err != nil {
				DBlog.Println("CreateComment err1:", err)
				return 0, comment_id //其他问题,插入失败
			}
			//插入成功
			return 1, comment_id //成功
		} else { //3其他问题,插入失败
			return 0, comment_id
		}
	} else { //3其他问题,插入失败
		return 0, comment_id
	}
}

//根据name password查询 （未注册0 已注册密码正确1 已注册密码错误2 其他情况3）  （User）
func SelectUsernamepassword(name string, password string) (int, User) {
	var user User
	sint, user := SelectUserOnname(name)
	if sint == 0 { //未注册0
		return 0, user
	} else if sint == 1 { //已注册
		if user.U_password == password { //密码正确1
			return 1, user
		} else { //密码错误2
			return 2, user
		}
	} else { //其他问题3
		return 3, user
	}
}

//根据post_id 删除帖子及帖子里的评论 (int型，1则成功，2则有其他问题)
func DeletePostOnid(post_id int) int {
	var Img_id string                                                                              //要删除的图片id
	rows, err := DB.Query(`DELETE FROM "comment" WHERE "post_id" = $1 RETURNING img_id;`, post_id) //删除帖子里的评论，顺带读出图片id
	if err != nil {                                                                                //有其他问题
		DBlog.Println("DeletePostOnid err1:", err)
		return 2
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&Img_id) //读出图片id
		if err != nil {
			DBlog.Println("DeletePostOnid err2:", err)
			return 2
		}
		deleteImg_produce(Img_id) //把要删除的图片id发到消息队列
	}
	err = DB.QueryRow(`DELETE FROM "post" WHERE "post_id" = $1 RETURNING img_id;`, post_id).Scan(&Img_id) //删除帖子，顺带读出图片id
	if err != nil {                                                                                       //有其他问题
		DBlog.Println("DeletePostOnid err3:", err)
		return 2
	}
	deleteImg_produce(Img_id) //把要删除的图片id发到消息队列
	DBlog.Printf("DeletePostOnid:	post_id %d 删除成功\n", post_id)
	return 1
}

//redis缓存中的也删掉	根据comment_id 删除评论 (int型，1则成功，2则有其他问题)
func DeleteCommentOnid(comment_id int) int {
	var Img_id string //图片也要删除
	err := DB.QueryRow(`DELETE FROM "comment" WHERE "comment_id" = $1 RETURNING img_id;`,
		comment_id).Scan(&Img_id)
	if err != nil { //有其他问题
		DBlog.Println("DeleteCommentOnid err1:", err)
		return 2
	} else { //删除成功
		Redis_DeleteCommentOnid(comment_id) //redis缓存中的也删掉
		deleteImg_produce(Img_id)           //把要删除的图片id发到消息队列
		return 1
	}
}

//根据用户u_id 获取属于该用户的所有帖子postids (int型，0则没有帖子，1则成功，2则有其他问题), []int）
func SelectPostidByuid(u_id int) (int, []int) {
	var postids []int
	rows, err := DB.Query(`SELECT "post_id" FROM "post" WHERE "u_id"=$1`, u_id)
	defer rows.Close()
	if err == sql.ErrNoRows { //没有帖子
		return 0, postids
	} else if err != nil { //2则有其他问题
		return 2, postids
	}
	for rows.Next() {
		var postid int
		err = rows.Scan(&postid)
		if err != nil {
			DBlog.Println("SelectPostidByuid:", err)
			return 2, postids //2则有其他问题
		}
		postids = append(postids, postid)
	}
	if len(postids) == 0 { //没有帖子
		return 0, postids
	}
	//1则成功
	return 1, postids
}

//根据对象类型，对象id，图片id 设置对应对象的图片id:即头像or镇楼图or评论图	(int型 0则失败,1则成功)
func UpdateObjectimgid(object string, object_id int, img_id string) int {
	objecthead := object //sql填充字段的第二个，因为user表的字段head是"u"而不是"user",需要转换
	if object == "user" {
		objecthead = "u"
	}
	sql := fmt.Sprintf(`UPDATE "%s" SET img_id = $1 WHERE	"%s_id" = $2;`, object, objecthead)
	_, err := DB.Exec(sql, img_id, object_id)
	if err != nil { //object不正确会报错，但是object_id不存在则不会报错
		DBlog.Println("UpdateObjectimgid err", err)
		return 0
	}
	return 1
}
