package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	code, forbiddenUser := dataLayer.SelectForbiddenUser(nil)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.GetBanUserIdsResponse{
			State:              definition.Success,
			StateMessage:       "查询封禁用户信息成功",
			BanUserIdAndReason: forbiddenUser,
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetBanIPs(c *gin.Context) {
	code, forbiddenIp := dataLayer.SelectForbiddenIp(nil)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.GetBanIPsResponse{
			State:          definition.Success,
			StateMessage:   "查询封禁IP信息成功",
			BanIPAndReason: forbiddenIp,
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
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

func AddBanUser(c *gin.Context) {
	var req definition.AddBanUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.CreateBanUser(nil, req.BanUserId, req.Reason)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.AddBanUserResponse{
			State:        definition.Success,
			StateMessage: "封禁成功",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func CancelBanUser(c *gin.Context) {
	var req definition.CancelBanUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.DeleteBanUser(nil, req.BanUserId)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.CancelBanUserResponse{
			State:        definition.Success,
			StateMessage: "解除封禁成功",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func AddBanIP(c *gin.Context) {
	var req definition.AddBanIPRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.CreateBanIP(nil, req.BanIP, req.Reason)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.AddBanIPResponse{
			State:        definition.Success,
			StateMessage: "封禁成功",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func CancelBanIp(c *gin.Context) {
	var req definition.CancelBanIpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.DeleteBanIP(nil, req.BanIP)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.CancelBanIpResponse{
			State:        definition.Success,
			StateMessage: "解除封禁成功",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
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

func SetApprovalUser(c *gin.Context) {
	var req definition.SetApprovalUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil || req.NeedApproval == nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.UpdateApprovalUser(nil, req.UserId, *req.NeedApproval)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.SetApprovalUserResponse{
			State:        definition.Success,
			StateMessage: "设置审批用户成功",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func SetTopPost(c *gin.Context) {
	var req definition.SetTopPostRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	var code definition.DBcode
	if req.IsTop == nil {
		SetParamErrorResponse(c)
		return
	} else if *req.IsTop == true {
		code = dataLayer.CreateTopPost(nil, req.PostId, req.Describe)
	} else if *req.IsTop == false {
		code = dataLayer.DeleteTopPost(nil, req.PostId)
	}

	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.SetTopPostResponse{
			State:        definition.Success,
			StateMessage: "设置精品帖子成功",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}

func GetNeedApprovalPost(c *gin.Context) {
	code, posts := dataLayer.SelectApprovalPost(nil)

	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.GetNeedApprovalPostResponse{
			State:         definition.Success,
			StateMessage:  "获取待审核帖子成功",
			ApprovalPosts: posts,
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}

}

func ApprovalPost(c *gin.Context) {
	var req definition.ApprovalPostRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SetParamErrorResponse(c)
		return
	}

	code := dataLayer.ApprovalPost(nil, req.ApprovalPostId)
	switch code {
	case definition.DB_SUCCESS:
		c.JSON(http.StatusOK, definition.ApprovalPostResponse{
			State:        definition.Success,
			StateMessage: "审批通过",
		})
	case definition.DB_ERROR:
		SetDBErrorResponse(c)
	default:
		SetServerErrorResponse(c)
	}
}
