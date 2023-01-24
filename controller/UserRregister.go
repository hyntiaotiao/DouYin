package controller

import (
	"DouYIn/service"
	"DouYIn/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"log"
)

type UserRegisterRequest struct {
	UserName string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserRegisterResponse struct {
	Response
	UserID int64  `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func Register(c *gin.Context) {
	var request UserRegisterRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}
	var response = &UserRegisterResponse{}
	UserID, error := service.Register(request.UserName, request.Password)
	if error != nil {
		log.Println(error)
	}
	token, error := utils.GenToken(UserID)
	if error != nil {
		panic(error)
	}
	response.UserID = UserID
	response.Token = token
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
