package repository

import (
	"DouYIn/common"
	"fmt"
	"log"
	"sync"

	"gorm.io/gorm"
)

var (
	commentOnce sync.Once
	commentDao  *CommentDao
)

type CommentDao struct{}

var comment_result []struct {
	Id            int64  `json:"id"`
	Publisher_id  int64  `json:"publisher_id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	Content       string `json:"content"`
	CreateDate    string `json:"create_date"`
}

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

func (commentDao *CommentDao) GetCommentList(videoID int64) ([]common.CommentVO, error) {
	commentListSQL := " select comment.id,comment.content,comment.create_time,comment.publisher_id," +
		" user.username as name,user.follow_count,user.follower_count," +
		" IFNULL( (SELECT 1 FROM fans WHERE fans.fans_id = 1 and fans.blogger_id = 1 LIMIT 1) , false ) as is_follow" +
		" from comment join user" +
		" on user.id = comment.publisher_id" +
		" where comment.video_id = " + fmt.Sprintf("%v", videoID) +
		" order by comment.create_time desc"
	Db.Raw(commentListSQL).Scan(&comment_result)
	var commentList = make([]common.CommentVO, len(comment_result))
	for i := 0; i < len(comment_result); i++ {
		commentList[i].Id = comment_result[i].Id
		commentList[i].Content = comment_result[i].Content
		commentList[i].CreateDate = comment_result[i].CreateDate
		commentList[i].User.Id = comment_result[i].Publisher_id
		commentList[i].User.FollowCount = comment_result[i].FollowCount
		commentList[i].User.FollowerCount = comment_result[i].FollowerCount
		commentList[i].User.Name = comment_result[i].Name
		commentList[i].User.IsFollow = comment_result[i].IsFollow
	}
	return commentList, nil
}

func (commentDao *CommentDao) InsertComment(comment *Comment) error {
	// 一个事务
	tx := Db.Begin()

	defer func() {
		if r := recover(); r != nil {
			log.Println("回滚")
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		log.Println("事务开启异常")
	}
	if err := tx.Select("content", "publisher_id", "video_id", "favorite_count").Create(comment).Error; err != nil {
		log.Println("插入评论回滚！")
		tx.Rollback()
	}
	video := &Video{ID: comment.VideoID}
	if err := tx.Model(&video).UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		log.Println("更新视频评论数回滚！")
		tx.Rollback()
	}
	tx.Commit()
	return nil
}

func (commentDao *CommentDao) DeleteComment(comment *Comment) error {
	// 一个事务
	tx := Db.Begin()

	defer func() {
		if r := recover(); r != nil {
			log.Println("回滚")
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		log.Println("事务开启异常")
	}

	if err := Db.Where("id = ?", comment.ID).Delete(comment).Error; err != nil {
		log.Println("删除评论回滚！")
		tx.Rollback()
	}

	video := &Video{ID: comment.VideoID}
	if err := tx.Model(&video).Where("comment_count != 0").UpdateColumn("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		log.Println("更新视频评论数回滚！")
		tx.Rollback()
	}
	tx.Commit()
	return nil
}
