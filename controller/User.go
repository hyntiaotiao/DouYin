package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"DouYIn/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserInfoRequest struct {
	UserId int64  `form:"user_id" json:"user_id" validator:"required,gt=0"`
	Token  string `form:"token" json:"token" validator:"required"`
}

type UserInfoResponse struct {
	common.Response
	User common.UserVO
}

type UserLoginRequest struct {
	Username string `form:"username" json:"username" validator:"required,min=6,max = 20"`
	Password string `form:"password" json:"password" validator:"required,min=6,max = 20"`
}

type UserLoginResponse struct {
	common.Response
	UserId int64  `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

type UserRegisterRequest struct {
	Username string `form:"username" json:"username" validator:"required,min=6,max = 20"` //用户名长度最短6最长20
	Password string `form:"password" json:"password" validator:"required,min=6,max = 20"` //密码长度最短6最长20
}

type UserRegisterResponse struct {
	common.Response
	UserId int64  `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

/*
用户信息
*/
func UserInfo(c *gin.Context) {
	var request UserInfoRequest
	var response = &UserInfoResponse{}
	if err := c.Bind(&request); err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败！"}
		c.JSON(400, response)
		log.Println("request参数绑定失败：", err)
		return
	}
	targetId := request.UserId
	user, err := service.GetByID(targetId)
	if err != nil {
		log.Println(err)
		c.JSON(400, common.Response{StatusCode: 1, StatusMsg: "该用户不存在！"})
		return
	}
	userVO := common.UserVO{}
	userVO.Id = user.ID
	userVO.Name = user.Username
	userVO.FollowCount = int64(user.FollowCount)
	userVO.FollowerCount = int64(user.FollowerCount)
	// 设置is_follow属性：
	//   1.获取 发起当前请求的用户的信息
	// 		a.如果是未登录的用户，is_follow的值应该为false；
	// 		b.如果是已经登录的用户，is_follow的值根据fans表中的数据决定
	userIdAny, _ := c.Get("UserID")
	curUserId, _ := strconv.ParseInt(fmt.Sprintf("%v", userIdAny), 0, 64)
	if curUserId == targetId {
		userVO.IsFollow = false
	} else {
		userVO.IsFollow = service.HasFollowed(targetId, curUserId)
	}

	response.StatusCode = 0
	response.StatusMsg = "success"
	response.User = userVO
	c.JSON(200, response)
}

/*
用户登录接口
*/
func Login(c *gin.Context) {
	var request UserLoginRequest
	var response = &UserLoginResponse{}
	if err := c.Bind(&request); err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败"}
		c.JSON(400, response)
		log.Println("request参数绑定失败")
		return
	}
	userId, err := service.Login(request.Username, request.Password)
	if err != nil {
		log.Println(err)
	}
	token, err := utils.GenToken(userId)
	if err != nil {
		log.Println(err)
	}
	response.UserId = userId
	response.Token = token
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}

/*
用户注册接口
*/
func Register(c *gin.Context) {
	var request UserRegisterRequest
	var response = &UserRegisterResponse{}
	//接收参数
	if err := c.Bind(&request); err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败"}
		c.JSON(400, response)
		log.Println("request参数绑定失败")
		return
	}
	userId, err := service.Register(request.Username, request.Password)
	if err != nil {
		log.Println("注册失败：", err)
		c.JSON(400, &common.Response{StatusCode: 1, StatusMsg: "注册失败！"})
		return
	}
	token, error := utils.GenToken(userId)
	if error != nil {
		panic(error)
	}
	response.UserId = userId
	response.Token = token
	response.StatusCode = 0 //0成功 1失败
	response.StatusMsg = "success"
	c.JSON(200, response)
}
