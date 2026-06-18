package domain

import "time"

// Post 话题领域模型
type Post struct {
	ID                int64
	ThreadId          int64
	ParentId          int64
	Type              int
	NodeId            int64
	UserId            int64
	Title             string
	Content           string
	ImageList         string
	IsPinned          bool
	Recommend         bool
	ViewCount         int64
	LikeCount         int64
	Status            int
	LastCommentUserId int64
	LastCommentTime   int64
	CreateTime        int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ExtraData         string
}

// IsRoot 判断是否是根帖子
func (p *Post) IsRoot() bool {
	return p.ParentId == 0
}

// PostTag 主题标签
type PostTag struct {
	ID              int64
	PostId          int64
	TagId           int64
	Status          int64
	LastCommentTime int64
	CreateTime      int64
}

// PostLike 话题点赞
type PostLike struct {
	ID         int64
	UserId     int64
	PostId     int64
	CreateTime int64
}
