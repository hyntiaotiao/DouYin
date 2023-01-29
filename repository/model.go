package repository

import (
	"time"
)

// User 用户数据表
type User struct {
	ID            int64     `gorm:"column:id;primary_key" json:"id"`                                          // 主键
	Username      string    `gorm:"column:username;NOT NULL" json:"username"`                                 // 用户名称
	Password      string    `gorm:"column:password;NOT NULL" json:"password"`                                 // 密码
	Gender        int       `gorm:"column:gender;NOT NULL" json:"gender"`                                     // 男1 女0
	FollowerCount int       `gorm:"column:follower_count;NOT NULL" json:"follower_count"`                     // 粉丝数
	FollowCount   int       `gorm:"column:follow_count;NOT NULL" json:"follow_count"`                         // 关注数
	Face          string    `gorm:"column:face;NOT NULL" json:"face"`                                         // 头像
	CreateTime    time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"` // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDeleted     int       `gorm:"column:is_deleted;default:0;NOT NULL" json:"is_deleted"`                   // 是否删除(0-未删, 1-已删)
}

// TableName 去除默认的表名复数，以user作为表名
func (User) TableName() string {
	return "user"
}

// UserVideo 用户视频表
type UserVideo struct {
	ID         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`                           // 主键
	UserID     int64     `gorm:"column:user_id;NOT NULL" json:"user_id"`                                   // 用户id
	VideoID    int64     `gorm:"column:video_id;NOT NULL" json:"video_id"`                                 // 视频id
	IsFavorite int       `gorm:"column:is_favorite;NOT NULL" json:"is_favorite"`                           // 1表示已点赞 0未点赞
	IsFollow   int       `gorm:"column:is_follow;NOT NULL" json:"is_follow"`                               // 是否关注该视频作者
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDeleted  int       `gorm:"column:is_deleted;default:0;NOT NULL" json:"is_deleted"`                   // 是否删除(0-未删, 1-已删)
}

// TableName 去除默认的表名复数，以user作为表名
func (UserVideo) TableName() string {
	return "user_video"
}

type Video struct {
	ID             int64     `gorm:"column:id;primary_key" json:"id"`                                          // 视频id
	AuthorID       string    `gorm:"column:author_id;NOT NULL" json:"author_id"`                               // 视频上传的作者ID
	Description    string    `gorm:"column:description" json:"description"`                                    // 视频简介
	PlayUrl        string    `gorm:"column:play_url;NOT NULL" json:"play_url"`                                 // 视频播放地址
	CoverUrl       string    `gorm:"column:cover_url;NOT NULL" json:"cover_url"`                               // 视频封面地址
	Title          string    `gorm:"column:title;NOT NULL" json:"title"`                                       // 标题
	FavouriteCount int       `gorm:"column:favourite_count;NOT NULL" json:"favourite_count"`                   // 点赞量
	PlayCounts     int       `gorm:"column:play_counts;NOT NULL" json:"play_counts"`                           // 播放次数
	CommentCount   int       `gorm:"column:comment_count;NOT NULL" json:"comment_count"`                       // 评论量
	CreateTime     time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"` // 创建时间
	UpdateTime     time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDeleted      int       `gorm:"column:is_deleted;default:0;NOT NULL" json:"is_deleted"`                   // 是否删除(0-未删, 1-已删)
}

// Fans 粉丝表
type Fans struct {
	ID         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`                           // 主键
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDeleted  int       `gorm:"column:is_deleted;default:0;NOT NULL" json:"is_deleted"`                   // 是否删除(0-未删, 1-已删)
	BloggerID  int64     `gorm:"column:blogger_id;NOT NULL" json:"blogger_id"`                             // 博主id
	FansID     int64     `gorm:"column:fans_id;NOT NULL" json:"fans_id"`                                   // 粉丝id
}
