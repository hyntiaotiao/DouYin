package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavoriteActionRequest struct {
	Token      string `form:"token" json:"token" validator:"required "`
	VideoId    int64  `form:"video_id" json:"video_id" validator:"required,gt=0"`
	ActionType int32  `form:"action_type" json:"action_type" validator:"required,gte=1,lte=2"` //1点咱2取消
}

type FavoriteActionResponse struct {
	common.Response
}

type FavoriteListRequest struct {
	Token  string `form:"token" json:"token" validator:"required"`
	UserId int64  `form:"user_id" json:"user_id" validator:"required,gt=0"`
}

type FavoriteListResponse struct {
	common.Response
	VideoList []common.VideoVO `json:"video_list"`
}

/*
赞操作
*/
func FavoriteAction(c *gin.Context) {
	var request FavoriteActionRequest
	var response = &FavoriteActionResponse{}
	if err := c.Bind(&request); err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败！"}
		c.JSON(400, response)
		log.Println("request参数绑定失败：", err)
		return
	}
	userIdAny, _ := c.Get("UserID")
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", userIdAny), 0, 64)
	err := service.FavoriteAction(userId, request.VideoId, request.ActionType)
	if err != nil {
		log.Println("赞操作失败：", err)
		response.StatusCode = 1
		response.StatusMsg = "赞操作失败"
		c.JSON(400, response)
		return
	}
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}

/*
喜欢列表
*/
func FavoriteList(c *gin.Context) {
	var request FavoriteListRequest
	var response = &FavoriteListResponse{}
	if err := c.Bind(&request); err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败！"}
		c.JSON(400, response)
		log.Println("request参数绑定失败：", err)
		return
	}
	userIdAny, _ := c.Get("UserID")
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", userIdAny), 0, 64)
	VideoList, err := service.FavoriteList(userId, request.UserId)
	if err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "赞操作失败"}
		c.JSON(400, response)
		log.Println("赞操作失败：", err)
		return
	}
	response.VideoList = VideoList
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
