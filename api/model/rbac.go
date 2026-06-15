package model

// Role 角色表
type Role struct {
    Model
    Name        string `gorm:"size:32;uniqueIndex;not null" json:"name"` // admin, editor, user
    DisplayName string `gorm:"size:64;not null" json:"displayName"`      // 管理员, 编辑, 普通用户
    Description string `gorm:"size:256" json:"description"`
    SortOrder   int    `gorm:"not null;default:0" json:"sortOrder"`
}

// Permission 权限表
type Permission struct {
    Model
    Code        string `gorm:"size:64;uniqueIndex;not null" json:"code"` // post:create, user:ban
    Name        string `gorm:"size:64;not null" json:"name"`             // 创建文章, 封禁用户
    Module      string `gorm:"size:32;index;not null" json:"module"`     // post, user, system
    Description string `gorm:"size:256" json:"description"`
}

// UserRole 用户-角色关联表（多对多）
type UserRole struct {
    UserID int64 `gorm:"primaryKey" json:"userId"`
    RoleID int64 `gorm:"primaryKey" json:"roleId"`
}

// RolePermission 角色-权限关联表（多对多）
type RolePermission struct {
    RoleID       int64 `gorm:"primaryKey" json:"roleId"`
    PermissionID int64 `gorm:"primaryKey" json:"permissionId"`
}
