package service

import "DouYIn/repository"

var (
	fansDao = repository.NewFansDaoInstance()
)

func HasFollowed(bloggerId int64, fansId int64) bool {
	return fansDao.HasFollowed(bloggerId, fansId)
}

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
