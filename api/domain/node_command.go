package domain

// CreateNodeCommand 创建节点命令
type CreateNodeCommand struct {
	Name        string
	Description string
	Icon        string
	SortNo      int
	Status      int
}

// UpdateNodeCommand 更新节点命令
type UpdateNodeCommand struct {
	ID          int64
	Name        string
	Description string
	Icon        string
	SortNo      int
	Status      int
}
