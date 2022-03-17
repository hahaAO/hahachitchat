package servicer

//
//import (
//	"code/Hahachitchat/dataLayer"
//	"net/http"
//)
//
//func defaulttest(w http.ResponseWriter, r *http.Request) {
//	Hearset(w, r)
//	str := []byte("nihaonihao!")
//	w.Write(str)
//}
//
//func StartService(port string) {
//	Mux1 := http.NewServeMux()
//	Mux1.HandleFunc("/", defaulttest)
//	Mux1.HandleFunc("/register", Register)
//	Mux1.HandleFunc("/login", Login)
//	Mux1.HandleFunc("/createpost", Createpost)
//	Mux1.HandleFunc("/createcomment", Createcomment)
//	Mux1.HandleFunc("/allpostid", Allpostid)
//
//	Mux1.HandleFunc("/selectpostonid", Selectpostonid)
//	Mux1.HandleFunc("/deletepostonid", Deletepostonid)
//	Mux1.HandleFunc("/allcommentidonpostid", Allcommentidonpostid)
//	Mux1.HandleFunc("/selectcommentonid", Selectcommentonid)
//	Mux1.HandleFunc("/deletecommentonid", Deletecommentonid)
//	Mux1.HandleFunc("/selectuseronid", Selectuseronid)
//	Mux1.HandleFunc("/allposthot", Allposthot)
//	Mux1.HandleFunc("/allpostidonuid", Allpostidonuid)
//	Mux1.HandleFunc("/uploadimg", Uploadimg)
//	Mux1.HandleFunc("/getimg/", Getimg)
//
//	server := &http.Server{
//		Addr:    port,
//		Handler: Mux1,
//	}
//	err := server.ListenAndServe()
//	if err != nil { //无法监听端口
//		dataLayer.Serverlog.Fatal("List", port)
//	}
//}
