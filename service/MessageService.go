package service

import "DouYIn/repository"

var (
	messageDao = repository.NewMessageDaoInstance()
)

func MessageChat(user1 int64, user2 int64) ([]repository.Message, error) {
	messageList, err := messageDao.GetChat(user1, user2)
	return messageList, err
}
