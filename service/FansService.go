package service

import "DouYIn/repository"

var (
	fansDao = repository.NewFansDaoInstance()
)

/*
查询是否有关注关系
*/

func HasFollowed(bloggerId int64, fansId int64) bool {
	return fansDao.HasFollowed(bloggerId, fansId)
}

//关注与取关操作

func FollowRelationAction(bloggerId int64, fansId int64, actionType int32) error {
	bool := HasFollowed(bloggerId, fansId)
	if actionType == 1 && !bool { //关注
		err := fansDao.InsertFollowRelation(bloggerId, fansId)
		return err
	} else if actionType == 2 && bool { //取关
		err := fansDao.DeleteFollowRelation(bloggerId, fansId)
		return err
	}
	return nil
}
