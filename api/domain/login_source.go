package domain

// LoginSource 第三方登录源领域模型
type LoginSource struct {
	ID         int64
	UserID     int64
	Avatar     string
	Nickname   string
	TargetType string
	TargetID   string
	ExtraData  string
	CreateTime int64
	UpdateTime int64
}
