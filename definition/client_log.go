package definition

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

var ClientLogChan chan string
var WsLogConnChan chan struct{}

type CustomResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type ResponseMessage struct {
	State        int    `json:"业务状态码"`
	StateMessage string `json:"响应消息"`
}

type ClientLog struct {
	IP              string          `json:"IP地址"`
	Uid             *uint64         `json:"用户ID(可空)"`
	OperateName     string          `json:"用户操作"`
	HttpStatusCode  int             `json:"HTTP响应"`
	ResponseMessage ResponseMessage `json:"操作结果"`
}

var URLMapOperateName = map[string]string{
	"/register": "用户注册",
	"/login": "用户登录",
	"/create-post": "用户发贴",
	"/create-comment": "发表评论",
	"/create-reply": "发表回复",
	"/delete-post": "删除贴子",
	"/delete-comment": "删除评论",
	"/delete-reply": "删除回复",
	"/uploadimg": "更换头像",
}
