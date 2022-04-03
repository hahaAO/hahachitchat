package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func GetAllUser(c *gin.Context) {
	code, user := dataLayer.AllUserMessage(nil)
	switch code {
	case definition.DB_EXIST:
		c.JSON(http.StatusOK, definition.AllUserResponse{
			State:        definition.Success,
			StateMessage: "查询所有用户信息成功",
			Users:        user,
		})
	case definition.DB_NOEXIST:
		c.JSON(http.StatusOK, definition.AllUserResponse{
			State:        definition.Success,
			StateMessage: "查询成功，没有用户",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetBanUser(c *gin.Context) {
	var res []uint64
	for idStr, _ := range definition.ForbiddenConfig.ForbiddenUser {
		id, _ := strconv.ParseUint(idStr, 10, 64)
		res = append(res, id)
	}

	c.JSON(http.StatusOK, definition.GetBanUserIdsResponse{
		State:        definition.Success,
		StateMessage: "查询所有用户信息成功",
		BanUsers:     res,
	})
}

func SetBanUser(c *gin.Context) {
	var req definition.SetBanUserIdsRequest
	err := c.ShouldBindJSON(&req);if err != nil {
		SetParamErrorResponse(c)
		return
	}

	banUser := make(map[string]struct{})
	for _, user := range req.BanUsers {
		banUser[strconv.FormatUint(user, 10)] = struct{}{}
	}
	var newForbiddenConfig definition.Forbidden
	newForbiddenConfig.ForbiddenIP=definition.ForbiddenConfig.ForbiddenIP
	newForbiddenConfig.ForbiddenUser=banUser

	jsonFile, err := os.OpenFile("./definition/forbidden.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer jsonFile.Close()
	if err != nil {
		dataLayer.Serverlog.Fatalln("jsonFile os.Open err: ", err)
		SetServerErrorResponse(c)
		return
	}

	if byte ,err := json.Marshal(&newForbiddenConfig); err != nil {
		dataLayer.Serverlog.Fatalln("jsonFile Unmarshal err: ", err)
		SetServerErrorResponse(c)
	}else {
		if _,err:=jsonFile.Write(byte);err!=nil{
			dataLayer.Serverlog.Fatalln("jsonFile.Write(byte) err: ", err)
			SetServerErrorResponse(c)
			return
		}
		c.JSON(http.StatusOK, definition.SetBanUserIdsResponse{
			State:        definition.Success,
			StateMessage: "设置封禁用户信息",
		})
	}
}
