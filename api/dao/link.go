package dao

import (
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var LinkDao = newLinkDao()

func newLinkDao() *linkDao {
	return &linkDao{}
}

type linkDao struct {
}

func (d *linkDao) Get(id int64) *model.Link {
	ret := &model.Link{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *linkDao) Take(where ...interface{}) *model.Link {
	ret := &model.Link{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *linkDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Link) {
	cnd.Find(db, &list)
	return
}

func (d *linkDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Link {
	ret := &model.Link{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *linkDao) List(cnd *querybuilder.QueryBuilder) (list []model.Link, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Link{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *linkDao) Create(t *model.Link) (err error) {
	err = db.Create(t).Error
	return
}

func (d *linkDao) Update(t *model.Link) (err error) {
	err = db.Save(t).Error
	return
}

func (d *linkDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Link{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *linkDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Link{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *linkDao) Delete(id int64) {
	db.Delete(&model.Link{}, "id = ?", id)
}
