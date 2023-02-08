package service

import (
	"DouYIn/common"
	"DouYIn/repository"
)

var (
	messageDao = repository.NewMessageDaoInstance()
)

func SendMessage(FromID int64, ToID int64, Content string) error {
	err := messageDao.SendMessage(FromID, ToID, Content)
	return err
}

func MessageChat(user1 int64, user2 int64) ([]common.MessageVO, error) {
	messageList, err := messageDao.GetChat(user1, user2)
	return messageList, err
}
