package service

import (
	"DouYIn/repository"
)

var (
	messageDao = repository.NewMessageDaoInstance()
)

// SendMessage
func SendMessage(FromID int64, ToID int64, Content string) error {
	err := messageDao.SendMessage(FromID, ToID, Content)
	return err
}
