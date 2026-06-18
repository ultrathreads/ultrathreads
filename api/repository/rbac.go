package repository

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// RbacRepository RBAC 数据访问契约
type RbacRepository interface {
	GetRole(id int64) *model.Role
	GetRoleByName(name string) *model.Role
	FindRoles(cnd *querybuilder.QueryBuilder) []model.Role
	ListRoles(cnd *querybuilder.QueryBuilder) ([]model.Role, *querybuilder.Paging)
	CreateRole(t *model.Role) error
	UpdateRole(t *model.Role) error
	DeleteRole(id int64) error
	GetUserRoleCodes(userID int64) []string

	GetPermission(id int64) *model.Permission
	GetPermissionByCode(code string) *model.Permission
	FindPermissions(cnd *querybuilder.QueryBuilder) []model.Permission
	CreatePermission(t *model.Permission) error
	UpdatePermission(t *model.Permission) error
	DeletePermission(id int64) error

	AssignRoleToUser(userID, roleID int64) error
	RevokeRoleFromUser(userID, roleID int64) error
	GetUserRoleIDs(userID int64) []int64

	AssignPermissionToRole(roleID, permID int64) error
	RevokePermissionFromRole(roleID, permID int64) error
	GetRolePermissionCodes(roleID int64) []string
	GetUserPermissionCodes(userID int64) []string
}

type rbacRepo struct {
	db *gorm.DB
}

func NewRbacRepository(db *gorm.DB) RbacRepository {
	return &rbacRepo{db: db}
}

// ==================== Role ====================

func (r *rbacRepo) GetRole(id int64) *model.Role {
	ret := &model.Role{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *rbacRepo) GetRoleByName(name string) *model.Role {
	ret := &model.Role{}
	if err := r.db.Take(ret, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *rbacRepo) FindRoles(cnd *querybuilder.QueryBuilder) (list []model.Role) {
	cnd.Find(r.db, &list)
	return
}

func (r *rbacRepo) ListRoles(cnd *querybuilder.QueryBuilder) (list []model.Role, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.Role{})
	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *rbacRepo) CreateRole(t *model.Role) error {
	return r.db.Create(t).Error
}

func (r *rbacRepo) UpdateRole(t *model.Role) error {
	return r.db.Save(t).Error
}

func (r *rbacRepo) DeleteRole(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.UserRole{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.RolePermission{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Role{}, "id = ?", id).Error
	})
}

func (r *rbacRepo) GetUserRoleCodes(userID int64) []string {
	var codes []string

	roleStmt := &gorm.Statement{DB: r.db}
	_ = roleStmt.Parse(&model.Role{})
	roleTable := roleStmt.Table

	urStmt := &gorm.Statement{DB: r.db}
	_ = urStmt.Parse(&model.UserRole{})
	urTable := urStmt.Table

	r.db.Table(roleTable).
		Select("DISTINCT "+roleTable+".name").
		Joins("JOIN "+urTable+" ON "+urTable+".role_id = "+roleTable+".id").
		Where(urTable+".user_id = ?", userID).
		Pluck("name", &codes)

	if codes == nil {
		return []string{}
	}
	return codes
}

// ==================== Permission ====================

func (r *rbacRepo) GetPermission(id int64) *model.Permission {
	ret := &model.Permission{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *rbacRepo) GetPermissionByCode(code string) *model.Permission {
	ret := &model.Permission{}
	if err := r.db.Take(ret, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *rbacRepo) FindPermissions(cnd *querybuilder.QueryBuilder) (list []model.Permission) {
	cnd.Find(r.db, &list)
	return
}

func (r *rbacRepo) CreatePermission(t *model.Permission) error {
	return r.db.Create(t).Error
}

func (r *rbacRepo) UpdatePermission(t *model.Permission) error {
	return r.db.Save(t).Error
}

func (r *rbacRepo) DeletePermission(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.RolePermission{}, "permission_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Permission{}, "id = ?", id).Error
	})
}

// ==================== UserRole 关联 ====================

func (r *rbacRepo) AssignRoleToUser(userID, roleID int64) error {
	ur := &model.UserRole{UserID: userID, RoleID: roleID}
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).
		FirstOrCreate(ur).Error
}

func (r *rbacRepo) RevokeRoleFromUser(userID, roleID int64) error {
	return r.db.Delete(&model.UserRole{}, "user_id = ? AND role_id = ?", userID, roleID).Error
}

func (r *rbacRepo) GetUserRoleIDs(userID int64) []int64 {
	var roleIDs []int64
	r.db.Table("user_roles").
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs)
	return roleIDs
}

// ==================== RolePermission 关联 ====================

func (r *rbacRepo) AssignPermissionToRole(roleID, permID int64) error {
	rp := &model.RolePermission{RoleID: roleID, PermissionID: permID}
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permID).
		FirstOrCreate(rp).Error
}

func (r *rbacRepo) RevokePermissionFromRole(roleID, permID int64) error {
	return r.db.Delete(&model.RolePermission{}, "role_id = ? AND permission_id = ?", roleID, permID).Error
}

func (r *rbacRepo) GetRolePermissionCodes(roleID int64) []string {
	var codes []string

	permStmt := &gorm.Statement{DB: r.db}
	_ = permStmt.Parse(&model.Permission{})
	permTable := permStmt.Table

	rpStmt := &gorm.Statement{DB: r.db}
	_ = rpStmt.Parse(&model.RolePermission{})
	rpTable := rpStmt.Table

	r.db.Table(permTable).
		Select(permTable+".code").
		Joins("JOIN "+rpTable+" ON "+rpTable+".permission_id = "+permTable+".id").
		Where(rpTable+".role_id = ?", roleID).
		Pluck("code", &codes)

	return codes
}

func (r *rbacRepo) GetUserPermissionCodes(userID int64) []string {
	var codes []string

	permStmt := &gorm.Statement{DB: r.db}
	_ = permStmt.Parse(&model.Permission{})
	permTable := permStmt.Table

	rpStmt := &gorm.Statement{DB: r.db}
	_ = rpStmt.Parse(&model.RolePermission{})
	rpTable := rpStmt.Table

	urStmt := &gorm.Statement{DB: r.db}
	_ = urStmt.Parse(&model.UserRole{})
	urTable := urStmt.Table

	r.db.Table(permTable).
		Select("DISTINCT "+permTable+".code").
		Joins("JOIN "+rpTable+" ON "+rpTable+".permission_id = "+permTable+".id").
		Joins("JOIN "+urTable+" ON "+urTable+".role_id = "+rpTable+".role_id").
		Where(urTable+".user_id = ?", userID).
		Pluck("code", &codes)

	if codes == nil {
		return []string{}
	}
	return codes
}
