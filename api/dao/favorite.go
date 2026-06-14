package dao

import (
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var FavoriteDao = newFavoriteDao()

func newFavoriteDao() *favoriteDao {
	return &favoriteDao{}
}

type favoriteDao struct {
}

func (d *favoriteDao) Get(id int64) *model.Favorite {
	ret := &model.Favorite{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *favoriteDao) Take(where ...interface{}) *model.Favorite {
	ret := &model.Favorite{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *favoriteDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Favorite) {
	cnd.Find(db, &list)
	return
}

func (d *favoriteDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Favorite {
	ret := &model.Favorite{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *favoriteDao) List(cnd *querybuilder.QueryBuilder) (list []model.Favorite, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Favorite{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total: count,
	}
	return
}

func (d *favoriteDao) Create(t *model.Favorite) (err error) {
	err = db.Create(t).Error
	return
}

func (d *favoriteDao) Update(t *model.Favorite) (err error) {
	err = db.Save(t).Error
	return
}

func (d *favoriteDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Favorite{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *favoriteDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Favorite{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *favoriteDao) Delete(id int64) {
	db.Delete(&model.Favorite{}, "id = ?", id)
}
