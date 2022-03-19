package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HearsetMiddleWare() gin.HandlerFunc { // 响应头设置，解决跨域问题
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Expose-Head", "Set-Cookie")

		method := c.Request.Method
		if method == "OPTIONS" { // 放行所有OPTIONS方法
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}

func AuthMiddleWare() gin.HandlerFunc { // 检查用户登录态
	return func(c *gin.Context) {
		session := dataLayer.GetSession(c.Request)
		if session == nil { //验证失败
			SetUnauthorizedResponse(c)
			c.Abort()
			return
		}
		// 有登录态
		c.Set("u_id", session.Id) // 写入 u_id 后续可以获取
		c.Next()
	}
}

func SetSessionMiddleWare() gin.HandlerFunc { // 用户有登录态则写入，无也放行
	return func(c *gin.Context) {
		session := dataLayer.GetSession(c.Request)
		if session == nil { // 无登录态
			c.Next()
			return
		}
		// 有登录态
		c.Set("u_id", session.Id) // 写入 u_id 后续可以获取
		c.Next()
	}
}
