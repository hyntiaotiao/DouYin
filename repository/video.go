package repository

import (
	"DouYIn/common"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

var (
	videoOnce sync.Once
	videoDao  *VideoDao
)

// VideoDao 即数据访问对象，直接对指定的“某个数据源”的增删改查的封装（这里是对User的增删改查）
type VideoDao struct{}

// NewVideoDaoInstance 返回一个UserDao实例
func NewVideoDaoInstance() *VideoDao {
	//不论NewUserDaoInstance()被调用多少次，Do中的内容只会调用一次 (实现了单例生成UserDao)
	videoOnce.Do(
		func() {
			//在Go语言中，对结构体进行&取地址操作时，视为对该类型进行一次 new 的实例化操作
			videoDao = &VideoDao{}
		})
	return videoDao
}

func (videoDao *VideoDao) Addvideo(authorId int64, playUrl string, coverUrl string, title string) error {
	newVideo := &Video{
		AuthorID: authorId,
		CoverUrl: "http://" + coverUrl,
		PlayUrl:  "http://" + playUrl,
		Title:    title,
	}
	// 重复视频校验（没做）
	err := db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Create(newVideo).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return err
}

func (videoDao VideoDao) GetVideos(amount int, UserID any, LatestTime int64) ([]common.Video, int64, error) {
	var result []struct {
		Id            int64  `json:"id"`
		AuthorId      int64  `json:"author_id"`
		Name          string `json:"name"`
		FollowCount   int64  `json:"follow_count"`
		FollowerCount int64  `json:"follower_count"`
		IsFollow      bool   `json:"is_follow"`
		PlayUrl       string `json:"play_url"`
		CoverUrl      string `json:"cover_url"`
		FavoriteCount int64  `json:"favorite_count"`
		CommentCount  int64  `json:"comment_count"`
		IsFavorite    bool   `json:"is_favorite"`
		Title         string `json:"title"`
	}
	var VideoListSQL string
	var NextTimeSQL string
	if UserID == -1 {
		VideoListSQL = "select video.id , video.play_url,video.cover_url,video.favourite_count,video.comment_count,video.title, video.author_id,user.username as name,user.follow_count," +
			" user.follower_count" +
			" from video inner join user" +
			" on author_id = user.id" +
			" where UNIX_TIMESTAMP(video.create_time) > " + strconv.FormatInt(LatestTime, 10) +
			" order by video.create_time limit " + strconv.Itoa(amount)
	} else {
		VideoListSQL = " select video.id,video.play_url,video.cover_url,video.title,video.comment_count,video.favourite_count," +
			" video.author_id,user.username as name,user.follow_count,user.follower_count," +
			" IFNULL( (SELECT 1 FROM	favorite WHERE favorite.user_id = " + fmt.Sprintf("%v", UserID) + " and favorite.video_id = video.id LIMIT 1) , false ) as is_favorite," +
			" IFNULL( (SELECT 1 FROM	fans WHERE fans.fans_id = " + fmt.Sprintf("%v", UserID) + " and fans.blogger_id = 1 LIMIT 1) , false ) as is_follow" +
			" from user join video" +
			" on video.author_id = user.id" +
			" where UNIX_TIMESTAMP(video.create_time) > " + strconv.FormatInt(LatestTime, 10) +
			" order by video.create_time  LIMIT " + strconv.Itoa(amount)
	}
	NextTimeSQL = " select UNIX_TIMESTAMP(video.create_time) as time" +
		" from video inner join user" +
		" on author_id = user.id" +
		" where UNIX_TIMESTAMP(video.create_time)>" + strconv.FormatInt(LatestTime, 10) +
		" order by time limit 1," + strconv.Itoa(amount-1)
	db.Raw(VideoListSQL).Scan(&result)
	var NextTime int64
	db.Raw(NextTimeSQL).Scan(&NextTime)
	var VideoList = make([]common.Video, len(result))
	for i := 0; i < len(result); i++ {
		VideoList[i].Author.Id = result[i].AuthorId
		VideoList[i].Author.FollowCount = result[i].FollowCount
		VideoList[i].Author.FollowerCount = result[i].FollowerCount
		VideoList[i].Author.Name = result[i].Name
		VideoList[i].Author.IsFollow = result[i].IsFollow
		VideoList[i].Id = result[i].Id
		VideoList[i].PlayUrl = result[i].PlayUrl
		VideoList[i].CoverUrl = result[i].CoverUrl
		VideoList[i].IsFavorite = result[i].IsFavorite
		VideoList[i].Title = result[i].Title
		VideoList[i].CommentCount = result[i].CommentCount
		VideoList[i].FavoriteCount = result[i].FavoriteCount

	}
	return VideoList, NextTime, nil
}

//gorm:"foreignkey:Id;references:UserID;
