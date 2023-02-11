package repository

import (
	"errors"
	"log"
	"sync"

	"gorm.io/gorm"
)

type LikeDao struct{}

var (
	likeOnce sync.Once
	likeDao  *LikeDao
)

func NewLikeDaoInstance() *LikeDao {
	//不论NewLikeDaoInstance()被调用多少次，Do中的内容只会调用一次
	likeOnce.Do(
		func() {
			//在Go语言中，对结构体进行&取地址操作时，视为对该类型进行一次 new 的实例化操作
			likeDao = &LikeDao{}
		})
	return likeDao
}

func (likeDao *LikeDao) GetLikeByUserIDAndVideoID(UserID int64, VideoId int64) (Favorite, error) {
	f := Favorite{}
	result := Db.Where("user_id = ? and video_id = ?", UserID, VideoId).Take(&f)
	//错误处理
	if result.Error != nil {
		//当 First、Last、Take 方法找不到记录时，GORM 会返回 ErrRecordNotFound 错误
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return f, errors.New("找不到指定的记录")
		}
		return f, errors.New("发生未知错误")
	}
	return f, nil
}

// InsertLike 插入点赞数据
func (likeDao *LikeDao) InsertLike(UserID int64, VideoId int64) error {
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

	favorite := &Favorite{UserID: UserID, VideoID: VideoId}
	if err := tx.Select("user_id", "video_id").Create(&favorite).Error; err != nil {
		log.Println("添加点赞回滚！错误信息：", err)
		tx.Rollback()
	}
	video := &Video{ID: VideoId}
	if err := tx.Model(&video).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error; err != nil {
		log.Println("更新视频点赞数回滚！错误信息：", err)
		tx.Rollback()
	}
	tx.Commit()
	return nil
}

// DeleteLike 取消点赞
func (likeDao *LikeDao) DeleteLike(UserID int64, VideoId int64) error {
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

	favorite := &Favorite{VideoID: VideoId, UserID: UserID}
	if err := tx.Where("user_id = ? and video_id = ?", UserID, VideoId).Delete(&favorite).Error; err != nil {
		log.Println("取消点赞回滚！错误信息：", err)
		tx.Rollback()
	}

	video := &Video{ID: VideoId}
	if err := tx.Model(&video).Where("favorite_count != 0").UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error; err != nil {
		log.Println("更新视频点赞数回滚！错误信息：", err)
		tx.Rollback()
	}
	tx.Commit()
	return nil
}
