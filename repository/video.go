package repository

import (
	"errors"
	"log"
	"sync"
)

var (
	videoOnce sync.Once
	videoDao  *VideoDao
)

// VideoDao 即数据访问对象，直接对指定的“某个数据源”的增删改查的封装（这里是对video的增删改查）
type VideoDao struct {
}

// NewVideoDaoInstance 返回一个UserVideoDao实例
func NewVideoDaoInstance() *VideoDao {
	//不论NewVideoDaoInstance()被调用多少次，Do中的内容只会调用一次 (实现了单例生成VideoDao)
	videoOnce.Do(
		func() {
			//在Go语言中，对结构体进行&取地址操作时，视为对该类型进行一次 new 的实例化操作
			videoDao = &VideoDao{}
		})
	return videoDao
}

func (videoDao *VideoDao) GetAllByAuthorID(authorID int64) ([]Video, error) {
	videos := []Video{}
	result := db.Where("author_id = ? and ", authorID).Find(&videos)

	//错误处理
	if result.Error != nil {
		log.Println("VideoDao GetAllByAuthorID ERROR") //控制台打印日志
		return videos, errors.New("发生未知错误")
	}
	return videos, nil
}

func (videoDao *VideoDao) InsertVideo(video *Video) error {
	result := db.Create(&video) // 通过数据的指针来创建
	if result.Error != nil {
		log.Println("VideoDao InsertVideo ERROR") //控制台打印日志
		return errors.New("发生未知错误")
	}
	return nil
}

func (videoDao *VideoDao) UpdateVideo(video *Video) error {
	return db.Save(video).Error
}
