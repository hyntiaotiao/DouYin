package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"log"

	"github.com/gin-gonic/gin"
)

type VideoFeedRequest struct {
	LatestTime int64  `form:"latest_time" json:"latest_time" binding:"omitempty"`
	Token      string `form:"token" json:"token" binding:"omitempty"`
}

type VideoFeedResponse struct {
	common.Response
	VideoList []common.VideoVO `json:"video_list"`
	NextTime  int64            `json:"next_time"`
}

/*
视频流接口
*/
func Feed(c *gin.Context) {
	var request VideoFeedRequest
	var response = &VideoFeedResponse{}
	//参数绑定
	if err := c.Bind(&request); err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败！"}
		c.JSON(400, response)
		log.Println("request参数绑定失败：", err)
		return
	}

	userIdAny, _ := c.Get("UserID")
	var userId int64
	if userIdAny == nil {
		userId = -1
	} else {
		userId = userIdAny.(int64)
	}
	VideoList, NextTime, err := service.Feed(15, userId, request.LatestTime)
	if err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "获取视频列表失败！"}
		c.JSON(400, response)
		log.Println("获取视频列表失败：", err)
		return
	}
	response.VideoList = VideoList
	response.NextTime = NextTime
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
