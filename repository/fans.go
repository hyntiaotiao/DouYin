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
	Db.Model(&Fans{}).Where("blogger_id = ? and fans_id = ?", bloggerId, fansId).Count(&count)
	if count != 0 {
		return true
	}
	return false
}

// SelectFolloweeList 返回关注列表
func (fansDao *FansDao) SelectFolloweeList(userId int64) []User {
	var followeeList []User
	Db.Debug().Table("user").Where("id in (?)", Db.Table("fans").Select("blogger_id").Where("fans_id = ?", userId)).Find(&followeeList)
	return followeeList
}

// SelectFollowerList 返回粉丝列表
func (fansDao *FansDao) SelectFollowerList(userId int64) []User {
	var followerList []User
	Db.Debug().Table("user").Where("id in (?)", Db.Table("fans").Select("fans_id").Where("blogger_id = ?", userId)).Find(&followerList)
	return followerList
}

// SelectFriendList 返回好友列表
func (fansDao *FansDao) SelectFriendList(userId int64) []User {
	var friendList []User
	Db.Debug().Table("user u").
		Select("u.*").
		Joins("join fans f1 on u.id = f1.blogger_id").
		Joins("Join fans f2 on u.id = f2.fans_id").
		Where("f1.fans_id = ? and f2.blogger_id = ?", userId, userId).Find(&friendList)
	return friendList
}

// insertFollowRelation 插入一条关注数据

func (fansDao *FansDao) InsertFollowRelation(bloggerId int64, fansId int64) error {
	fans := Fans{BloggerID: bloggerId, FansID: fansId}
	result := Db.Select("blogger_id", "fans_id").Create(&fans) // 通过数据的指针来创建
	return result.Error
}

// deleteFollowRelation 删除一条关注数据

func (fansDao *FansDao) DeleteFollowRelation(bloggerId int64, fansId int64) error {
	fans := Fans{}
	result := Db.Where("blogger_id = ? and fans_id = ?", bloggerId, fansId).Delete(&fans) // 通过数据的指针来删除
	return result.Error
}
