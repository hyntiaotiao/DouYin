package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"github.com/gin-gonic/gin"
	"log"
)

type VideoFeedRequest struct {
	LatestTime int64  `form:"latest_time" json:"latest_time" binding:"omitempty"`
	Token      string `form:"token" json:"token" binding:"omitempty"`
}

type VideoFeedResponse struct {
	common.Response
	VideoList []common.Video `json:"video_list"`
	NextTime  int64          `json:"next_time"`
}

func Feed(c *gin.Context) {
	var request VideoFeedRequest
	//参数绑定
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}

	var response = &VideoFeedResponse{}
	response.StatusCode = 0
	response.StatusMsg = "success"

	UserID, err := c.Get("UserID")
	if !err {
		UserID = -1
	}
	VideoList, NextTime, error := service.Feed(15, UserID, request.LatestTime)
	if error != nil {
		response.StatusCode = 1
		response.StatusMsg = "error"
		log.Println(error)
		c.JSON(200, response)
		return
	}
	response.VideoList = VideoList
	response.NextTime = NextTime
	c.JSON(200, response)
}
