package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

//这里FansID为本机已登录的user_id，BloggerId为将要进行关注操作的对方的user_id

type FansRequest struct {
	Token      string `form:"blogger_id" json:"token" binding:"required"`
	BloggerId  int64  `form:"blogger_id" json:"to_user_id" binding:"required"`
	ActionType int32  `form:"action_type" json:"action_type" binding:"required"`
}

type FansResponse struct {
	common.Response
}

/*
关注和取关操作
*/

func RelationAction(c *gin.Context) {
	var request FansRequest
	var response = &FansResponse{}
	if err := c.Bind(&request); err != nil {
		response.StatusCode = 1
		response.StatusMsg = "参数解析失败"
		c.JSON(400, response)
		log.Println("关注操作request参数绑定失败")
		return
	}
	FansIDAny, _ := c.Get("FansID")
	FansID, _ := strconv.ParseInt(fmt.Sprintf("%v", FansIDAny), 0, 64)
	err := service.FollowRelationAction(request.BloggerId, FansID, request.ActionType)
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
