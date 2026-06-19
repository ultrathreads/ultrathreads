package domain

// CreateTagCommand 创建标签命令
type CreateTagCommand struct {
	Name        string
	Description string
	Status      int
}

// UpdateTagCommand 更新标签命令
type UpdateTagCommand struct {
	ID          int64
	Name        string
	Description string
	Status      int
}
