//提供存取图片的服务，未来可能移动到另一个服务器
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/nsqio/go-nsq"
)

const pdnsqdAddr = "127.0.0.1:4150"           //生产者只能指定一个nsqd
var cmnsqdsAddrs = []string{"127.0.0.1:4150"} //消费者能指定多个nsqd
var deleteImg_producer *nsq.Producer          //删除图片id的生产者

//创建图片文件夹
func init() {
	os.Mkdir("./imgdoc", os.ModePerm)
	cfg := nsq.NewConfig()
	deleteImg_producer, err = nsq.NewProducer(pdnsqdAddr, cfg) //初始化创建生产者 绑定nsqd
	if err != nil {
		log.Fatal(err)
	}
	go deleteImg_consum() //启动一个协程去订阅id删除图片
}

func uploadimg(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	if r.Method == "POST" { //注意内容不是json 而是"multipart/form-data"
		var goods struct { //响应体里的东西
			State int `json:"state"` //失败返回0,成功返回1,内容类型不正确返回2
		}
		err := r.ParseMultipartForm(8388608) //解析表单 即最大8M  8*1024*1024
		if err != nil {                      //解析表单出错
			Imglog.Println("uploadimg err解析表单出错", err)
			goods.State = 0
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		}
		ctlen := len("multipart/form-data")                                 //只需要前面一截判断内容类型,后面一截是标识
		if r.Header.Get("Content-Type")[0:ctlen] != "multipart/form-data" { //内容类型不正确
			Imglog.Println("uploadimg err内容类型不正确", err)
			goods.State = 2
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		}
		imgfilehear := r.MultipartForm.File["image"][0] //获取表单里的图片文件
		imgfile, err := imgfilehear.Open()              //把解码出的表单文件当成一个文件打开
		if err != nil {                                 //文件打开失败
			Imglog.Println("uploadimg err文件打开失败", err)
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
			Imglog.Println("uploadimg err复制保存失败", err)
			goods.State = 0
			goodsjson, _ := json.Marshal(goods)
			w.Write(goodsjson)
			return
		}
		object := r.MultipartForm.Value["object"][0]                          //对象的类型
		object_id, err := strconv.Atoi(r.MultipartForm.Value["object_id"][0]) //对象的id
		if err != nil {                                                       //object_id不能转化为int型
			Imglog.Println("uploadimg err object_id不能转化为int型", err)
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
		// 	// Imglog.Println("getimg", err)//可能是没有该图片
		// 	w.WriteHeader(404)
		// 	return
		// }
		//现在直接用这个发图片
		http.ServeFile(w, r, img_f)
		return
	}
}

//把要删除的图片id放进指定nsqd的topic
func deleteImg_produce(id string) {
	if id == "" { //空则不用发送 发送空的东西到消息队列会引发错误
		return
	}
	if err := deleteImg_producer.Publish("deleteImgid", []byte(id)); err != nil { //把图片id推送到队列
		Imglog.Println("deleteImg_produce error: " + err.Error()) //推送出错则记录
	}
}

//创建消费者 从订阅的消息队列中获取要删除的图片id
func deleteImg_consum() {
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("deleteImgid", "comsumer01", cfg) //创建消费者 管道comsumer01 绑定主题deleteImgid
	if err != nil {
		Imglog.Fatal(err)
	}
	//添加处理回调
	consumer.AddHandler(nsq.HandlerFunc(
		func(message *nsq.Message) error {
			id := string(message.Body)         //拿到id
			err := os.Remove("./imgdoc/" + id) //转化为路径并删除
			if err != nil {
				Imglog.Println("deleteImg_consumer Remove err:", err) //没有删除成功有两种情况：操作出错，图片不存在
				return nil                                            //默认为图片不存在,不用再返回消息队列
			}
			Imglog.Println("delete OK Img:", id) //删除成功
			return nil
		}))
	//用消费者 连接订阅的nsqd
	if err := consumer.ConnectToNSQDs(cmnsqdsAddrs); err != nil {
		Imglog.Fatal(err, " deleteImg_consumer err")
	}
	<-consumer.StopChan
}
