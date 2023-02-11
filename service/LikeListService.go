package service

import (
	"DouYIn/common"
	"log"
)

// Feed 返回指定投稿时间之后的amount个视屏
func LikeList(MyID int64, UserId int64) ([]common.Video, error) {
	videos, error := videoDao.GetLikeList(MyID, UserId)
	if error != nil {
		log.Println("videoDao.GetLikeList 出错")
	}
	return videos, error
}
