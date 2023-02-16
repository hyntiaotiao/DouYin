package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type CommentActionRequest struct {
	VideoId     int64  `form:"video_id" json:"video_id"`
	Token       string `form:"token" json:"token"`
	ActionType  int32  `form:"action_type" json:"action_type"`
	CommentId   int64  `form:"comment_id" json:"comment_id,omitempty"`
	CommentText string `form:"comment_text" json:"comment_text,omitempty"`
}

type CommentActionResponse struct {
	common.Response
	Comment common.CommentVO `json:"comment,omitempty"`
}

type CommentListRequest struct {
	VideoId int64  `form:"video_id" json:"video_id"`
	Token   string `form:"token" json:"token"`
}

type CommentListResponse struct {
	common.Response
	CommentList []common.CommentVO `json:"comment_list" binding:"required"`
}

/*
评论操作
*/
func CommentAction(c *gin.Context) {
	client := resty.New()
	var request CommentActionRequest
	var response = &CommentActionResponse{}
	if err := c.Bind(&request); err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败！"}
		c.JSON(400, response)
		log.Println("request参数绑定失败：", err)
		return
	}
	userIdAny, _ := c.Get("UserID")
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", userIdAny), 0, 64)
	comment, err := service.CommentAction(request.ActionType, request.CommentId, request.VideoId, userId, request.CommentText)
	if err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "获取评论失败！"}
		c.JSON(400, response)
		log.Println("获取评论失败：", err)
		return
	}
	log.Println("操作的评论：", comment)
	if request.ActionType == 1 {
		user := common.UserVO{}
		userId_ := strconv.FormatInt(userId, 10)
		_, err := client.R().
			SetResult(&user).
			SetQueryParams(map[string]string{
				"user_id": userId_,
				"token":   request.Token,
			}).
			Get("http://localhost:8080/douyin/user/")
		if err != nil {
			log.Println("获取user_id=", userId_, "的用户信息失败：", err)
		}
		response.Comment.User = user
		response.Comment.Id = comment.ID
		response.Comment.Content = comment.Content
		response.Comment.CreateDate = comment.CreateTime.Format("2006-01-02 15:04:05")
		response.StatusMsg = "发布评论成功！"
	} else if request.ActionType == 2 {
		response.StatusMsg = "删除评论成功！"
	}
	response.StatusCode = 0
	c.JSON(200, response)
}

/*
视频评论列表
*/
func CommentList(c *gin.Context) {
	var request CommentListRequest
	var response = &CommentListResponse{}
	if err := c.Bind(&request); err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "request参数绑定失败！"}
		c.JSON(400, response)
		log.Println("request参数绑定失败：", err)
		return
	}

	commentList, err := service.CommentList(request.VideoId)
	if err != nil {
		response.Response = common.Response{StatusCode: 1, StatusMsg: "获取评论列表失败！"}
		c.JSON(400, response)
		log.Println("获取评论列表失败：", err)
		return
	}
	response.StatusCode = 0
	response.StatusMsg = "success"
	response.CommentList = commentList
	c.JSON(200, response)
}
