package repository

import (
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

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

func (commentDao *CommentDao) GetCommentList(videoID int64) ([]Comment, error) {
	var comments []Comment
	commentListSQL := " select comment.id,comment.content,comment.create_time,comment.publisher_id from comment" +
		" where comment.video_id = " + fmt.Sprintf("%v", videoID) +
		" order by comment.create_time desc"
	Db.Raw(commentListSQL).Scan(&comments)
	return comments, nil
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
