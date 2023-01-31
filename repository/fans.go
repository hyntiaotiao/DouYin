package repository

import (
	"sync"
)

var (
	fansOnce sync.Once
	fansDao  *FansDao
)

type FansDao struct{}

// NewFansDaoInstance 返回一个FansDao实例
func NewFansDaoInstance() *FansDao {
	fansOnce.Do(
		func() {
			fansDao = &FansDao{}
		})
	return fansDao
}

// HasFollowed 查询fansId对应的用户是否关注了bloggerId对应的用户
func (fansDao *FansDao) HasFollowed(bloggerId int64, fansId int64) bool {
	var count int64
	db.Model(&Fans{}).Where("blogger_id = ? and fans_id = ?", bloggerId, fansId).Count(&count)
	if count != 0 {
		return true
	}
	return false
}

// insertFollowRelation 插入一条关注数据

func (fansDao *FansDao) InsertFollowRelation(bloggerId int64, fansId int64) error {
	fans := Fans{BloggerID: bloggerId, FansID: fansId}
	result := db.Select("blogger_id", "fans_id").Create(&fans) // 通过数据的指针来创建
	return result.Error
}

// deleteFollowRelation 删除一条关注数据

func (fansDao *FansDao) DeleteFollowRelation(bloggerId int64, fansId int64) error {
	fans := Fans{}
	result := db.Where("blogger_id = ? and fans_id = ?", bloggerId, fansId).Delete(&fans) // 通过数据的指针来删除
	return result.Error
}
