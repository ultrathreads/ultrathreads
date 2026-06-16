package dao

import (
	"errors"
	"gorm.io/gorm"
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

type nodeDao struct {
    db *gorm.DB
}

func NewNodeDao(db *gorm.DB) *nodeDao {
    return &nodeDao{db: db}
}

func (d *nodeDao) Get(id int64) *model.Node {
	ret := &model.Node{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) FindByIds(ids []int64) ([]model.Node) {
	if len(ids) == 0 {
		return nil
	}

	// 可选优化：对 ids 去重，避免生成冗余的 IN 条件

	qb := querybuilder.NewQueryBuilder().In("id", ids)
	
	// 必须捕获并返回 Find 内部的错误
	results := d.Find(qb)
	
	return results
}

func (d *nodeDao) Take(where ...interface{}) *model.Node {
	ret := &model.Node{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Node) {
	cnd.Find(d.db, &list)
	return
}

func (d *nodeDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Node {
	ret := &model.Node{}
	if err := cnd.FindOne(d.db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) List(cnd *querybuilder.QueryBuilder) (list []model.Node, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.Node{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total: count,
	}
	return
}

func (d *nodeDao) Create(t *model.Node) error {
	return d.db.Create(t).Error
}

func (d *nodeDao) Update(t *model.Node) error {
	return d.db.Save(t).Error
}

func (d *nodeDao) Updates(id int64, columns map[string]interface{}) error {
	return d.db.Model(&model.Node{}).Where("id = ?", id).Updates(columns).Error
}

func (d *nodeDao) UpdateColumn(id int64, name string, value interface{}) error {
	return d.db.Model(&model.Node{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// IncrField 原子递增/递减指定字段（Service 层解耦 gorm 的关键方法）
func (d *nodeDao) IncrField(id int64, field string, delta int) error {
	return d.db.Model(&model.Node{}).
		Where("id = ?", id).
		UpdateColumn(field, gorm.Expr(field+" + ?", delta)).Error
}

// Delete 返回 error，不再静默吞掉删除失败
func (d *nodeDao) Delete(id int64) error {
	result := d.db.Delete(&model.Node{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}