package dao

import (
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var PostDao = newPostDao()

func newPostDao() *postDao {
	return &postDao{}
}

type postDao struct {
}

func (d *postDao) Get(id int64) *model.Post {
	ret := &model.Post{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *postDao) Take(where ...interface{}) *model.Post {
	ret := &model.Post{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *postDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Post) {
	cnd.Find(db, &list)
	return
}

func (d *postDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Post {
	ret := &model.Post{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *postDao) List(cnd *querybuilder.QueryBuilder) (list []model.Post, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Post{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *postDao) Count(cnd *querybuilder.QueryBuilder) int {
	return cnd.Count(db, &model.Post{})
}

func (d *postDao) Create(t *model.Post) (err error) {
	err = db.Create(t).Error
	return
}

func (d *postDao) Update(t *model.Post) (err error) {
	err = db.Save(t).Error
	return
}

func (d *postDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Post{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *postDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Post{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *postDao) Delete(id int64) {
	db.Delete(&model.Post{}, "id = ?", id)
}
