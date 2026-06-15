package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var RbacDao = newRbacDao()

func newRbacDao() *rbacDao {
	return &rbacDao{}
}

type rbacDao struct{}

// ==================== Role ====================

func (d *rbacDao) GetRole(id int64) *model.Role {
	ret := &model.Role{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *rbacDao) GetRoleByName(name string) *model.Role {
	ret := &model.Role{}
	if err := db.Take(ret, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *rbacDao) FindRoles(cnd *querybuilder.QueryBuilder) (list []model.Role) {
	cnd.Find(db, &list)
	return
}

func (d *rbacDao) ListRoles(cnd *querybuilder.QueryBuilder) (list []model.Role, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Role{})
	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *rbacDao) CreateRole(t *model.Role) error {
	return db.Create(t).Error
}

func (d *rbacDao) UpdateRole(t *model.Role) error {
	return db.Save(t).Error
}

func (d *rbacDao) DeleteRole(id int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 删除角色时级联清理关联表
		if err := tx.Delete(&model.UserRole{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.RolePermission{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Role{}, "id = ?", id).Error
	})
}

// ==================== Permission ====================

func (d *rbacDao) GetPermission(id int64) *model.Permission {
	ret := &model.Permission{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *rbacDao) GetPermissionByCode(code string) *model.Permission {
	ret := &model.Permission{}
	if err := db.Take(ret, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *rbacDao) FindPermissions(cnd *querybuilder.QueryBuilder) (list []model.Permission) {
	cnd.Find(db, &list)
	return
}

func (d *rbacDao) CreatePermission(t *model.Permission) error {
	return db.Create(t).Error
}

func (d *rbacDao) UpdatePermission(t *model.Permission) error {
	return db.Save(t).Error
}

func (d *rbacDao) DeletePermission(id int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.RolePermission{}, "permission_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Permission{}, "id = ?", id).Error
	})
}

// ==================== UserRole 关联 ====================

func (d *rbacDao) AssignRoleToUser(userID, roleID int64) error {
	ur := &model.UserRole{UserID: userID, RoleID: roleID}
	return db.Where("user_id = ? AND role_id = ?", userID, roleID).
		FirstOrCreate(ur).Error
}

func (d *rbacDao) RevokeRoleFromUser(userID, roleID int64) error {
	return db.Delete(&model.UserRole{}, "user_id = ? AND role_id = ?", userID, roleID).Error
}

func (d *rbacDao) GetUserRoleIDs(userID int64) []int64 {
	var roleIDs []int64
	db.Table("user_roles").
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs)
	return roleIDs
}

// ==================== RolePermission 关联 ====================

func (d *rbacDao) AssignPermissionToRole(roleID, permID int64) error {
	rp := &model.RolePermission{RoleID: roleID, PermissionID: permID}
	return db.Where("role_id = ? AND permission_id = ?", roleID, permID).
		FirstOrCreate(rp).Error
}

func (d *rbacDao) RevokePermissionFromRole(roleID, permID int64) error {
	return db.Delete(&model.RolePermission{}, "role_id = ? AND permission_id = ?", roleID, permID).Error
}

func (d *rbacDao) GetRolePermissionCodes(roleID int64) []string {
	var codes []string
	db.Table("permissions").
		Select("permissions.code").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Pluck("code", &codes)
	return codes
}

func (d *rbacDao) GetUserPermissionCodes(userID int64) []string {
	var codes []string
	db.Table("permissions").
		Select("DISTINCT permissions.code").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("code", &codes)
	return codes
}