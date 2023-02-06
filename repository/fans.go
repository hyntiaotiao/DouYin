package repository

import (
	"sync"
)

var (
	fansOnce sync.Once
	fansDao  *FansDao
)

type FansDao struct{}

// NewFansDaoInstance 返回一个FansDao实例
func NewFansDaoInstance() *FansDao {
	fansOnce.Do(
		func() {
			fansDao = &FansDao{}
		})
	return fansDao
}

// HasFollowed 查询fansId对应的用户是否关注了bloggerId对应的用户
func (fansDao *FansDao) HasFollowed(bloggerId int64, fansId int64) bool {
	var count int64
	db.Model(&Fans{}).Where("blogger_id = ? and fans_id = ?", bloggerId, fansId).Count(&count)
	if count != 0 {
		return true
	}
	return false
}

// SelectFolloweeList 返回关注列表
func (fansDao *FansDao) SelectFolloweeList(userId int64) []User {
	var followeeList []User
	db.Debug().Table("user").Where("id in (?)", db.Table("fans").Select("blogger_id").Where("fans_id = ?", userId)).Find(&followeeList)
	return followeeList
}

// SelectFollowerList 返回粉丝列表
func (fansDao *FansDao) SelectFollowerList(userId int64) []User {
	var followerList []User
	db.Debug().Table("user").Where("id in (?)", db.Table("fans").Select("fans_id").Where("blogger_id = ?", userId)).Find(&followerList)
	return followerList
}

// SelectFriendList 返回好友列表
func (fansDao *FansDao) SelectFriendList(userId int64) []User {
	var friendList []User
	db.Debug().Table("user u").
		Select("u.*").
		Joins("join fans f1 on u.id = f1.blogger_id").
		Joins("Join fans f2 on u.id = f2.fans_id").
		Where("f1.fans_id = ? and f2.blogger_id = ?", userId, userId).Find(&friendList)
	return friendList
}
