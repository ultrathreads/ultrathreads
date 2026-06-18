package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// LinkRepository 链接数据访问契约
type LinkRepository interface {
	Get(id int64) *model.Link
	Take(where ...interface{}) *model.Link
	Find(cnd *querybuilder.QueryBuilder) []model.Link
	FindOne(cnd *querybuilder.QueryBuilder) *model.Link
	List(cnd *querybuilder.QueryBuilder) ([]model.Link, *querybuilder.Paging)
	Count(cnd *querybuilder.QueryBuilder) int64
	Create(t *model.Link) error
	Update(t *model.Link) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64) error
}

type linkRepo struct {
	db *gorm.DB
}

func NewLinkDao(db *gorm.DB) LinkRepository {
	return &linkRepo{db: db}
}

func (r *linkRepo) Get(id int64) *model.Link {
	ret := &model.Link{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *linkRepo) Take(where ...interface{}) *model.Link {
	ret := &model.Link{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *linkRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.Link) {
	cnd.Find(r.db, &list)
	return
}

func (r *linkRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.Link {
	ret := &model.Link{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *linkRepo) List(cnd *querybuilder.QueryBuilder) (list []model.Link, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.Link{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *linkRepo) Count(cnd *querybuilder.QueryBuilder) int64 {
	return cnd.Count(r.db, &model.Link{})
}

func (r *linkRepo) Create(t *model.Link) error {
	return r.db.Create(t).Error
}

func (r *linkRepo) Update(t *model.Link) error {
	return r.db.Save(t).Error
}

func (r *linkRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.Link{}).Where("id = ?", id).Updates(columns).Error
}

func (r *linkRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.Link{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *linkRepo) Delete(id int64) error {
	return r.db.Delete(&model.Link{}, "id = ?", id).Error
}