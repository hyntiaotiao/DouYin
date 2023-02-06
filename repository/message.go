package repository

import (
	"fmt"
	"sync"
)

var (
	messageOnce sync.Once
	messageDao  *MessageDao
)

type MessageDao struct{}

func NewMessageDaoInstance() *MessageDao {
	messageOnce.Do(
		func() {
			messageDao = &MessageDao{}
		})
	return messageDao
}

func (messageDao *MessageDao) GetChat(userID1 int64, userID2 int64) ([]Message, error) {
	var messageList []Message
	messageListSQL := " select id,send_user_id,receive_user_id,content,create_time from message" +
		" where (send_user_id = " + fmt.Sprintf("%v", userID1) + "and receive_user_id = " + fmt.Sprintf("%v", userID2) + ")" +
		" or (send_user_id = " + fmt.Sprintf("%v", userID2) + "and receive_user_id = " + fmt.Sprintf("%v", userID1) + ")" +
		" order by create_time desc"
	db.Raw(messageListSQL).Scan(&messageList)
	return messageList, nil
}
