//主函数,启动http服务，连接数据库
package main

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/servicer"
	"net/http"
)

const (
	socket = ":15656"
)

func defaulttest(w http.ResponseWriter, r *http.Request) {
	servicer.Hearset(w, r)
	str := []byte("nihaonihao!")
	w.Write(str)
}

func main() {
	dataLayer.DB_open()
	defer dataLayer.DB_close()
	dataLayer.Redis_open()
	defer dataLayer.Redis_close()
	Mux1 := http.NewServeMux()
	Mux1.HandleFunc("/", defaulttest)
	Mux1.HandleFunc("/register", servicer.Register)
	Mux1.HandleFunc("/login", servicer.Login)
	Mux1.HandleFunc("/createpost", servicer.Createpost)
	Mux1.HandleFunc("/createcomment", servicer.Createcomment)
	Mux1.HandleFunc("/allpostid", servicer.Allpostid)

	Mux1.HandleFunc("/selectpostonid", servicer.Selectpostonid)
	Mux1.HandleFunc("/deletepostonid", servicer.Deletepostonid)
	Mux1.HandleFunc("/allcommentidonpostid", servicer.Allcommentidonpostid)
	Mux1.HandleFunc("/selectcommentonid", servicer.Selectcommentonid)
	Mux1.HandleFunc("/deletecommentonid", servicer.Deletecommentonid)
	Mux1.HandleFunc("/selectuseronid", servicer.Selectuseronid)
	Mux1.HandleFunc("/allposthot", servicer.Allposthot)
	Mux1.HandleFunc("/allpostidonuid", servicer.Allpostidonuid)
	Mux1.HandleFunc("/uploadimg", servicer.Uploadimg)
	Mux1.HandleFunc("/getimg/", servicer.Getimg)

	server := &http.Server{
		Addr:    socket,
		Handler: Mux1,
	}
	err := server.ListenAndServe()
	if err != nil { //无法监听端口
		dataLayer.Serverlog.Fatal("List", socket)
	}
}
