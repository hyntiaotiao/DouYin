package controller

import (
	"DouYIn/service"
	"DouYIn/utils"
	"github.com/gin-gonic/gin"
	"log"
)

type UserLoginRequest struct {
	UserName string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Response
	UserID int64  `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func Login(c *gin.Context) {
	var request UserLoginRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}
	var response = &UserLoginResponse{}
	UserID, error := service.Login(request.UserName, request.Password)
	if error != nil {
		log.Println(error)
	}
	token, error := utils.GenToken(UserID)
	if error != nil {
		log.Println(error)
	}
	response.UserID = UserID
	response.Token = token
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
