package service

import (
	"DouYIn/repository"
)

var (
	videoDao = repository.NewVideoDaoInstance()
)

func GetPublishList(userId int64) ([]repository.Video, error) {
	videos, err := videoDao.GetAllByAuthorID(userId)
	if err != nil {
		return videos, err
	}
	return videos, nil
}
