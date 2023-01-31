package repository

import (
	"errors"
	"gorm.io/gorm"
	"sync"
)

type LikeDao struct{}

var (
	likeOnce sync.Once
	likeDao  *LikeDao
)

func NewLikeDaoInstance() *LikeDao {
	//不论NewLikeDaoInstance()被调用多少次，Do中的内容只会调用一次
	likeOnce.Do(
		func() {
			//在Go语言中，对结构体进行&取地址操作时，视为对该类型进行一次 new 的实例化操作
			likeDao = &LikeDao{}
		})
	return likeDao
}

func (likeDao *LikeDao) GetLikeByUserIDAndVideoID(UserID int64, VideoId int64) (User, error) {
	f := User{}
	result := db.Where("user_id = ? and video_id = ?", UserID, VideoId).Take(&f)
	//错误处理
	if result.Error != nil {
		//当 First、Last、Take 方法找不到记录时，GORM 会返回 ErrRecordNotFound 错误
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return f, errors.New("找不到指定的记录")
		}
		return f, errors.New("发生未知错误")
	}
	return f, nil
}

// InsertLike 插入点赞数据
func (likeDao *LikeDao) InsertLike(UserID int64, VideoId int64) error {
	favorite := Favorite{UserID: UserID, VideoID: VideoId}
	result := db.Select("user_id", "video_id").Create(&favorite) // 通过数据的指针来创建
	return result.Error
}

// DeleteLike 取消点赞
func (likeDao *LikeDao) DeleteLike(UserID int64, VideoId int64) error {
	favorite := Favorite{}
	result := db.Where("user_id = ? and video_id = ?", UserID, VideoId).Delete(&favorite) // 通过数据的指针来创建
	return result.Error
}

//// GetLikeInfo 根据userId,videoId查询点赞信息
//func (likeDao *LikeDao) GetLikeInfo(userId int64, videoId int64) (LikeDao, error) {
//	//创建一条空like结构体，用来存储查询到的信息
//	var likeInfo LikeDao
//	//根据userid,videoId查询是否有该条信息，如果有，存储在likeInfo,返回查询结果
//	err := db.Model(LikeDao{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
//		First(&likeInfo).Error
//	if err != nil {
//		//查询数据为0，打印"can't find data"，返回空结构体，这时候就应该要考虑是否插入这条数据了
//		if "record not found" == err.Error() {
//			log.Println("can't find data")
//			return LikeDao{}, nil
//		} else {
//			//如果查询数据库失败，返回获取likeInfo信息失败
//			log.Println(err.Error())
//			return likeInfo, errors.New("get likeInfo failed")
//		}
//	}
//	return likeInfo, nil
//}
//
//// GetLikeUserIdList 根据videoId获取点赞userId
//func (likeDao *LikeDao) GetLikeUserIdList(videoId int64) ([]int64, error) {
//	var likeUserIdList []int64 //存所有该视频点赞用户id；
//	//查询likes表对应视频id点赞用户，返回查询结果
//	err := db.Model(LikeDao{}).Where(map[string]interface{}{"video_id": videoId, "cancel": 0}).
//		Pluck("user_id", &likeUserIdList).Error
//	//查询过程出现错误，返回默认值0，并输出错误信息
//	if err != nil {
//		log.Println(err.Error())
//		return nil, errors.New("get likeUserIdList failed")
//	} else {
//		//没查询到或者查询到结果，返回数量以及无报错
//		return likeUserIdList, nil
//	}
//}
