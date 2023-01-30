package service

import (
	"DouYIn/repository"
)

var (
	likeDao = repository.NewLikeDaoInstance()
)

// FavouriteAction 点赞与取消赞操作
func FavouriteAction(userId int64, videoId int64, actionType int32) error {
	_, err := likeDao.GetLikeByUserIDAndVideoID(userId, videoId)
	if actionType == 1 && err != nil { //点赞
		err = likeDao.InsertLike(userId, videoId)
		return err
	} else if actionType == 2 { //取消点赞
		err = likeDao.DeleteLike(userId, videoId)
		return err
	}
	return nil
}
