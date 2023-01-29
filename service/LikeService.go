package service

import (
	"DouYIn/repository"
	"errors"
)

var (
	likeDao = repository.NewLikeDaoInstance()
)

//点赞与取消赞操作

func FavouriteAction(userId int64, videoId int64, actionType int32) (repository.LikeDao, error) {
	likeInfo, err := likeDao.GetLikeInfo(userId, videoId)
	//查询信息出错
	if err != nil {
		return likeInfo, errors.New("点赞信息查询失败")
	}
	//如果未查询到信息，返回的是空结构体，此时插入数据
	if likeInfo == (repository.LikeDao{}) {
		err := likeDao.InsertLike(likeInfo)
		if err != nil {
			return likeInfo, errors.New("插入点赞信息失败")
		}
	}
	//如果查询到了点赞信息，说明需要更新点赞信息
	err = likeDao.UpdateLike(userId, videoId, actionType)
	if err != nil {
		return likeInfo, errors.New("更新点赞信息失败")
	}
	return likeInfo, nil
}
