package controller

import (
	"DouYin/common"
	"DouYin/service"
	"github.com/gin-gonic/gin"
	"log"
)

var messageIdSequence = int64(1)

type MessageActionRequest struct {
	token    string `form:"token" json:"token" validator:"required"`
	ToUserID int64  `form:"to_user_id" json:"to_user_id" validator:"required,gt=0"`
	//“1”发送消息
	ActionType string `form:"action_type" json:"action_type" validator:"required,gte=1,lte1"`
	Content    string `form:"content" json:"content" validator:"required"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	var request MessageActionRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}
	FromUserID, _ := c.Get("UserID")
	if request.ActionType != "1" {
		c.JSON(400, "action_type值不为1")
	}
	err := service.SendMessage(FromUserID.(int64), request.ToUserID, request.Content)
	if err != nil {
		c.JSON(500, &common.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(200, &common.Response{StatusCode: 0, StatusMsg: "success"})
}
