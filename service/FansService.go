package service

import (
	"DouYIn/repository"
	"gorm.io/gorm"
)

var (
	fansDao = repository.NewFansDaoInstance()
)

// FindFolloweeList 查询用户的关注列表
func FindFolloweeList(userId int64) []repository.User {
	return fansDao.SelectFolloweeList(userId)
}

// FindFollowerList 查询用户的粉丝列表
func FindFollowerList(userId int64) []repository.User {
	return fansDao.SelectFollowerList(userId)
}

// FindFriendList 查询用户的好友列表
func FindFriendList(userId int64) []repository.User {
	return fansDao.SelectFriendList(userId)
}

/*
查询是否已经关注
*/

func HasFollowed(bloggerId int64, fansId int64) bool {
	return fansDao.HasFollowed(bloggerId, fansId)
}

//关注与取关操作

func FollowRelationAction(bloggerId int64, fansId int64, actionType int) error {
	bool := HasFollowed(bloggerId, fansId)

	err := repository.Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if actionType == 1 && !bool { //关注
			err := fansDao.InsertFollowRelation(bloggerId, fansId)
			if err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			err = userDao.FollowerInsc(bloggerId)
			if err != nil {
				return err
			}
			err = userDao.FollowInsc(fansId)
			if err != nil {
				return err
			}
		} else if actionType == 2 && bool { //取关
			err := fansDao.DeleteFollowRelation(bloggerId, fansId)
			if err != nil {
				return err
			}
			err = userDao.FollowerDesc(bloggerId)
			if err != nil {
				return err
			}
			err = userDao.FollowDesc(fansId)
			if err != nil {
				return err
			}
		}
		// 返回 nil 提交事务
		return nil
	})
	return err
}
