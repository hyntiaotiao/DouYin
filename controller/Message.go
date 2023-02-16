package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"log"

	"github.com/gin-gonic/gin"
)

var messageIdSequence = int64(1)

type MessageActionRequest struct {
	Token      string `form:"token" json:"token" binding:"required"`
	ToUserId   int64  `form:"to_user_id" json:"to_user_id" binding:"required"`
	ActionType string `form:"action_type" json:"action_type" binding:"required"` //“1”发送消息
	Content    string `form:"content" json:"content" binding:"required"`
}

type MessageChatRequest struct {
	Token    string `form:"token" json:"token" binding:"required"`
	ToUserId int64  `form:"to_user_id" json:"to_user_id" binding:"required"`
}

type MessageChatResponse struct {
	common.Response
	MessageList []common.MessageVO `form:"message_list" json:"message_list"`
}

/*
消息操作
*/
func MessageAction(c *gin.Context) {
	var request MessageActionRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(400, common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败！"})
		log.Println("request参数绑定失败：", err)
		return
	}
	FromUserId, _ := c.Get("UserID")
	if request.ActionType != "1" {
		c.JSON(400, "action_type值不为1")
	}
	err := service.SendMessage(FromUserId.(int64), request.ToUserId, request.Content)
	if err != nil {
		c.JSON(500, &common.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(200, &common.Response{StatusCode: 0, StatusMsg: "success"})
}

/*
聊天记录
*/
func MessageChat(c *gin.Context) {
	var request MessageChatRequest
	var response MessageChatResponse

	if err := c.Bind(&request); err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败！"}
		c.JSON(400, response)
		log.Println("request参数绑定失败：", err)
		return
	}

	userIdAny, _ := c.Get("UserID")
	user1 := userIdAny.(int64)
	user2 := request.ToUserId
	messageList, err := service.MessageChat(user1, user2)
	if err != nil {
		log.Println("获取聊天列表失败：", err)
	}
	response.MessageList = messageList
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
