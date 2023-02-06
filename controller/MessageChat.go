package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"log"

	"github.com/gin-gonic/gin"
)

type MessageChatRequest struct {
	Token    string `form:"token" json:"token" binding:"required"`
	ToUserID int64  `form:"to_user_id" json:"to_user_id" binding:"required"`
}

type MessageChatResponse struct {
	common.Response
	MessageList []common.Message `form:"message_list" json:"message_list"`
}

func MessageChat(c *gin.Context) {
	var request MessageChatRequest
	var response MessageChatResponse

	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}

	temp_user1, _ := c.Get("UserID")
	if temp_user1 != nil {
		user1 := temp_user1.(int64)
		user2 := request.ToUserID
		messageList, err := service.MessageChat(user1, user2)
		if err != nil {
			log.Println(err)
		}

		messageList_ := make([]common.Message, len(messageList))
		for i := range messageList {
			message := messageList[i]
			messageList_[i].Content = message.Content
			messageList_[i].CreateTime = message.CreateTime.Format("2006-01-02 15:04:05")
			messageList_[i].FromUserId = user1
			messageList_[i].ToUserID = user2
			messageList_[i].Id = message.ID
		}
		response.MessageList = messageList_
		response.StatusCode = 0
		response.StatusMsg = "success"
		c.JSON(200, response)
	} else {
		response.StatusCode = 1
		response.StatusMsg = "当前用户未登录！"
		c.JSON(200, response)
	}
}
