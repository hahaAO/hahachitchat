//主函数,启动http服务，连接数据库
package main

import (
	// "fmt"

	"net/http"
)

const (
	socket = ":15656"
)

func defaulttest(w http.ResponseWriter, r *http.Request) {
	hearset(w, r)
	str := []byte("nihaonihao!")
	w.Write(str)
}

func main() {
	//DB_open()
	//defer DB_close()
	//Redis_open()
	//defer Redis_close()
	Mux1 := http.NewServeMux()
	Mux1.HandleFunc("/", defaulttest)
	Mux1.HandleFunc("/register", register)
	Mux1.HandleFunc("/login", login)
	Mux1.HandleFunc("/createpost", createpost)
	Mux1.HandleFunc("/createcomment", createcomment)
	Mux1.HandleFunc("/allpostid", allpostid)

	Mux1.HandleFunc("/selectpostonid", selectpostonid)
	Mux1.HandleFunc("/deletepostonid", deletepostonid)
	Mux1.HandleFunc("/allcommentidonpostid", allcommentidonpostid)
	Mux1.HandleFunc("/selectcommentonid", selectcommentonid)
	Mux1.HandleFunc("/deletecommentonid", deletecommentonid)
	Mux1.HandleFunc("/selectuseronid", selectuseronid)
	Mux1.HandleFunc("/allposthot", allposthot)
	Mux1.HandleFunc("/allpostidonuid", allpostidonuid)
	Mux1.HandleFunc("/uploadimg", uploadimg)
	Mux1.HandleFunc("/getimg/", getimg)

	server := &http.Server{
		Addr:    socket,
		Handler: Mux1,
	}
	err := server.ListenAndServe()
	if err != nil { //无法监听端口
		Serverlog.Fatal("List", socket)
	}
}
