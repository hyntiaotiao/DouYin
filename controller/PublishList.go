package controller

import (
	"DouYIn/repository"
	"DouYIn/service"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
)

type publishListRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

type publishListResponse struct {
	Response
	User        User               `json:"user" binding:"required"`
	PublishList []repository.Video `json:"publishList" binding:"required"`
}

func PublishList(c *gin.Context) {
	var request publishListRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}
	var response = &publishListResponse{}
	userID := request.UserID
	response.PublishList, _ = service.GetPublishList(userID)
	// response.User 请求user对象
	response.User.Id = userID
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
