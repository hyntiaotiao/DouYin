package service

import (
	"DouYIn/common"
	"DouYIn/repository"
	"log"
)

var (
	favoriteDao = repository.NewFavoriteDaoInstance()
)

// Feed 返回指定投稿时间之后的amount个视屏
func FavoriteList(MyID int64, UserId int64) ([]common.VideoVO, error) {
	videos, error := videoDao.GetLikeList(MyID, UserId)
	if error != nil {
		log.Println("videoDao.GetLikeList 出错")
	}
	return videos, error
}

// FavoriteAction 点赞与取消赞操作
func FavoriteAction(userId int64, videoId int64, actionType int32) error {
	// _, err := likeDao.GetLikeByUserIDAndVideoID(userId, videoId)
	if actionType == 1 { //点赞
		err := favoriteDao.InsertFavorite(userId, videoId)
		return err
	} else if actionType == 2 { //取消点赞
		err := favoriteDao.DeleteFavorite(userId, videoId)
		return err
	}
	return nil
}
