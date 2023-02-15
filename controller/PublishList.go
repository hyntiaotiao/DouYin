package controller

import (
	"DouYin/common"
	"DouYin/service"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
)

type PublishListRequest struct {
	UserID int64  `form:"user_id" json:"user_id" validator:"required,gt=0"`
	Token  string `form:"token" json:"token" validator:"required" `
}

type PublishListResponse struct {
	common.Response
	PublishList []common.Video `json:"video_list" binding:"required"`
}

func PublishList(c *gin.Context) {
	var request PublishListRequest
	var response = &PublishListResponse{}
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}

	userID := request.UserID
	videoList, _ := service.PublishList(userID)

	log.Println(videoList)

	// response
	response.StatusCode = 0
	response.StatusMsg = "success"
	response.PublishList = videoList
	c.JSON(200, response)
}
