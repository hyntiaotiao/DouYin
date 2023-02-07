package repository

import (
	"sync"
)

var (
	messageOnce sync.Once
	messageDao  *MessageDao
)

// MessageDao 即数据访问对象，直接对指定的“某个数据源”的增删改查的封装（这里是对Message的增删改查）
type MessageDao struct{}

// NewMessageDaoInstance 返回一个MessageDao实例
func NewMessageDaoInstance() *MessageDao {
	//不论NewMessageDaoInstance()被调用多少次，Do中的内容只会调用一次 (实现了单例生成MessageDao)
	messageOnce.Do(
		func() {
			//在Go语言中，对结构体进行&取地址操作时，视为对该类型进行一次 new 的实例化操作
			messageDao = &MessageDao{}
		})
	return messageDao
}

func (messageDao *MessageDao) SendMessage(FromID int64, ToID int64, Content string) error {
	var message = Message{
		Content:       Content,
		SendUserID:    FromID,
		ReceiveUserID: ToID,
	}
	result := Db.Select("content", "send_user_id", "receive_user_id").Create(&message)
	return result.Error
}
