package dao

import (
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var CommentDao = newCommentDao()

func newCommentDao() *commentDao {
	return &commentDao{}
}

type commentDao struct {
}

func (d *commentDao) Get(id int64) *model.Comment {
	ret := &model.Comment{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *commentDao) Take(where ...interface{}) *model.Comment {
	ret := &model.Comment{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *commentDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Comment) {
	cnd.Find(db, &list)
	return
}

func (d *commentDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Comment {
	ret := &model.Comment{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *commentDao) List(cnd *querybuilder.QueryBuilder) (list []model.Comment, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Comment{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *commentDao) Count(cnd *querybuilder.QueryBuilder) int {
	return cnd.Count(db, &model.Comment{})
}

func (d *commentDao) Create(t *model.Comment) (err error) {
	err = db.Create(t).Error
	return
}

func (d *commentDao) Update(t *model.Comment) (err error) {
	err = db.Save(t).Error
	return
}

func (d *commentDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Comment{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *commentDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Comment{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *commentDao) Delete(id int64) {
	db.Delete(&model.Comment{}, "id = ?", id)
}
