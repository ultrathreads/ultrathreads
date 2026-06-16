package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewRbacDao(db *gorm.DB) *rbacDao {
	return &rbacDao{db: db}
}

type rbacDao struct {
	db *gorm.DB
}

// ==================== Role ====================

func (d *rbacDao) GetRole(id int64) *model.Role {
	ret := &model.Role{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *rbacDao) GetRoleByName(name string) *model.Role {
	ret := &model.Role{}
	if err := d.db.Take(ret, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *rbacDao) FindRoles(cnd *querybuilder.QueryBuilder) (list []model.Role) {
	cnd.Find(d.db, &list)
	return
}

func (d *rbacDao) ListRoles(cnd *querybuilder.QueryBuilder) (list []model.Role, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.Role{})
	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *rbacDao) CreateRole(t *model.Role) error {
	return d.db.Create(t).Error
}

func (d *rbacDao) UpdateRole(t *model.Role) error {
	return d.db.Save(t).Error
}

func (d *rbacDao) DeleteRole(id int64) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
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

// GetUserRoleCodes 获取用户所有角色标识（name）
func (d *rbacDao) GetUserRoleCodes(userID int64) []string {
	var codes []string

	roleStmt := &gorm.Statement{DB: d.db}
	_ = roleStmt.Parse(&model.Role{})
	roleTable := roleStmt.Table

	urStmt := &gorm.Statement{DB: d.db}
	_ = urStmt.Parse(&model.UserRole{})
	urTable := urStmt.Table

	d.db.Table(roleTable).
		Select("DISTINCT " + roleTable + ".name").
		Joins("JOIN "+urTable+" ON "+urTable+".role_id = "+roleTable+".id").
		Where(urTable+".user_id = ?", userID).
		Pluck("name", &codes)

	// 保证返回非 nil，与 GetUserPermissionCodes 行为一致
	if codes == nil {
		return []string{}
	}
	return codes
}

// ==================== Permission ====================

func (d *rbacDao) GetPermission(id int64) *model.Permission {
	ret := &model.Permission{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *rbacDao) GetPermissionByCode(code string) *model.Permission {
	ret := &model.Permission{}
	if err := d.db.Take(ret, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *rbacDao) FindPermissions(cnd *querybuilder.QueryBuilder) (list []model.Permission) {
	cnd.Find(d.db, &list)
	return
}

func (d *rbacDao) CreatePermission(t *model.Permission) error {
	return d.db.Create(t).Error
}

func (d *rbacDao) UpdatePermission(t *model.Permission) error {
	return d.db.Save(t).Error
}

func (d *rbacDao) DeletePermission(id int64) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.RolePermission{}, "permission_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Permission{}, "id = ?", id).Error
	})
}

// ==================== UserRole 关联 ====================

func (d *rbacDao) AssignRoleToUser(userID, roleID int64) error {
	ur := &model.UserRole{UserID: userID, RoleID: roleID}
	return d.db.Where("user_id = ? AND role_id = ?", userID, roleID).
		FirstOrCreate(ur).Error
}

func (d *rbacDao) RevokeRoleFromUser(userID, roleID int64) error {
	return d.db.Delete(&model.UserRole{}, "user_id = ? AND role_id = ?", userID, roleID).Error
}

func (d *rbacDao) GetUserRoleIDs(userID int64) []int64 {
	var roleIDs []int64
	d.db.Table("user_roles").
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs)
	return roleIDs
}

// ==================== RolePermission 关联 ====================

func (d *rbacDao) AssignPermissionToRole(roleID, permID int64) error {
	rp := &model.RolePermission{RoleID: roleID, PermissionID: permID}
	return d.db.Where("role_id = ? AND permission_id = ?", roleID, permID).
		FirstOrCreate(rp).Error
}

func (d *rbacDao) RevokePermissionFromRole(roleID, permID int64) error {
	return d.db.Delete(&model.RolePermission{}, "role_id = ? AND permission_id = ?", roleID, permID).Error
}

func (d *rbacDao) GetRolePermissionCodes(roleID int64) []string {
	var codes []string

	// 通过 Statement 解析出带前缀的真实表名
	permStmt := &gorm.Statement{DB: d.db}
	_ = permStmt.Parse(&model.Permission{})
	permTable := permStmt.Table

	rpStmt := &gorm.Statement{DB: d.db}
	_ = rpStmt.Parse(&model.RolePermission{})
	rpTable := rpStmt.Table

	d.db.Table(permTable).
		Select(permTable + ".code").
		Joins("JOIN "+rpTable+" ON "+rpTable+".permission_id = "+permTable+".id").
		Where(rpTable+".role_id = ?", roleID).
		Pluck("code", &codes)
	return codes
}

func (d *rbacDao) GetUserPermissionCodes(userID int64) []string {
	var codes []string

	permStmt := &gorm.Statement{DB: d.db}
	_ = permStmt.Parse(&model.Permission{})
	permTable := permStmt.Table

	rpStmt := &gorm.Statement{DB: d.db}
	_ = rpStmt.Parse(&model.RolePermission{})
	rpTable := rpStmt.Table

	urStmt := &gorm.Statement{DB: d.db}
	_ = urStmt.Parse(&model.UserRole{})
	urTable := urStmt.Table

	d.db.Table(permTable).
		Select("DISTINCT " + permTable + ".code").
		Joins("JOIN "+rpTable+" ON "+rpTable+".permission_id = "+permTable+".id").
		Joins("JOIN "+urTable+" ON "+urTable+".role_id = "+rpTable+".role_id").
		Where(urTable+".user_id = ?", userID).
		Pluck("code", &codes)
	return codes
}