package repository

import (
	"errors"
	"gorm.io/gorm"
	"sync"
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

func (likeDao *LikeDao) GetLikeByUserIDAndVideoID(UserID int64, VideoId int64) (User, error) {
	f := User{}
	result := db.Where("user_id = ? and video_id = ?", UserID, VideoId).Take(&f)
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
	favorite := Favorite{UserID: UserID, VideoID: VideoId}
	result := db.Select("user_id", "video_id").Create(&favorite) // 通过数据的指针来创建
	return result.Error
}

// DeleteLike 取消点赞
func (likeDao *LikeDao) DeleteLike(UserID int64, VideoId int64) error {
	favorite := Favorite{}
	result := db.Where("user_id = ? and video_id = ?", UserID, VideoId).Delete(&favorite) // 通过数据的指针来创建
	return result.Error
}
