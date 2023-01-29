package repository

import (
	"errors"
	"log"
	"sync"
)

// Like 表的结构

type LikeDao struct {
	Id      int64 `gorm:"column:id;primary_key" json:"id"`                //自增主键
	UserId  int64 `gorm:"column:user_id;NOT NULL" json:"user_id"`         //点赞用户id
	VideoId int64 `gorm:"column:video_id;NOT NULL" json:"video_id"`       //视频id
	Cancel  int   `gorm:"column:cancel;default:0;NOT NULL" json:"cancel"` //是否点赞，0为点赞，1为取消赞
}

// TableName 修改表名映射
func (LikeDao) TableName() string {
	return "likes"
}

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

// UpdateLike 根据userId，videoId,actionType点赞或者取消赞
func (likeDao *LikeDao) UpdateLike(userId int64, videoId int64, actionType int32) error {
	//更新当前用户观看视频的点赞状态“cancel”，返回错误结果
	err := db.Model(LikeDao{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("cancel", actionType).Error
	//如果出现错误，返回更新数据库失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("update data fail")
	}
	//更新操作成功
	return nil
}

// InsertLike 插入点赞数据
func (likeDao *LikeDao) InsertLike(likeData LikeDao) error {
	//创建点赞数据，默认为点赞，cancel为0，返回错误结果
	err := db.Model(LikeDao{}).Create(&likeData).Error
	//如果有错误结果，返回插入失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert data fail")
	}
	return nil
}

// GetLikeInfo 根据userId,videoId查询点赞信息
func (likeDao *LikeDao) GetLikeInfo(userId int64, videoId int64) (LikeDao, error) {
	//创建一条空like结构体，用来存储查询到的信息
	var likeInfo LikeDao
	//根据userid,videoId查询是否有该条信息，如果有，存储在likeInfo,返回查询结果
	err := db.Model(LikeDao{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		First(&likeInfo).Error
	if err != nil {
		//查询数据为0，打印"can't find data"，返回空结构体，这时候就应该要考虑是否插入这条数据了
		if "record not found" == err.Error() {
			log.Println("can't find data")
			return LikeDao{}, nil
		} else {
			//如果查询数据库失败，返回获取likeInfo信息失败
			log.Println(err.Error())
			return likeInfo, errors.New("get likeInfo failed")
		}
	}
	return likeInfo, nil
}

// GetLikeUserIdList 根据videoId获取点赞userId
func (likeDao *LikeDao) GetLikeUserIdList(videoId int64) ([]int64, error) {
	var likeUserIdList []int64 //存所有该视频点赞用户id；
	//查询likes表对应视频id点赞用户，返回查询结果
	err := db.Model(LikeDao{}).Where(map[string]interface{}{"video_id": videoId, "cancel": 0}).
		Pluck("user_id", &likeUserIdList).Error
	//查询过程出现错误，返回默认值0，并输出错误信息
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("get likeUserIdList failed")
	} else {
		//没查询到或者查询到结果，返回数量以及无报错
		return likeUserIdList, nil
	}
}
