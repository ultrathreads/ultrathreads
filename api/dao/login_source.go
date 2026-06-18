package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// LoginSourceRepository 登录来源数据访问契约
type LoginSourceRepository interface {
	Get(id int64) *model.LoginSource
	Take(where ...interface{}) *model.LoginSource
	Find(cnd *querybuilder.QueryBuilder) []model.LoginSource
	FindOne(cnd *querybuilder.QueryBuilder) *model.LoginSource
	List(cnd *querybuilder.QueryBuilder) ([]model.LoginSource, *querybuilder.Paging)
	Create(t *model.LoginSource) error
	Update(t *model.LoginSource) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
}

type loginSourceRepo struct {
	db *gorm.DB
}

func NewLoginSourceDao(db *gorm.DB) LoginSourceRepository {
	return &loginSourceRepo{db: db}
}

func (r *loginSourceRepo) Get(id int64) *model.LoginSource {
	ret := &model.LoginSource{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *loginSourceRepo) Take(where ...interface{}) *model.LoginSource {
	ret := &model.LoginSource{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *loginSourceRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.LoginSource) {
	cnd.Find(r.db, &list)
	return
}

func (r *loginSourceRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.LoginSource {
	ret := &model.LoginSource{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		return nil
	}
	return ret
}

func (r *loginSourceRepo) List(cnd *querybuilder.QueryBuilder) (list []model.LoginSource, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.LoginSource{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *loginSourceRepo) Create(t *model.LoginSource) error {
	return r.db.Create(t).Error
}

func (r *loginSourceRepo) Update(t *model.LoginSource) error {
	return r.db.Save(t).Error
}

func (r *loginSourceRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.LoginSource{}).Where("id = ?", id).Updates(columns).Error
}

func (r *loginSourceRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.LoginSource{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *loginSourceRepo) Delete(id int64) {
	r.db.Delete(&model.LoginSource{}, "id = ?", id)
}