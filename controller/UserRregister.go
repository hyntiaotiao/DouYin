package controller

import (
	"DouYin/common"
	"DouYin/service"
	"DouYin/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"log"
)

type UserRegisterRequest struct {
	UserName string `form:"username" json:"username" validator:"required,min=6,max = 20"` //用户名长度最短6最长20
	Password string `form:"password" json:"password" validator:"required,min=6,max = 20"` //密码长度最短6最长20
}

type UserRegisterResponse struct {
	common.Response
	UserID int64  `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func Register(c *gin.Context) {
	var request UserRegisterRequest
	//接收参数
	if err := c.Bind(&request); err != nil {
		c.JSON(400, &common.Response{StatusCode: 1, StatusMsg: "用户名或密码格式错误"})
		log.Println("request参数绑定失败")
		return
	}
	var response = &UserRegisterResponse{}
	//传入的密码将在service层加密后再存入数据库
	UserID, err := service.Register(request.UserName, request.Password)
	if err != nil {
		c.JSON(400, &common.Response{StatusCode: 1, StatusMsg: "注册失败"})
		log.Println(err)
		return
	}
	//生成token
	token, err := utils.GenToken(UserID)
	if err != nil {
		c.JSON(400, &common.Response{StatusCode: 1, StatusMsg: "注册失败"})
		log.Println(err)
	}
	//注册成功
	response.UserID = UserID
	response.Token = token
	response.StatusCode = 0 //0成功 1失败
	response.StatusMsg = "success"
	c.JSON(200, response)
}
