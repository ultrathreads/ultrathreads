package model

// UserReadState 用户版块阅读状态表
type UserReadState struct {
	Model
	UserID     int64 `gorm:"not null;unique_index:uk_user_node" json:"userId"`
	NodeID     int64 `gorm:"not null;unique_index:uk_user_node" json:"nodeId"`
	LastReadAt int64 `gorm:"not null;default:0" json:"lastReadAt"` // Unix 秒级时间戳
}

func (UserReadState) TableName() string {
	return "user_read_states"
}