package controller

import (
	"DouYin/common"
	"DouYin/service"
	"DouYin/utils"
	"github.com/gin-gonic/gin"
	"log"
)

type UserLoginRequest struct {
	UserName string `form:"username" json:"username" validator:"required,min=6,max = 20"`
	Password string `form:"password" json:"password" validator:"required,min=6,max = 20"`
}

type UserLoginResponse struct {
	common.Response
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
	UserID, err := service.Login(request.UserName, request.Password)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := utils.GenToken(UserID)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response.UserID = UserID
	response.Token = token
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
