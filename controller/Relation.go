package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowListRequest struct {
	UserId int64  `json:"user_id" validator:"required,gt=0"`
	Token  string `json:"token" validator:"required"`
}

type FollowListResponse struct {
	common.Response
	UserList []common.UserVO `json:"user_list"`
}

type FriendUser struct {
	Message string `json:"message"`  //和该好友的最新聊天消息
	MsgType int64  `json:"msg_type"` // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
	*common.UserVO
}

type FriendListResponse struct {
	common.Response
	UserList []FriendUser `json:"user_list"`
}

type RelationActionRequest struct {
	Token    string `form:"token" json:"token" validator:"required"`
	ToUserID int64  `form:"to_user_id" json:"to_user_id" validator:"required,gt=0"`
	// 1关注 2取消
	ActionType int `form:"action_type" json:"action_type" validator:"required,gt=0,lt=3"`
}

/*
用户关注列表：获取登录用户关注的所有用户列表，两种情况下调用该接口：
1. 用户A查看自己的关注列表
2. 用户B查看用户A的关注列表
*/
func FollowList(c *gin.Context) {
	targetId, _ := strconv.Atoi(c.Query("user_id"))
	curUserId, _ := c.Get("UserID")
	var response = &FollowListResponse{}
	followeeList := service.FindFolloweeList(int64(targetId))
	response.StatusCode = 0
	response.StatusMsg = "success"
	for _, user := range followeeList {
		var userVO common.UserVO
		userVO.Id = user.ID
		userVO.Name = user.Username
		userVO.FollowCount = int64(user.FollowCount)
		userVO.FollowerCount = int64(user.FollowerCount)
		if int64(targetId) == curUserId.(int64) {
			// 用户A查看自己的关注列表
			userVO.IsFollow = true
		} else {
			// 用户B查看用户A的关注列表
			userVO.IsFollow = service.HasFollowed(user.ID, curUserId.(int64))
		}
		response.UserList = append(response.UserList, userVO)
	}
	c.JSON(200, response)
}

/*
用户粉丝列表：获取登录用户关注的所有粉丝列表，两种情况下调用该接口：
1. 用户A查看自己的粉丝列表
2. 用户B查看用户A粉丝列表
*/
func FollowerList(c *gin.Context) {
	targetId, _ := strconv.Atoi(c.Query("user_id"))
	curUserId, _ := c.Get("UserID")
	var response = &FollowListResponse{}
	followerList := service.FindFollowerList(int64(targetId))
	response.StatusCode = 0
	response.StatusMsg = "success"
	for _, user := range followerList {
		var userVO common.UserVO
		userVO.Id = user.ID
		userVO.Name = user.Username
		userVO.FollowCount = int64(user.FollowCount)
		userVO.FollowerCount = int64(user.FollowerCount)
		userVO.IsFollow = service.HasFollowed(user.ID, curUserId.(int64))
		response.UserList = append(response.UserList, userVO)
	}
	c.JSON(200, response)
}

/*
用户好友列表：获取登录用户关注的所有粉丝列表，一种情况下调用该接口：
1. 用户A查看自己的好友列表
*/
func FriendList(c *gin.Context) {
	targetId, _ := strconv.Atoi(c.Query("user_id"))
	curUserId, _ := c.Get("UserID")
	var response = &FriendListResponse{}
	if int64(targetId) != curUserId.(int64) {
		c.JSON(404, &FriendListResponse{common.Response{StatusCode: 1, StatusMsg: "用户id有误"}, nil})
		return
	}

	friendList := service.FindFriendList(curUserId.(int64))
	response.StatusCode = 0
	response.StatusMsg = "success"

	for _, user := range friendList {
		var userVO common.UserVO
		userVO.Id = user.ID
		userVO.Name = user.Username
		userVO.FollowCount = int64(user.FollowCount)
		userVO.FollowerCount = int64(user.FollowerCount)
		userVO.IsFollow = service.HasFollowed(user.ID, curUserId.(int64))

		var friendUser FriendUser
		friendUser.UserVO = &userVO

		//这里写死了“最新消息”和“消息类型”，等“消息”接口写好之后直接替换即可
		friendUser.Message = "这是最新的消息！"
		friendUser.MsgType = 0
		response.UserList = append(response.UserList, friendUser)
	}
	c.JSON(200, response)
}

/*
关注操作
*/
func RelationAction(c *gin.Context) {
	var request RelationActionRequest
	var response = &common.Response{}
	//参数绑定
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}
	//获取粉丝ID（即当前登录用户）
	FansID, _ := c.Get("UserID")
	log.Println(FansID.(int64), request.ToUserID, request.ActionType)
	err := service.FollowRelationAction(request.ToUserID, FansID.(int64), request.ActionType)
	if err != nil {
		log.Println("关注/取关操作失败", err)
		response.StatusCode = 1
		response.StatusMsg = "关注/取关操作失败"
		c.JSON(400, response)
		return
	}
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}
