package domain

// CreateLinkCommand 创建链接命令
type CreateLinkCommand struct {
	Title   string
	URL     string
	Logo    string
	Summary string
	Status  int
}

// UpdateLinkCommand 更新链接命令
type UpdateLinkCommand struct {
	ID      int64
	Title   string
	URL     string
	Logo    string
	Summary string
	Status  int
}
