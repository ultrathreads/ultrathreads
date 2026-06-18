package domain

// CreatePostCommand 创建帖子命令
type CreatePostCommand struct {
	NodeSlug  string
	Title     string
	Content   string
	Tags      []string
	ImageList []string
}

// UpdatePostCommand 更新帖子命令
type UpdatePostCommand struct {
	Slug      string
	Title     *string
	Content   *string
	NodeSlug  *string
	Tags      []string
	ImageList []string
}

// CreateReplyCommand 创建回复命令
type CreateReplyCommand struct {
	Slug       string
	Content    string
	ImageList  []string
	ParentSlug string
}

// UpdateReplyCommand 更新回复命令
type UpdateReplyCommand struct {
	Slug    string
	Content *string
}
