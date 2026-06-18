package mock

import (
	"fmt"

	"ultrathreads/model"
)

// 预置角色定义
var seedRoles = []model.Role{
	{Name: "admin", DisplayName: "超级管理员", Description: "拥有系统所有权限", SortOrder: 1},
	{Name: "editor", DisplayName: "内容编辑", Description: "管理文章与帖子内容", SortOrder: 2},
	{Name: "user", DisplayName: "普通用户", Description: "基础前台访问权限", SortOrder: 3},
}

// 预置权限定义（按模块划分）
var seedPermissions = []model.Permission{
	// system 模块
	{Code: "admin:panel:access", Name: "后台管理准入", Module: "system", Description: "允许进入后台管理面板"}, // ✅ 新增
	{Code: "dashboard:view", Name: "查看仪表盘", Module: "system", Description: "访问后台首页"},
	{Code: "site:config", Name: "站点配置", Module: "system", Description: "修改站点基础设置"},
	// post 模块
	{Code: "post:create", Name: "创建文章", Module: "post", Description: "发布新文章/帖子"},
	{Code: "post:edit", Name: "编辑文章", Module: "post", Description: "编辑任意文章/帖子"},
	{Code: "post:delete", Name: "删除文章", Module: "post", Description: "删除任意文章/帖子"},
	// user 模块
	{Code: "user:manage", Name: "管理用户", Module: "user", Description: "封禁/解封/编辑用户信息"},
}

// rolePermMapping 定义每个角色拥有的权限 Code
var rolePermMapping = map[string][]string{
	"admin": {
		"admin:panel:access",
		"dashboard:view", "site:config",
		"post:create", "post:edit", "post:delete",
		"user:manage",
	},
	"editor": {
		"dashboard:view",
		"post:create", "post:edit", "post:delete",
	},
	"user": {
		"post:create",
	},
}

// RbacTableSeeder - 初始化 RBAC 相关表数据
func RbacTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.RolePermission{})
		dropAndCreateTable(&model.UserRole{})
		dropAndCreateTable(&model.Permission{})
		dropAndCreateTable(&model.Role{})
	}

	// ========== 1. 插入角色 ==========
	roleMap := make(map[string]*model.Role)
	for i := range seedRoles {
		role := &seedRoles[i]
		if err := rbacDao.CreateRole(role); err != nil {
			fmt.Printf("mock role error [%s]: %v\n", role.Name, err)
			continue
		}
		roleMap[role.Name] = role
		fmt.Println("✅ Created role:", role.DisplayName)
	}

	// ========== 2. 插入权限 ==========
	permMap := make(map[string]*model.Permission)
	for i := range seedPermissions {
		perm := &seedPermissions[i]
		if err := rbacDao.CreatePermission(perm); err != nil {
			fmt.Printf("mock permission error [%s]: %v\n", perm.Code, err)
			continue
		}
		permMap[perm.Code] = perm
		fmt.Println("✅ Created permission:", perm.Code)
	}

	// ========== 3. 绑定 角色-权限 ==========
	for roleName, permCodes := range rolePermMapping {
		role, ok := roleMap[roleName]
		if !ok {
			fmt.Printf("⚠️  skip role-perm binding: role [%s] not found\n", roleName)
			continue
		}
		boundCount := 0
		for _, code := range permCodes {
			perm, ok := permMap[code]
			if !ok {
				fmt.Printf("⚠️  skip role-perm binding: permission [%s] not found\n", code)
				continue
			}
			if err := rbacDao.AssignPermissionToRole(role.ID, perm.ID); err != nil {
				fmt.Printf("mock role_permission error [%s->%s]: %v\n", roleName, code, err)
				continue
			}
			boundCount++
		}
		fmt.Printf("✅ Bound %d permissions to role [%s]\n", boundCount, roleName)
	}

	// ========== 4. 绑定 用户-角色 ==========
	// admin(i=0) -> ID=1, ultrathreads(i=1) -> ID=2
	userRoleBindings := []struct {
		UserID   int64
		RoleName string
	}{
		{UserID: 1, RoleName: "admin"},
		{UserID: 2, RoleName: "editor"},
	}

	for _, binding := range userRoleBindings {
		role, ok := roleMap[binding.RoleName]
		if !ok {
			fmt.Printf("⚠️  skip user-role binding: role [%s] not found\n", binding.RoleName)
			continue
		}
		if err := rbacDao.AssignRoleToUser(binding.UserID, role.ID); err != nil {
			fmt.Printf("mock user_role error [user:%d->%s]: %v\n", binding.UserID, binding.RoleName)
		} else {
			fmt.Printf("✅ Bound user [ID:%d] to role [%s]\n", binding.UserID, binding.RoleName)
		}
	}
}
