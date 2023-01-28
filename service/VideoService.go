package service

import (
	"DouYIn/repository"
	"log"
)

var (
	VideoDao = repository.NewvideoDaoInstance()
)

func AddVideo(authodId int, playUrl string, coverUrl string, title string) error {

	//1. 获取视频地址和图片地址
	playUrl = "D:\\newproject\\go\\douyin\\video" + playUrl
	coverUrl = "D:\\newproject\\go\\douyin\\video" + coverUrl
	err := VideoDao.Addvideo(authodId, playUrl, coverUrl, title)
	if err != nil {
		log.Println("发布视频失败，请稍后再试")
	}
	return nil
}
