package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type LikeListRequest struct {
	Token  string `form:"token" json:"token" binding:"required"`
	UserID int64  `form:"user_id" json:"user_id" binding:"required"`
}

type LikeListResponse struct {
	common.Response
	VideoList []common.Video `json:"video_list"`
}

// FavoriteList LikeList 点赞列表
func FavoriteList(c *gin.Context) {
	var request LikeListRequest
	var response = &LikeListResponse{}
	if err := c.Bind(&request); err != nil {
		response.StatusCode = 1
		response.StatusMsg = "参数解析失败"
		c.JSON(400, response)
		log.Println("request参数绑定失败")
		return
	}
	UserIDAny, _ := c.Get("UserID")
	UserID, _ := strconv.ParseInt(fmt.Sprintf("%v", UserIDAny), 0, 64)
	VideoList, error := service.LikeList(UserID, request.UserID)
	if error != nil {
		log.Println("赞操作失败", error)
		response.StatusCode = 1
		response.StatusMsg = "赞操作失败"
		c.JSON(400, response)
		return
	}
	response.VideoList = VideoList
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
