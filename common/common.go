package common

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type VideoVO struct {
	Id            int64  `json:"id"`
	Author        UserVO `json:"author" gorm:"references:ID"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type CommentVO struct {
	Id         int64  `json:"id"`
	User       UserVO `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type UserVO struct {
	Id            int64  `json:"id" gorm:""`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type MessageVO struct {
	Id         int64  `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
	FromUserId int64  `json:"from_user_id"`
	ToUserId   int64  `json:"to_user_id"`
}
