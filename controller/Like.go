package controller

import (
	"DouYIn/service"
	"github.com/gin-gonic/gin"
	"log"
)

type LikeRequest struct {
	Token      string `json:"token" binding:"required"`
	UserId     int64  `form:"user_id" json:"user_id" binding:"required"`
	VideoId    int64  `form:"video_id" json:"video_id" binding:"required"`
	ActionType int32  `form:"cancel" json:"cancel" binding:"required"`
}

type LikeResponse struct {
	Response
	UserID  int64 `json:"user_id" binding:"required"`
	VideoId int64 `json:"video_id" binding:"required"`
}

// FavoriteAction 点赞或者取消赞操作;
func FavoriteAction(c *gin.Context) {
	var request LikeRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("赞操作request参数绑定失败")
		return
	}
	var response = &LikeResponse{}
	likeInfo, err := service.FavouriteAction(request.UserId, request.VideoId, request.ActionType)
	if err != nil {
		log.Println("赞操作失败", err)
	}
	response.UserID = likeInfo.UserId
	response.VideoId = likeInfo.VideoId
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
