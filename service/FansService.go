package service

import "DouYIn/repository"

var (
	fansDao = repository.NewFansDaoInstance()
)

func HasFollowed(bloggerId int64, fansId int64) bool {
	return fansDao.HasFollowed(bloggerId, fansId)
}