package service

import (
	"DouYIn/repository"
)

var (
	userVideoDao = repository.NewUserVideoDaoInstance()
	// videoDao = repository.NewVideoDaoInstance()
)

type PublishList struct {
	userId    int64
	VideoList []repository.Video
}

func GetPublishList(userId int64) (PublishList, error) {
	publishList := PublishList{} // 发布列表对象
	uservideos, err := userVideoDao.GetByUserId(userId)
	if err != nil {
		return publishList, err
	}
	publishList.userId = userId
	for _, uservideo := range uservideos {
		videoId := uservideo.VideoID
		println(videoId) // 查询video
	}
	return publishList, nil
}
