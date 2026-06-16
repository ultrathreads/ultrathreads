package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewFavoriteDao(db *gorm.DB) *favoriteDao {
	return &favoriteDao{db: db}
}

type favoriteDao struct {
	db *gorm.DB
}

func (d *favoriteDao) Get(id int64) *model.Favorite {
	ret := &model.Favorite{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *favoriteDao) Take(where ...interface{}) *model.Favorite {
	ret := &model.Favorite{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *favoriteDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Favorite) {
	cnd.Find(d.db, &list)
	return
}

func (d *favoriteDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Favorite {
	ret := &model.Favorite{}
	if err := cnd.FindOne(d.db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *favoriteDao) List(cnd *querybuilder.QueryBuilder) (list []model.Favorite, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.Favorite{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *favoriteDao) Create(t *model.Favorite) (err error) {
	err = d.db.Create(t).Error
	return
}

func (d *favoriteDao) Update(t *model.Favorite) (err error) {
	err = d.db.Save(t).Error
	return
}

func (d *favoriteDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = d.db.Model(&model.Favorite{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *favoriteDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = d.db.Model(&model.Favorite{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *favoriteDao) Delete(id int64) {
	d.db.Delete(&model.Favorite{}, "id = ?", id)
}