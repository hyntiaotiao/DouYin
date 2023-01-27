package controller

import (
	"DouYIn/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
)

type publishListRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

type publishListResponse struct {
	Response
	PublishList service.PublishList `json:"publishList" binding:"required"`
}

func PublishList(c *gin.Context) {
	var response = &publishListResponse{}
	response.PublishList = service.getPublishList()
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
