package repository

import (
	"errors"
	"log"
	"sync"
)

var (
	userVideoOnce sync.Once
	userVideoDao  *UserVideoDao
)

// UserVideoDao 即数据访问对象，直接对指定的“某个数据源”的增删改查的封装（这里是对user_video的增删改查）
type UserVideoDao struct {
}

// NewUserVideoDaoInstance 返回一个UserVideoDao实例
func NewUserVideoDaoInstance() *UserVideoDao {
	//不论NewUserVideoDaoInstance()被调用多少次，Do中的内容只会调用一次 (实现了单例生成UserVideoDao)
	userVideoOnce.Do(
		func() {
			//在Go语言中，对结构体进行&取地址操作时，视为对该类型进行一次 new 的实例化操作
			userVideoDao = &UserVideoDao{}
		})
	return userVideoDao
}

func (userVideoDao *UserVideoDao) GetByUserId(userId int64) ([]UserVideo, error) {
	userVideos := []UserVideo{}
	result := db.Where("userId = ?", userId).Find(&userVideos)

	//错误处理
	if result.Error != nil {
		log.Println("UserDao GetByUserId ERROR") //控制台打印日志
		return userVideos, errors.New("发生未知错误")
	}
	return userVideos, nil
}

func (userVideoDao *UserVideoDao) InsertUserVideo(userVideo *UserVideo) error {
	result := db.Create(&userVideo) // 通过数据的指针来创建
	if result.Error != nil {
		log.Println("UserDao InsertUserVideo ERROR") //控制台打印日志
		return errors.New("发生未知错误")
	}
	return nil
}

func (userVideoDao *UserVideoDao) UpdateUserVideo(userVideo *UserVideo) error {
	return db.Save(userVideo).Error
}
