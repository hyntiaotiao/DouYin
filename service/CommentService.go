package service

import (
	"DouYIn/common"
	"DouYIn/repository"
)

var (
	commentDao = repository.NewCommentDaoInstance()
)

// CommentList  点赞与取消赞操作
func CommentList(videoId int64) ([]common.CommentVO, error) {
	commentList, error := commentDao.GetCommentList(videoId)
	return commentList, error
}

func CommentAction(action_type int32, comment_id int64, video_id int64, publisher_id int64, comment_context string) (*repository.Comment, error) {
	// 插入评论
	comment := &repository.Comment{ID: comment_id, VideoID: video_id, PublisherID: publisher_id, Content: comment_context, AuthorID: 1, FavoriteCount: 0}
	var err error
	if action_type == 1 {
		err = commentDao.InsertComment(comment)
	} else if action_type == 2 {
		err = commentDao.DeleteComment(comment)
	}
	return comment, err
}
