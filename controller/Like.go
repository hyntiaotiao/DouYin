package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type LikeRequest struct {
	Token      string `form:"video_id" json:"token" binding:"required"`
	VideoID    int64  `form:"video_id" json:"video_id" binding:"required"`
	ActionType int32  `form:"action_type" json:"action_type" binding:"required"`
}

type LikeResponse struct {
	common.Response
}

// FavoriteAction 点赞或者取消赞操作;
func FavoriteAction(c *gin.Context) {
	var request LikeRequest
	var response = &LikeResponse{}
	if err := c.Bind(&request); err != nil {
		response.StatusCode = 1
		response.StatusMsg = "参数解析失败"
		c.JSON(400, response)
		log.Println("赞操作request参数绑定失败")
		return
	}
	UserIDAny, _ := c.Get("UserID")
	UserID, _ := strconv.ParseInt(fmt.Sprintf("%v", UserIDAny), 0, 64)
	err := service.FavouriteAction(UserID, request.VideoID, request.ActionType)
	if err != nil {
		log.Println("赞操作失败", err)
		response.StatusCode = 1
		response.StatusMsg = "赞操作失败"
		c.JSON(400, response)
		return
	}
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
