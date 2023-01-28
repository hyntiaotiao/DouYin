package repository

import (
	"sync"
)

var (
	videoOnce sync.Once
	videoDao  *VideoDao
)

type VideoDao struct {
}

func NewvideoDaoInstance() *VideoDao {
	userOnce.Do(func() {
		videoDao = &VideoDao{}
	})
	return videoDao
}

// 新增视频
func (videoDao *VideoDao) Addvideo(authorId int, playUrl string, coverUrl string, title string) error {
	newVideo := &Video{
		Author_id: authorId,
		CoverUrl:  coverUrl,
		PlayUrl:   playUrl,
		Title:     title,
	}
	videoResult := db.Create(newVideo)
	return videoResult.Error
}
