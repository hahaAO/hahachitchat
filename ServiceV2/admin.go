package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
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
		StateMessage: "查询封禁用户信息成功",
		BanUsers:     res,
	})
}

func GetBanIPs(c *gin.Context) {
	var res []uint64
	for idStr, _ := range definition.ForbiddenConfig.ForbiddenIP {
		id, _ := strconv.ParseUint(idStr, 10, 64)
		res = append(res, id)
	}

	c.JSON(http.StatusOK, definition.GetBanIPsResponse{
		State:        definition.Success,
		StateMessage: "查询封禁IP信息成功",
		BanIPList:     res,
	})
}

func PostStatisticsPieChart(c *gin.Context) {
	var req definition.PostStatisticsPieChartRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	startTime := time.UnixMilli(req.StartTimeSTP)
	endTime := time.UnixMilli(req.EndTimeSTP)

	fmt.Println(startTime)
	fmt.Println(endTime)

	if startTime.After(endTime) {
		SetParamErrorResponse(c)
		return
	}

	code, countSmallTalk, countStudyShare, countMarket := dataLayer.PostZoneCount(nil, startTime, endTime)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.PostStatisticsPieChartResponse{
			State:           definition.Success,
			StateMessage:    "查询分区统计成功",
			CountSmallTalk:  countSmallTalk,
			CountStudyShare: countStudyShare,
			CountMarket:     countMarket,
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)

	}
}

func PostStatisticsLineChart(c *gin.Context) {
	code, res := dataLayer.PostEverydayCount(nil)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.PostStatisticsLineChartResponse{
			State:          definition.Success,
			StateMessage:   "查询每日统计成功",
			PostCountByDay: res,
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)

	}
}

func PostStatisticsBarChart(c *gin.Context) {
	var req definition.PostStatisticsBarChartRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	startTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	code, res := dataLayer.PostEveryHourCount(nil, startTime)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.PostStatisticsBarChartResponse{
			State:           definition.Success,
			StateMessage:    "查询每小时统计成功",
			PostCountByHour: res,
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)

	}
}

func SetBanUser(c *gin.Context) {
	var req definition.SetBanUserIdsRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	banUser := make(map[string]struct{})
	for _, user := range req.BanUsers {
		banUser[strconv.FormatUint(user, 10)] = struct{}{}
	}
	var newForbiddenConfig definition.Forbidden
	newForbiddenConfig.ForbiddenIP = definition.ForbiddenConfig.ForbiddenIP
	newForbiddenConfig.ForbiddenUser = banUser

	jsonFile, err := os.OpenFile("./definition/forbidden.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer jsonFile.Close()
	if err != nil {
		dataLayer.Serverlog.Fatalln("jsonFile os.Open err: ", err)
		SetServerErrorResponse(c)
		return
	}

	if byte, err := json.Marshal(&newForbiddenConfig); err != nil {
		dataLayer.Serverlog.Fatalln("jsonFile Unmarshal err: ", err)
		SetServerErrorResponse(c)
	} else {
		if _, err := jsonFile.Write(byte); err != nil {
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

func SetBanIPs(c *gin.Context) {
	var req definition.SetBanIPsRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	banIPs := make(map[string]struct{})
	for _, ip := range req.BanIPList {
		banIPs[strconv.FormatUint(ip, 10)] = struct{}{}
	}
	var newForbiddenConfig definition.Forbidden
	newForbiddenConfig.ForbiddenUser = definition.ForbiddenConfig.ForbiddenUser
	newForbiddenConfig.ForbiddenIP = banIPs

	jsonFile, err := os.OpenFile("./definition/forbidden.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer jsonFile.Close()
	if err != nil {
		dataLayer.Serverlog.Fatalln("jsonFile os.Open err: ", err)
		SetServerErrorResponse(c)
		return
	}

	if byte, err := json.Marshal(&newForbiddenConfig); err != nil {
		dataLayer.Serverlog.Fatalln("jsonFile Unmarshal err: ", err)
		SetServerErrorResponse(c)
	} else {
		if _, err := jsonFile.Write(byte); err != nil {
			dataLayer.Serverlog.Fatalln("jsonFile.Write(byte) err: ", err)
			SetServerErrorResponse(c)
			return
		}
		c.JSON(http.StatusOK, definition.SetBanIPsResponse{
			State:        definition.Success,
			StateMessage: "设置封禁IP信息成功",
		})
	}
}

func SilenceUser(c *gin.Context) {
	var req definition.SilenceUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}
	code := dataLayer.UpdateSilenceUser(nil, req.UserId, req.DisableSendMsgTime)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.SilenceUserResponse{
			State:        definition.Success,
			StateMessage: "禁言用户成功",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}
