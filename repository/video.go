package repository

import (
	"DouYIn/common"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
)

var (
	videoOnce sync.Once
	videoDao  *VideoDao
)

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

// VideoDao 即数据访问对象，直接对指定的“某个数据源”的增删改查的封装（这里是对video的增删改查）
type VideoDao struct {
}

// NewVideoDaoInstance 返回一个UserVideoDao实例
func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

func (videoDao *VideoDao) InsertVideo(video *Video) error {
	result := db.Create(&video) // 通过数据的指针来创建
	if result.Error != nil {
		log.Println("VideoDao InsertVideo ERROR") //控制台打印日志
		return errors.New("发生未知错误")
	}
	return nil
}

func (videoDao *VideoDao) UpdateVideo(video *Video) error {
	return db.Save(video).Error
}

func (videoDao *VideoDao) GetPublishList(UserID int64) ([]common.Video, error) {
	VideoListSQL := " select video.id,video.play_url,video.cover_url,video.title,video.comment_count,video.favourite_count," +
		" video.author_id,user.username as name,user.follow_count,user.follower_count," +
		" IFNULL( (SELECT 1 FROM	favorite WHERE favorite.user_id = " + fmt.Sprintf("%v", UserID) + " and favorite.video_id = video.id LIMIT 1) , false ) as is_favorite," +
		" IFNULL( (SELECT 1 FROM	fans WHERE fans.fans_id = " + fmt.Sprintf("%v", UserID) + " and fans.blogger_id = 1 LIMIT 1) , false ) as is_follow" +
		" from user join video" +
		" on video.author_id = user.id" +
		" order by video.create_time"
	db.Raw(VideoListSQL).Scan(&result)
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
	return VideoList, nil
}

func (videoDao VideoDao) GetVideos(amount int, UserID any, LatestTime int64) ([]common.Video, int64, error) {
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
