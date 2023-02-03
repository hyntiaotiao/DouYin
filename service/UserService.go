package service

import (
	"DouYIn/repository"
	"errors"
	"log"
)

var (
	userDao = repository.NewUserDaoInstance()
)

func Register(Username string, Password string) (int64, error) {
	user, err := userDao.InsertUser(Username, Password)
	if err != nil {
		log.Println("service.Register error")
		return user.ID, err
	}
	return user.ID, nil
}

func Login(Username string, Password string) (int64, error) {
	user, err := userDao.GetByUsername(Username)
	if err != nil {
		return 0, err
	}
	if user.Password != Password { //密码不正确
		return 0, errors.New("密码错误")
	}
	return user.ID, nil
}

func GetUserByID(id int64) (repository.User, error) {
	//获取数据
	//组装数据
	user, err := userDao.GetByID(id)
	if err != nil {
		log.Println("service.GetById error")
		return user, err
	}
	return user, nil
}
