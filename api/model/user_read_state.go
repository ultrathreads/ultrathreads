package model

// UserReadState 用户版块阅读状态表
type UserReadState struct {
	Model
	UserID     int64 `gorm:"not null;uniqueIndex:uk_user_node,priority:1" json:"userId"`
	NodeID     int64 `gorm:"not null;uniqueIndex:uk_user_node,priority:2" json:"nodeId"`
	LastReadAt int64 `gorm:"not null;default:0" json:"lastReadAt"`
}