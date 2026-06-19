package domain

// UpdateUserCommand 更新用户命令
type UpdateUserCommand struct {
	Slug        string
	Nickname    string
	Avatar      string
	Website     string
	Description string
	Level       int
}
