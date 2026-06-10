package dao

import (
	"errors"
	"gorm.io/gorm"
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var NodeDao = newNodeDao()

func newNodeDao() *nodeDao {
	return &nodeDao{}
}

type nodeDao struct{}

func (d *nodeDao) Get(id int64) *model.Node {
	ret := &model.Node{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) Take(where ...interface{}) *model.Node {
	ret := &model.Node{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Node) {
	cnd.Find(db, &list)
	return
}

func (d *nodeDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Node {
	ret := &model.Node{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) List(cnd *querybuilder.QueryBuilder) (list []model.Node, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Node{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *nodeDao) Create(t *model.Node) error {
	return db.Create(t).Error
}

func (d *nodeDao) Update(t *model.Node) error {
	return db.Save(t).Error
}

func (d *nodeDao) Updates(id int64, columns map[string]interface{}) error {
	return db.Model(&model.Node{}).Where("id = ?", id).Updates(columns).Error
}

func (d *nodeDao) UpdateColumn(id int64, name string, value interface{}) error {
	return db.Model(&model.Node{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// IncrField 原子递增/递减指定字段（Service 层解耦 gorm 的关键方法）
func (d *nodeDao) IncrField(id int64, field string, delta int) error {
	return db.Model(&model.Node{}).
		Where("id = ?", id).
		UpdateColumn(field, gorm.Expr(field+" + ?", delta)).Error
}

// Delete 返回 error，不再静默吞掉删除失败
func (d *nodeDao) Delete(id int64) error {
	result := db.Delete(&model.Node{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}