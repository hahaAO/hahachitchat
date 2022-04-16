package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"code/Hahachitchat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ForbiddenMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := definition.ForbiddenConfig.ForbiddenIP[c.ClientIP()]
		if ok {
			SetForbiddenResponse(c)
			c.Abort()
			return
		}

		if c.Request.URL.Path == "/login" {
			return
		}

		session := utils.GetSession(c.Request)
		if session == nil { // cookie 中没有 session
			return
		} else {
			code, uId := dataLayer.Redis_SelectSessionidByRandid(*session)
			if code != definition.DB_SUCCESS { // session 错误
				return
			}
			if uId != nil {
				_, ok = definition.ForbiddenConfig.ForbiddenUser[*uId]
				if ok {
					SetForbiddenResponse(c)
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}

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
		session := utils.GetSession(c.Request)
		if session == nil { // cookie 中没有 session
			SetUnauthorizedResponse(c)
			c.Abort()
			return
		}
		code, uId := dataLayer.Redis_SelectSessionidByRandid(*session)
		if code != definition.DB_SUCCESS { // session 错误
			SetUnauthorizedResponse(c)
			c.Abort()
			return
		}
		// 有登录态
		c.Set("u_id", *uId) // 写入 u_id 后续可以获取
		c.Next()
	}
}

func SetSessionMiddleWare() gin.HandlerFunc { // 用户有登录态则写入，无也放行
	return func(c *gin.Context) {
		session := utils.GetSession(c.Request)
		if session == nil { // cookie 中没有 session
			c.Next()
			return
		}
		code, uId := dataLayer.Redis_SelectSessionidByRandid(*session)
		if code != definition.DB_SUCCESS { // session 错误
			c.Next()
			return
		}
		// 有登录态
		c.Set("u_id", *uId) // 写入 u_id 后续可以获取
		c.Next()
	}
}
