package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewLoginSourceDao(db *gorm.DB) *loginSourceDao {
	return &loginSourceDao{db: db}
}

type loginSourceDao struct {
	db *gorm.DB
}

func (d *loginSourceDao) Get(id int64) *model.LoginSource {
	ret := &model.LoginSource{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *loginSourceDao) Take(where ...interface{}) *model.LoginSource {
	ret := &model.LoginSource{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *loginSourceDao) List(cnd *querybuilder.QueryBuilder) (list []model.LoginSource, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.LoginSource{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *loginSourceDao) Create(t *model.LoginSource) (err error) {
	err = d.db.Create(t).Error
	return
}

func (d *loginSourceDao) Update(t *model.LoginSource) (err error) {
	err = d.db.Save(t).Error
	return
}

func (d *loginSourceDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = d.db.Model(&model.LoginSource{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *loginSourceDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = d.db.Model(&model.LoginSource{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *loginSourceDao) Delete(id int64) {
	d.db.Delete(&model.LoginSource{}, "id = ?", id)
}