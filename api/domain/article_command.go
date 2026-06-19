package domain

// CreateArticleCommand 创建文章命令
type CreateArticleCommand struct {
	UserID  int64
	Title   string
	Summary string
	Content string
	Tags    string
}

// UpdateArticleCommand 更新文章命令
type UpdateArticleCommand struct {
	ID      int64
	Title   string
	Summary string
	Content string
	Tags    string
}
