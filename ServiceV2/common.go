package ServiceV2

import (
	"code/Hahachitchat/definition"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 放一些通用的返回
func SetForbiddenResponse(c *gin.Context) {
	c.JSON(http.StatusForbidden, definition.CommonResponse{
		State:        definition.BadRequest,
		StateMessage: "服务器出错",
	})
}

func SetUnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, definition.CommonResponse{
		State:        definition.NoPermission,
		StateMessage: "未登录,无法操作",
	})
}

func SetGetUidErrorResponse(c *gin.Context) {
	c.JSON(http.StatusOK, definition.CommonResponse{
		State:        definition.ServerError,
		StateMessage: "session获取u_id失败",
	})
}

func SetParamErrorResponse(c *gin.Context) {
	c.JSON(http.StatusBadRequest, definition.CommonResponse{
		State:        definition.ServerError,
		StateMessage: "请求参数解析失败，参数不正确",
	})
}

func SetServerErrorResponse(c *gin.Context) {
	c.JSON(http.StatusOK, definition.CommonResponse{
		State:        definition.ServerError,
		StateMessage: "服务器出错",
	})
}

func SetDBErrorResponse(c *gin.Context) {
	c.JSON(http.StatusOK, definition.CommonResponse{
		State:        definition.ServerError,
		StateMessage: "数据库出错",
	})
}

func SetDBParamErrorResponse(c *gin.Context) {
	c.JSON(http.StatusOK, definition.CommonResponse{
		State:        definition.ServerError,
		StateMessage: "数据库存储的数据解析出错",
	})
}
