package domain

// Role 角色领域模型
type Role struct {
	ID          int64
	Name        string
	DisplayName string
	Description string
	SortOrder   int
}

// Permission 权限领域模型
type Permission struct {
	ID          int64
	Code        string
	Name        string
	Module      string
	Description string
}

// UserRole 用户-角色关联
type UserRole struct {
	UserID int64
	RoleID int64
}

// RolePermission 角色-权限关联
type RolePermission struct {
	RoleID       int64
	PermissionID int64
}
