package dao

import (
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var PostLikeDao = newPostLikeDao()

func newPostLikeDao() *postLikeDao {
	return &postLikeDao{}
}

type postLikeDao struct {
}

func (d *postLikeDao) Get(id int64) *model.PostLike {
	ret := &model.PostLike{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *postLikeDao) Take(where ...interface{}) *model.PostLike {
	ret := &model.PostLike{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *postLikeDao) Find(cnd *querybuilder.QueryBuilder) (list []model.PostLike) {
	cnd.Find(db, &list)
	return
}

func (d *postLikeDao) FindOne(cnd *querybuilder.QueryBuilder) *model.PostLike {
	ret := &model.PostLike{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *postLikeDao) List(cnd *querybuilder.QueryBuilder) (list []model.PostLike, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.PostLike{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *postLikeDao) Create(t *model.PostLike) (err error) {
	err = db.Create(t).Error
	return
}

func (d *postLikeDao) Update(t *model.PostLike) (err error) {
	err = db.Save(t).Error
	return
}

func (d *postLikeDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.PostLike{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *postLikeDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.PostLike{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *postLikeDao) Delete(id int64) {
	db.Delete(&model.PostLike{}, "id = ?", id)
}
