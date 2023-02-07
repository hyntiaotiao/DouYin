package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserInfoRequest struct {
	UserId int64  `form:"user_id" json:"user_id"`
	Token  string `form:"token" json:"token"`
}

type UserInfoResponse struct {
	common.Response
	User common.UserVO
}

func UserInfo(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	var response = &UserInfoResponse{}
	user, err := service.GetUserByID(int64(userId))
	if err != nil {
		log.Println(err)
		response := common.Response{StatusCode: 1, StatusMsg: "用户不存在"}
		c.JSON(404, response)
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
	curUserId, _ := c.Get("UserID")
	if curUserId == userId {
		userVO.IsFollow = false
	} else {
		userVO.IsFollow = service.HasFollowed(int64(userId), curUserId.(int64))
	}

	response.StatusCode = 0
	response.StatusMsg = "success"
	response.User = userVO
	c.JSON(200, response)
}
