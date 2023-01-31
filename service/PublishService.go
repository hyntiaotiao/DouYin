package service

import (
	"DouYIn/common"
)

func PublishList(userId int64) ([]common.Video, error) {
	videoList, err := videoDao.GetPublishList(userId)
	if err != nil {
		return videoList, err
	}
	return videoList, nil
}
