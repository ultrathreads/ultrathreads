package model

// 话题节点
type Node struct {
	Model
	Name        string `gorm:"size:32;unique" json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	SortNo      int    `gorm:"index:idx_sort_no" json:"sortNo" form:"sortNo"`
	Status      int    `gorm:"not null" json:"status" form:"status"`
	TopicCount   int64  `gorm:"not null" json:"topicCount" form:"topicCount"`
	CreateTime  int64  `json:"createTime" form:"createTime"`
}
