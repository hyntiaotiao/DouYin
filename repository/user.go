package repository

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	userOnce sync.Once
	userDao  *UserDao
)

// UserDao 即数据访问对象，直接对指定的“某个数据源”的增删改查的封装（这里是对User的增删改查）
type UserDao struct{}

// NewUserDaoInstance 返回一个UserDao实例
func NewUserDaoInstance() *UserDao {
	//不论NewUserDaoInstance()被调用多少次，Do中的内容只会调用一次 (实现了单例生成UserDao)
	userOnce.Do(
		func() {
			//在Go语言中，对结构体进行&取地址操作时，视为对该类型进行一次 new 的实例化操作
			userDao = &UserDao{}
		})
	return userDao
}

// GetById 根据用户id查询user对象
// 如果不存在对应的用户，则方法返回的error非空
func (userDao *UserDao) GetById(id int64) (User, error) {
	u := User{}
	result := db.Where("id = ?", id).Take(&u)

	//错误处理
	if result.Error != nil {
		log.Println("UserDao GetByID ERROR") //控制台打印日志
		//当 First、Last、Take 方法找不到记录时，GORM 会返回 ErrRecordNotFound 错误
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return u, errors.New("找不到指定的记录")
		}
		return u, errors.New("发生未知错误")
	}
	return u, nil
}

func (userDao *UserDao) GetByUsername(username string) (User, error) {
	u := User{}
	result := db.Where("username = ?", username).Take(&u)

	//错误处理
	if result.Error != nil {
		log.Println("UserDao. GetByUsername ERROR") //控制台打印日志
		//当 First、Last、Take 方法找不到记录时，GORM 会返回 ErrRecordNotFound 错误
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return u, errors.New("找不到指定的记录")
		}
		return u, result.Error
	}
	return u, nil
}

func (userDao *UserDao) InsertUser(username string, password string) (User, error) {
	user := User{Username: username, Password: password}
	result := db.Select("username", "password").Create(&user) // 通过数据的指针来创建
	return user, result.Error
}
