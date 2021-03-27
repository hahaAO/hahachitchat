//提供存取图片的服务，未来可能移动到另一个服务器
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

//创建图片文件夹
func init() {
	os.Mkdir("./imgdoc", os.ModePerm)
}

func uploadimg(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "POST" { //注意内容不是json 而是"multipart/form-data"
		var goods struct { //响应体里的东西
			State int `json:"state"` //失败返回0,成功返回1,内容类型不正确返回2
		}
		err := r.ParseMultipartForm(8388608) //解析表单 即最大8M  8*1024*1024
		if err != nil {                      //解析表单出错
			Serverlog.Println("uploadimg err解析表单出错", err)
			goods.State = 0
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		}
		ctlen := len("multipart/form-data")                                 //只需要前面一截判断内容类型,后面一截是标识
		if r.Header.Get("Content-Type")[0:ctlen] != "multipart/form-data" { //内容类型不正确
			Serverlog.Println("uploadimg err内容类型不正确", err)
			goods.State = 2
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		}
		imgfilehear := r.MultipartForm.File["image"][0] //获取表单里的图片文件
		imgfile, err := imgfilehear.Open()              //把解码出的表单文件当成一个文件打开
		if err != nil {                                 //文件打开失败
			Serverlog.Println("uploadimg err文件打开失败", err)
			goods.State = 0
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		}
		img_id := timeRandId() //图片唯一id
		filedocandname := fmt.Sprintf("./imgdoc/%s", img_id)
		saveFile, _ := os.Create(filedocandname) //创建文件
		_, err = io.Copy(saveFile, imgfile)      //复制保存
		if err != nil {                          //复制保存失败
			Serverlog.Println("uploadimg err复制保存失败", err)
			goods.State = 0
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		}
		object := r.MultipartForm.Value["object"][0]                          //对象的类型
		object_id, err := strconv.Atoi(r.MultipartForm.Value["object_id"][0]) //对象的id
		if err != nil {                                                       //object_id不能转化为int型
			Serverlog.Println("uploadimg err object_id不能转化为int型", err)
			goods.State = 0
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		}
		sint := UpdateObjectimgid(object, object_id, img_id)
		if sint == 1 {
			goods.State = 1
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		} else {
			goods.State = 0
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		}

	}
}

func getimg(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "image/*")
		img_id := r.FormValue("img_id") //获取图片id
		img_f := fmt.Sprintf("./imgdoc/%s", img_id)
		// content, err := ioutil.ReadFile(img_f) //读取图片到内存
		// if err != nil {
		// 	// Serverlog.Println("getimg", err)//可能是没有该图片
		// 	w.WriteHeader(404)
		// 	return
		// }
		//现在直接用这个发图片
		http.ServeFile(w, r, img_f)
		return
	}
}
