package service

import (
	"DouYIn/repository"
	"DouYIn/utils"
	"errors"
	"log"
)

var (
	userDao = repository.NewUserDaoInstance()
)

func GetByID(id int64) (repository.User, error) {
	//获取数据
	//组装数据
	user, err := userDao.GetByID(id)
	if err != nil {
		log.Println("service.GetById error")
		return user, err
	}
	return user, nil
}

func GetByUserName(username string) (repository.User, error) {
	user, err := userDao.GetByUsername(username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func Register(Username string, Password string) (int64, error) {
	encryptedPassword, err := utils.GeneratePassword(Password)
	if err != nil {
		log.Println("注册时，密码加密失败")
		return 0, err
	}
	user, err := userDao.InsertUser(Username, encryptedPassword)
	if err != nil {
		log.Println("service.Register error")
		return user.ID, err
	}
	return user.ID, nil
}

func Login(Username string, Password string) (int64, error) {
	//先检查用户是否存在
	user, err := userDao.GetByUsername(Username)
	if err != nil {
		return 0, errors.New("用户不存在")
	}
	//验证密码
	err = utils.VerifyPassword(Password, user.Password)
	if err != nil { //密码不正确
		return 0, errors.New("密码错误")
	}
	return user.ID, nil
}
