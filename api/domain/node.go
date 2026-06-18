package domain

// Node 话题节点领域模型
type Node struct {
	ID          int64
	Name        string
	Description string
	Icon        string
	SortNo      int
	Status      int
	TopicCount  int64
	CreateTime  int64
}
