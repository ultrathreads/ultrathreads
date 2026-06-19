package repository

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// UserRepository 用户数据访问契约
type UserRepository interface {
	Get(id int64) *model.User
	FindByIds(ids []int64) []model.User
	Take(where ...interface{}) *model.User
	Find(cnd *querybuilder.QueryBuilder) []model.User
	FindOne(cnd *querybuilder.QueryBuilder) *model.User
	List(cnd *querybuilder.QueryBuilder) ([]model.User, *querybuilder.Paging)
	Count(cnd *querybuilder.QueryBuilder) int64
	Create(t *model.User) error
	Update(t *model.User) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64) error
	GetByEmail(email string) *model.User
	GetByUsername(username string) *model.User
	// 以下方法封装事务/多表操作，避免 service 层依赖 *gorm.DB
	CreateWithAvatar(user *model.User, avatarUrl string, isAdmin bool) error
	CreateFromLoginSource(user *model.User, loginSourceId int64, avatarUrl string) error
	IncrColumn(id int64, column string, delta int) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Get(id int64) *model.User {
	ret := &model.User{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *userRepo) FindByIds(ids []int64) []model.User {
	if len(ids) == 0 {
		return nil
	}
	qb := querybuilder.NewQueryBuilder().In("id", ids)
	return r.Find(qb)
}

func (r *userRepo) Take(where ...interface{}) *model.User {
	ret := &model.User{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *userRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.User) {
	cnd.Find(r.db, &list)
	return
}

func (r *userRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.User {
	ret := &model.User{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *userRepo) List(cnd *querybuilder.QueryBuilder) (list []model.User, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.User{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *userRepo) Count(cnd *querybuilder.QueryBuilder) int64 {
	return cnd.Count(r.db, &model.User{})
}

func (r *userRepo) Create(t *model.User) error {
	return r.db.Create(t).Error
}

func (r *userRepo) Update(t *model.User) error {
	return r.db.Save(t).Error
}

func (r *userRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(columns).Error
}

func (r *userRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *userRepo) Delete(id int64) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

func (r *userRepo) GetByEmail(email string) *model.User {
	return r.Take("email = ?", email)
}

func (r *userRepo) GetByUsername(username string) *model.User {
	return r.Take("username = ?", username)
}

// CreateWithAvatar 创建用户并更新头像（事务：创建用户 + 更新头像/管理员等级）
func (r *userRepo) CreateWithAvatar(user *model.User, avatarUrl string, isAdmin bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		updateColumns := map[string]interface{}{
			"avatar": avatarUrl,
		}
		if isAdmin {
			updateColumns["level"] = model.UserLevelAdmin
		}
		return tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(updateColumns).Error
	})
}

// CreateFromLoginSource 通过登录来源创建用户（事务：创建用户 + 绑定登录来源 + 更新头像）
func (r *userRepo) CreateFromLoginSource(user *model.User, loginSourceId int64, avatarUrl string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.LoginSource{}).Where("id = ?", loginSourceId).
			UpdateColumn("user_id", user.ID).Error; err != nil {
			return err
		}
		return tx.Model(&model.User{}).Where("id = ?", user.ID).
			UpdateColumn("avatar", avatarUrl).Error
	})
}

// IncrColumn 自增指定字段
func (r *userRepo) IncrColumn(id int64, column string, delta int) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn(column, gorm.Expr(column+" + ?", delta)).Error
}
