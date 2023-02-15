package controller

import (
	"DouYin/common"
	"DouYin/service"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type CommentActionRequest struct {
	VideoId     int64  `form:"video_id" json:"videoId" validator:"required,gt=0"`
	Token       string `form:"token" json:"token" validator:"required"`
	ActionType  int32  `form:"action_type" json:"action_type" validator:"required,gte=1,lte=2"` //2删除评论，1发表评论
	CommentID   int64  `form:"comment_id" json:"comment_id,omitempty" validator:"required,gt=0"`
	CommentText string `form:"comment_text" json:"comment_text,omitempty" validator:"omitempty"`
}

type CommentActionResponse struct {
	common.Response
	Comment common.Comment `json:"comment,omitempty"`
}

type CommentListRequest struct {
	VideoId int64  `form:"video_id" json:"videoId"`
	Token   string `form:"token" json:"token"`
}

type CommentListResponse struct {
	common.Response
	CommentList []common.Comment `json:"comment_list" binding:"required"`
}

func CommentAction(c *gin.Context) {
	client := resty.New()
	var request CommentActionRequest
	var response = &CommentActionResponse{}
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}
	temp_curUserID, _ := c.Get("UserID")
	if temp_curUserID != nil {
		curUserID := temp_curUserID.(int64)
		comment, err := service.CommentAction(request.ActionType, request.CommentID, request.VideoId, curUserID, request.CommentText)
		if err != nil {
			log.Println("评论操作失败")
			response.StatusCode = 1
			response.StatusMsg = "评论操作失败"
		} else {
			if request.ActionType == 1 {
				user := common.User{}
				_, err := client.R().
					SetResult(&user).
					SetQueryParams(map[string]string{
						"user_id": strconv.FormatInt(curUserID, 10),
						"token":   request.Token,
					}).
					Get("http://localhost:8080/douyin/user/")
				if err != nil {
					log.Println("获取用户信息失败")
					log.Println(err)
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
		}
		c.JSON(200, response)
	} else {
		c.JSON(400, gin.H{"error": "请先登录！"})
	}

}

func CommentList(c *gin.Context) {
	client := resty.New()
	var request CommentListRequest
	var response = &CommentListResponse{}
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}

	commentList, err1 := service.CommentList(request.VideoId)
	if err1 != nil {
		c.JSON(400, gin.H{"error": err1})
		log.Println("获取评论列表失败")
		return
	}
	commentList_ := make([]common.Comment, len(commentList))
	log.Println("评论列表：", commentList)
	for i := range commentList {
		comment := commentList[i]
		user := common.User{}
		_, err2 := client.R().
			SetResult(&user).
			SetQueryParams(map[string]string{
				"user_id": strconv.FormatInt(comment.PublisherID, 10),
				"token":   request.Token,
			}).
			Get("http://localhost:8080/douyin/user/")
		if err2 != nil {
			log.Println("获取用户信息失败")
			log.Println(err2)
		}
		commentList_[i].User = user
		commentList_[i].CreateDate = comment.CreateTime.Format("2006-01-02 15:04:05")
		commentList_[i].Id = comment.ID
		commentList_[i].Content = comment.Content
	}
	response.StatusCode = 0
	response.StatusMsg = "success"
	response.CommentList = commentList_
	c.JSON(200, response)
}
