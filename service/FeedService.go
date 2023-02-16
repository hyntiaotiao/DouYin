package service

import (
	"DouYIn/common"
	"DouYIn/repository"
	"log"
)

var (
	videoDao = repository.NewVideoDaoInstance()
)

// Feed 返回指定投稿时间之后的amount个视屏
func Feed(amount int, UserId int64, LatestTime ...int64) ([]common.VideoVO, int64, error) {
	var latestTime int64 = 0
	if len(LatestTime) == 1 {
		latestTime = LatestTime[0]
	}
	videos, nextTime, error := videoDao.GetVideos(amount, UserId, latestTime)
	if error != nil {
		log.Println("videoDao.GetVideos 出错")
	}
	if len(videos) == 0 { //说明传入的时间戳刚好是最后一条视屏的投稿时间，导致没有查到数据
		videos, nextTime, error = videoDao.GetVideos(amount, UserId, 0)
		if error != nil {
			log.Println("videoDao.GetVideos 出错")
		}
	}
	if len(videos) < amount { //说明在指定投稿时间之后视屏已经不足amount个了，所以把nextTime设置为0，下次从头开始查询
		nextTime = 0
	}
	// fmt.Println(videos[0])
	return videos, nextTime, error
}
