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
