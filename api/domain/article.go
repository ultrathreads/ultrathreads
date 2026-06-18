package domain

import "time"

// Article 文章领域模型
type Article struct {
	ID          int64
	UserId      int64
	Title       string
	Summary     string
	Content     string
	ContentType string
	Status      int
	Share       bool
	SourceUrl   string
	ViewCount   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ArticleTag 文章标签
type ArticleTag struct {
	ID        int64
	ArticleId int64
	TagId     int64
	Status    int64
	CreatedAt time.Time
}
