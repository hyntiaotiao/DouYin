package controller

import (
	"DouYIn/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FollowListRequest struct {
	UserId int64  `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

type FollowListResponse struct {
	Response
	UserList []UserVO `json:"user_list"`
}

// FollowList 获取登录用户关注的所有用户列表，两种情况下调用该接口：
//  1. 用户A查看自己的关注列表
//  2. 用户B查看用户A的关注列表
func FollowList(c *gin.Context) {
	targetId, _ := strconv.Atoi(c.Query("user_id"))
	curUserId, _ := c.Get("UserID")
	var response = &FollowListResponse{}
	followeeList := service.FindFolloweeList(int64(targetId))
	response.StatusCode = 0
	response.StatusMsg = "success"
	for _, user := range followeeList {
		var userVO UserVO
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

// FansList 获取登录用户关注的所有粉丝列表，两种情况下调用该接口：
//  1. 用户A查看自己的粉丝列表
//  2. 用户B查看用户A粉丝列表
func FollowerList(c *gin.Context) {
	targetId, _ := strconv.Atoi(c.Query("user_id"))
	curUserId, _ := c.Get("UserID")
	var response = &FollowListResponse{}
	followerList := service.FindFollowerList(int64(targetId))
	response.StatusCode = 0
	response.StatusMsg = "success"
	for _, user := range followerList {
		var userVO UserVO
		userVO.Id = user.ID
		userVO.Name = user.Username
		userVO.FollowCount = int64(user.FollowCount)
		userVO.FollowerCount = int64(user.FollowerCount)
		userVO.IsFollow = service.HasFollowed(user.ID, curUserId.(int64))
		response.UserList = append(response.UserList, userVO)
	}
	c.JSON(200, response)
}
