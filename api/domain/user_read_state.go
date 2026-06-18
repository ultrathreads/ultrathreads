package domain

// UserReadState 用户版块阅读状态领域模型
type UserReadState struct {
	ID         int64
	UserID     int64
	NodeID     int64
	LastReadAt int64
}
