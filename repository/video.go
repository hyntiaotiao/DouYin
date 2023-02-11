package repository

import (
	"DouYIn/common"
	"fmt"
	"strconv"
	"sync"

	"gorm.io/gorm"
)

var (
	videoOnce sync.Once
	videoDao  *VideoDao
)

type result = struct {
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

var video_result []struct {
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

func (videoDao *VideoDao) GetPublishList(UserID int64) ([]common.Video, error) {
	VideoListSQL := " select video.id,video.play_url,video.cover_url,video.title,video.comment_count,video.favorite_count," +
		" video.author_id,user.username as name,user.follow_count,user.follower_count," +
		" IFNULL( (SELECT 1 FROM	favorite WHERE favorite.user_id = " + fmt.Sprintf("%v", UserID) + " and favorite.video_id = video.id LIMIT 1) , false ) as is_favorite," +
		" IFNULL( (SELECT 1 FROM	fans WHERE fans.fans_id = " + fmt.Sprintf("%v", UserID) + " and fans.blogger_id = 1 LIMIT 1) , false ) as is_follow" +
		" from user join video" +
		" on video.author_id = user.id" +
		" order by video.create_time"
	Db.Raw(VideoListSQL).Scan(&video_result)
	var VideoList = make([]common.Video, len(video_result))
	for i := 0; i < len(video_result); i++ {
		VideoList[i].Author.Id = video_result[i].AuthorId
		VideoList[i].Author.FollowCount = video_result[i].FollowCount
		VideoList[i].Author.FollowerCount = video_result[i].FollowerCount
		VideoList[i].Author.Name = video_result[i].Name
		VideoList[i].Author.IsFollow = video_result[i].IsFollow
		VideoList[i].Id = video_result[i].Id
		VideoList[i].PlayUrl = video_result[i].PlayUrl
		VideoList[i].CoverUrl = video_result[i].CoverUrl
		VideoList[i].IsFavorite = video_result[i].IsFavorite
		VideoList[i].Title = video_result[i].Title
		VideoList[i].CommentCount = video_result[i].CommentCount
		VideoList[i].FavoriteCount = video_result[i].FavoriteCount
	}
	return VideoList, nil
}

func (videoDao VideoDao) GetVideos(amount int, UserID any, LatestTime int64) ([]common.Video, int64, error) {
	var VideoListSQL string
	var NextTimeSQL string
	if UserID == -1 {
		VideoListSQL = "select video.id , video.play_url,video.cover_url,video.favorite_count,video.comment_count,video.title, video.author_id,user.username as name,user.follow_count," +
			" user.follower_count" +
			" from video inner join user" +
			" on author_id = user.id" +
			" where UNIX_TIMESTAMP(video.create_time) > " + strconv.FormatInt(LatestTime, 10) +
			" order by video.create_time limit " + strconv.Itoa(amount)
	} else {
		VideoListSQL = " select video.id,video.play_url,video.cover_url,video.title,video.comment_count,video.favorite_count," +
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
	Db.Raw(VideoListSQL).Scan(&video_result)
	var NextTime int64
	Db.Raw(NextTimeSQL).Scan(&NextTime)
	var VideoList = make([]common.Video, len(video_result))
	for i := 0; i < len(video_result); i++ {
		VideoList[i].Author.Id = video_result[i].AuthorId
		VideoList[i].Author.FollowCount = video_result[i].FollowCount
		VideoList[i].Author.FollowerCount = video_result[i].FollowerCount
		VideoList[i].Author.Name = video_result[i].Name
		VideoList[i].Author.IsFollow = video_result[i].IsFollow
		VideoList[i].Id = video_result[i].Id
		VideoList[i].PlayUrl = video_result[i].PlayUrl
		VideoList[i].CoverUrl = video_result[i].CoverUrl
		VideoList[i].IsFavorite = video_result[i].IsFavorite
		VideoList[i].Title = video_result[i].Title
		VideoList[i].CommentCount = video_result[i].CommentCount
		VideoList[i].FavoriteCount = video_result[i].FavoriteCount
	}
	return VideoList, NextTime, nil
}

func (videoDao VideoDao) GetLikeList(MyID int64, UserID int64) ([]common.Video, error) {

	var res []result
	VideoListSQL := "select video.id,video.play_url,video.cover_url,video.title,video.comment_count,video.favorite_count ," +
		" video.author_id,user.username as name,user.follow_count,user.follower_count,1 as is_favorite," +
		" IFNULL( (SELECT 1 FROM fans WHERE fans.fans_id =  " + fmt.Sprintf("%v", MyID) +
		" and fans.blogger_id =  " + fmt.Sprintf("%v", UserID) +
		" LIMIT 1) , false ) as is_follow" +
		" from user join video" +
		" on video.author_id = user.id" +
		" where IFNULL( (SELECT 1 FROM favorite WHERE favorite.user_id =  " + fmt.Sprintf("%v", MyID) +
		" and favorite.video_id = video.id LIMIT 1) , false )  = 1"
	Db.Raw(VideoListSQL).Scan(&res)
	var VideoList = make([]common.Video, len(res))
	for i := 0; i < len(res); i++ {
		VideoList[i].Author.Id = res[i].AuthorId
		VideoList[i].Author.FollowCount = res[i].FollowCount
		VideoList[i].Author.FollowerCount = res[i].FollowerCount
		VideoList[i].Author.Name = res[i].Name
		VideoList[i].Author.IsFollow = res[i].IsFollow
		VideoList[i].Id = res[i].Id
		VideoList[i].PlayUrl = res[i].PlayUrl
		VideoList[i].CoverUrl = res[i].CoverUrl
		VideoList[i].IsFavorite = res[i].IsFavorite
		VideoList[i].Title = res[i].Title
		VideoList[i].CommentCount = res[i].CommentCount
		VideoList[i].FavoriteCount = res[i].FavoriteCount
	}
	return VideoList, nil
}

func (videoDao *VideoDao) Addvideo(authorId int64, playUrl string, coverUrl string, title string) error {
	newVideo := &Video{
		AuthorID: authorId,
		CoverUrl: "http://" + coverUrl,
		PlayUrl:  "http://" + playUrl,
		Title:    title,
	}
	// 重复视频校验（没做）
	err := Db.Transaction(func(tx *gorm.DB) error {
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
