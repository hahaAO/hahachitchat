//处理前端的请求以及调用数据库函数对操作数据库进行操作
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

//统一设置响应头的跨域和内容类型
func hearset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func register(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "POST" {
		var receiver struct { //接收 请求体里的东西
			U_name     string `json:"u_name"`
			U_password string `json:"u_password"`
			U_nickname string `json:"u_nickname"`
		}
		var goods struct { //响应体里的东西
			State int `json:"state"`
			U_id  int `json:"u_id"`
		}
		body := make([]byte, r.ContentLength)
		r.Body.Read(body) // 调用 Read 方法读取请求实体并将返回内容存放到上面创建的字节切片
		err := json.Unmarshal(body, &receiver)
		if err != nil {
			Serverlog.Println("json err:", err)
			Serverlog.Println("body:", string(body)) //用于查看请求体里的东西
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
		sint, suser := SelectUserOnname(receiver.U_name)
		if sint == 1 { //已注册
			goods.State = 0
			goods.U_id = suser.U_id
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println(goods.U_id,"已注册")
			return
		} else if sint == 0 { //未注册
			cint, cuser := CreateUser(receiver.U_name, receiver.U_password, receiver.U_nickname)
			if cint == 1 { //注册成功
				goods.State = 1
				goods.U_id = cuser.U_id
				goods_byte, _ := json.Marshal(goods)
				w.Write(goods_byte)
				Serverlog.Println(goods.U_id, "注册成功")
				return
			} else if cint == 0 { //注册失败
				goods.State = 2
				goods_byte, _ := json.Marshal(goods)
				w.Write(goods_byte)
				// Serverlog.Println(goods.U_id,"未注册,注册失败")
				return
			}
		} else if sint == 3 {
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println(goods.U_id,"sql有问题，注册失败")
			return
		}
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "GET" {
		var receiver struct { //接收 请求参数里的东西
			U_name     string `json:"u_name"`
			U_password string `json:"u_password"`
		}
		var goods struct { //响应体里的东西
			State      int    `json:"state"`
			U_id       int    `json:"u_id"`
			U_nickname string `json:"u_nickname"`
		}
		receiver.U_name = r.FormValue("u_name")
		receiver.U_password = r.FormValue("u_password")
		sint, suser := SelectUsernamepassword(receiver.U_name, receiver.U_password)
		if sint == 1 { //1已注册密码准确
			goods.State = 1
			goods.U_id = suser.U_id
			goods.U_nickname = suser.U_nickname
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println(goods.U_id, "登陆成功")
			return
		} else if sint == 2 { //2已注册密码错误
			goods.State = 2
			goods.U_id = suser.U_id
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("密码错误")
			return
		} else if sint == 0 { //0未注册
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("未注册")
			return
		} else if sint == 3 { //3则有其他问题)
			goods.State = 3
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("其他问题")
			return
		}
	}
}
func createpost(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "POST" {
		var receiver struct { //接收 请求体里的东西
			U_id         int    `json:"u_id"`
			Post_name    string `json:"post_name"`
			Post_txt     string `json:"post_txt"`
			Post_txthtml string `json:"post_txthtml"` //帖子内容的html
		}
		var goods struct { //响应体里的东西
			State   int `json:"state"`
			Post_id int `json:"post_id"`
		}
		body := make([]byte, r.ContentLength)
		r.Body.Read(body) // 调用 Read 方法读取请求实体并将返回内容存放到上面创建的字节切片
		err := json.Unmarshal(body, &receiver)
		if err != nil {
			Serverlog.Println("errjson", err)
			Serverlog.Println("body:", string(body)) //用于查看请求体里的东西
			goods.State = 3                          //3则有其他问题
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
		cint, cpost_id := CreatePost(receiver.U_id, receiver.Post_name, receiver.Post_txt, receiver.Post_txthtml)
		if cint == 0 { //0则失败
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else if cint == 1 { //1则成功
			goods.State = 1 //创建成功
			goods.Post_id = cpost_id
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else if cint == 2 { //2则无此人id
			goods.State = 2 //则无此人id
			goods.Post_id = cpost_id
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else { //3则有其他问题
			goods.State = 3 //有其他问题
			goods.Post_id = cpost_id
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
	}
}
func createcomment(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "POST" {
		var receiver struct { //接收 请求体里的东西
			Post_id     int    `json:"post_id"`
			U_id        int    `json:"u_id"`
			Comment_txt string `json:"comment_txt"`
		}
		var goods struct { //响应体里的东西
			State      int `json:"state"` //(int型，0则失败，1则成功，2则无此人id，3则无帖子id，4则有其他问题)
			Comment_id int `json:"comment_id"`
		}
		body := make([]byte, r.ContentLength)
		r.Body.Read(body) // 调用 Read 方法读取请求实体并将返回内容存放到上面创建的字节切片
		// Serverlog.Println("body:",string(body))//用于查看请求体里的东西
		err := json.Unmarshal(body, &receiver)
		if err != nil {
			Serverlog.Println("errjson", err)
			goods.State = 4 //4则有其他问题
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
		cint, ccomid := CreateComment(receiver.Post_id, receiver.U_id, receiver.Comment_txt)
		if cint == 1 { //成功
			goods.State = 1
			goods.Comment_id = ccomid
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println("创建评论成功")
			return
		} else if cint == 2 { //无此人id
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println("无此人id")
			return
		} else if cint == 3 { //无此帖子id
			goods.State = 3
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println("无此帖子id")
			return
		} else { //失败
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println("插入失败")
			return
		}
	}
}
func allselectpostid(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "GET" {
		var goods struct { //响应体里的东西
			State   int   `json:"state"` //(int型，0则失败，1则成功，2则有其他问题)
			Postids []int `json:"postids"`
		}
		aint, aposts := AllSelectPost()
		if aint == 1 { //查询成功
			goods.State = 1
			along := len(aposts)
			for i := 0; i < along; i++ {
				apost := aposts[i]
				goods.Postids = append(goods.Postids, apost.Post_id)
			}
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("查询成功")
			return
		} else if aint == 0 { //没有帖子
			goods.State = 1
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("没有帖子")
			return
		} else {
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("查询失败")
			return
		}
	}
}
func selectpostonid(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "GET" {
		var goods struct { //响应体里的东西
			State        int       `json:"state"` //(int型，0则失败，1则成功，2则有其他问题),
			U_id         int       `json:"u_id"`
			Post_name    string    `json:"post_name"`
			Post_txt     string    `json:"post_txt"`
			Post_time    time.Time `json:"post_time"`
			Post_txthtml string    `json:"post_txthtml"`
			Img_id       string    `json:"img_id"`
		}
		Post_id, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil { //参数不能转为int
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println("参数不能转为int", err)
			return
		}
		sint, spost := SelectPostOnid(Post_id)
		if sint == 1 { //查到了
			goods.State = 1
			goods.U_id = spost.U_id
			goods.Post_name = spost.Post_name
			goods.Post_txt = spost.Post_txt
			goods.Post_time = spost.Post_time
			goods.Post_txthtml = spost.Post_txthtml
			goods.Img_id = spost.Img_id
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("查到了")
			return
		} else if sint == 0 {
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("无此帖子id")
			return
		} else {
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("有其他问题")
			return
		}
	}
}
func deletepostonid(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "POST" {
		var receiver struct { //接收 请求体里的东西
			Post_id int `json:"post_id"`
		}
		var goods struct { //响应体里的东西
			State int `json:"state"` //(int型，0则失败没有该帖子，1则成功，2则有其他问题)
		}
		body := make([]byte, r.ContentLength)
		r.Body.Read(body) // 调用 Read 方法读取请求实体并将返回内容存放到上面创建的字节切片
		// Serverlog.Println("body:",string(body))//用于查看请求体里的东西
		err := json.Unmarshal(body, &receiver)
		if err != nil {
			Serverlog.Println("errjson", err)
			goods.State = 2 //2则有其他问题
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
		dint := DeletePostOnid(receiver.Post_id)
		if dint == 1 { //成功
			goods.State = 1
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println(receiver.Post_id, "删除成功")
			return
		} else { //删除失败
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			// Serverlog.Println("删除失败")
			return
		}
	}
}
func allcommentidonpostid(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "GET" {
		var goods struct { //响应体里的东西
			State      int   `json:"state"` //(int型，0则失败没有该帖子，1则成功，2则有其他问题)
			Commentids []int `json:"commentids"`
		}
		Post_id, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil { //参数不能转为int
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println("参数不能转为int", err)
			return
		}
		aint, acommentids := AllCommentidOnpostid(Post_id)
		if aint == 0 { //没有评论
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else if aint == 1 { //成功
			goods.State = 1
			goods.Commentids = acommentids
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else { //出问题
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
	}
}
func selectcommentonid(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "GET" {
		var goods struct { //响应体里的东西
			State        int       `json:"state"` //(int型，0则无此评论id，1则成功，2则有其他问题)
			U_id         int       `json:"u_id"`
			Comment_txt  string    `json:"comment_txt"`
			Comment_time time.Time `json:"comment_time"`
			Img_id       string    `json:"img_id"`
		}
		Comment_id, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil { //参数不能转为int
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println("参数不能转为int", err)
			return
		}
		sint, scomment := SelectCommentOnid(Comment_id)
		if sint == 0 { //无此评论id
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else if sint == 1 { //有此评论id
			goods.State = 1
			goods.U_id = scomment.U_id
			goods.Comment_txt = scomment.Comment_txt
			goods.Comment_time = scomment.Comment_time
			goods.Img_id = scomment.Img_id
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return

		} else { //其他问题
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
	}
}
func deletecommentonid(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "POST" {
		var receiver struct { //接收 请求体里的东西
			Comment_id int `json:"comment_id"`
		}
		var goods struct { //响应体里的东西
			State int `json:"state"` //(int型，0则失败没有该评论，1则成功，2则有其他问题)
		}
		body := make([]byte, r.ContentLength)
		r.Body.Read(body) // 调用 Read 方法读取请求实体并将返回内容存放到上面创建的字节切片
		// Serverlog.Println("body:",string(body))//用于查看请求体里的东西
		err := json.Unmarshal(body, &receiver)
		if err != nil {
			Serverlog.Println("errjson", err)
			goods.State = 2 //2则有其他问题
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
		dint := DeleteCommentOnid(receiver.Comment_id)
		if dint == 1 {
			goods.State = 1 //1则成功
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else { //2则有其他问题
			goods.State = 2 //2则有其他问题
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
	}
}
func selectuseronid(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "GET" {
		var goods struct { //响应体里的东西
			State      int       `json:"state"` //(int型，0则没有此人，1则成功，2则有其他问题)
			U_nickname string    `json:"u_nickname"`
			U_time     time.Time `json:"u_time"`
			Img_id     string    `json:"img_id"`
		}
		U_id, err := strconv.Atoi(r.FormValue("u_id"))
		if err != nil { //参数不能转为int
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println("参数不能转为int", err)
			return
		}
		sint, suser := SelectUserOnid(U_id)
		if sint == 1 { //已注册
			goods.State = 1 //1则成功
			goods.U_nickname = suser.U_nickname
			goods.U_time = suser.U_time
			goods.Img_id = suser.Img_id
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else if sint == 0 { //无此人
			goods.State = 0 //无此人
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else {
			goods.State = 2 //其他问题
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
	}
}
func allposthot(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "GET" {

		var goods struct { //响应体里的东西
			State    int             `json:"state"` //(int型，0则失败，1则成功)
			Hot_desc []Post_idandhot `json:"hot_desc"`
		}

		goods.Hot_desc, err = Allposthot()

		if err != nil {
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else {
			plong := len(goods.Hot_desc)
			for i := 0; i < plong-1; i++ {
				for j := 0; j < plong-i-1; j++ {
					if goods.Hot_desc[j].Post_hot < goods.Hot_desc[j+1].Post_hot {
						goods.Hot_desc[j], goods.Hot_desc[j+1] = goods.Hot_desc[j+1], goods.Hot_desc[j]
					}
				}
			}
			goods.State = 1
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
	}
}
func selectpostidbyuid(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "GET" {
		var goods struct { //响应体里的东西
			State   int   `json:"state"` //(int型，0则没有帖子，1则成功，2则有其他问题),
			Postids []int `json:"postids"`
		}
		u_id, err := strconv.Atoi(r.FormValue("u_id")) //接受参数
		if err != nil {                                //参数不能转为int
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			Serverlog.Println("参数不能转为int", err)
			return
		}
		sint, spostids := SelectPostidByuid(u_id)
		if sint == 1 {
			goods.State = 1
			goods.Postids = spostids
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else if sint == 0 {
			goods.State = 0
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		} else {
			goods.State = 2
			goods_byte, _ := json.Marshal(goods)
			w.Write(goods_byte)
			return
		}
	}
}
