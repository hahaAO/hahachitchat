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
	StateCode    int    `json:"state_code"`
	StateMessage string `json:"state_message"`
}

type ClientLog struct {
	IP              string          `json:"ip"`
	Uid             *uint64         `json:"u_id"`
	OperateName     string          `json:"operate_name"`
	HttpStatusCode  int             `json:"http_status_code"`
	ResponseMessage ResponseMessage `json:"response_message"`
	OperateTime     string          `json:"operate_time"`
}

var URLMapOperateName = map[string]string{
	"/register":       "用户注册",
	"/login":          "用户登录",
	"/reset-password": "重置密码",
	"/create-post":    "用户发贴",
	"/create-comment": "发表评论",
	"/create-reply":   "发表回复",
	"/delete-post":    "删除贴子",
	"/delete-comment": "删除评论",
	"/delete-reply":   "删除回复",
	"/uploadimg":      "更换头像",
}
