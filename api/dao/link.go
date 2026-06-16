package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewLinkDao(db *gorm.DB) *linkDao {
	return &linkDao{db: db}
}

type linkDao struct {
	db *gorm.DB
}

func (d *linkDao) Get(id int64) *model.Link {
	ret := &model.Link{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *linkDao) Take(where ...interface{}) *model.Link {
	ret := &model.Link{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *linkDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Link) {
	cnd.Find(d.db, &list)
	return
}

func (d *linkDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Link {
	ret := &model.Link{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *linkDao) List(cnd *querybuilder.QueryBuilder) (list []model.Link, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.Link{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *linkDao) Count(cnd *querybuilder.QueryBuilder) int64 {
	return cnd.Count(d.db, &model.Link{})
}

func (d *linkDao) Create(t *model.Link) error {
	return d.db.Create(t).Error
}

func (d *linkDao) Update(t *model.Link) error {
	return d.db.Save(t).Error
}

func (d *linkDao) Updates(id int64, columns map[string]interface{}) error {
	return d.db.Model(&model.Link{}).Where("id = ?", id).Updates(columns).Error
}

func (d *linkDao) UpdateColumn(id int64, name string, value interface{}) error {
	return d.db.Model(&model.Link{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (d *linkDao) Delete(id int64) error {
	return d.db.Delete(&model.Link{}, "id = ?", id).Error
}