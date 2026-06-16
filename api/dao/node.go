package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// NodeRepository 节点数据访问契约
type NodeRepository interface {
	Get(id int64) *model.Node
	Take(where ...interface{}) *model.Node
	Find(cnd *querybuilder.QueryBuilder) []model.Node
	FindOne(cnd *querybuilder.QueryBuilder) *model.Node
	FindByIds(ids []int64) []model.Node
	List(cnd *querybuilder.QueryBuilder) ([]model.Node, *querybuilder.Paging)
	Create(node *model.Node) error
	Update(node *model.Node) error
	Updates(id int64, fields map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	IncrField(id int64, field string, delta int) error
	Delete(id int64) error
	Transaction(fn func(txRepo NodeRepository) error) error
}

type nodeDao struct {
	db *gorm.DB
}

func NewNodeDao(db *gorm.DB) NodeRepository {
	return &nodeDao{db: db}
}

func (d *nodeDao) Get(id int64) *model.Node {
	ret := &model.Node{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) Take(where ...interface{}) *model.Node {
	ret := &model.Node{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) Find(cnd *querybuilder.QueryBuilder) []model.Node {
	var list []model.Node
	cnd.Find(d.db, &list)
	return list
}

func (d *nodeDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Node {
	ret := &model.Node{}
	if err := cnd.FindOne(d.db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) FindByIds(ids []int64) []model.Node {
	if len(ids) == 0 {
		return nil
	}
	qb := querybuilder.NewQueryBuilder().In("id", ids)
	return d.Find(qb)
}

func (d *nodeDao) List(cnd *querybuilder.QueryBuilder) ([]model.Node, *querybuilder.Paging) {
	var list []model.Node
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.Node{})

	paging := &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return list, paging
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

// IncrField 原子递增/递减指定字段
func (d *nodeDao) IncrField(id int64, field string, delta int) error {
	return d.db.Model(&model.Node{}).
		Where("id = ?", id).
		UpdateColumn(field, gorm.Expr(field+" + ?", delta)).Error
}

// Delete 软删除或硬删除（根据 model 是否包含 gorm.DeletedAt 自动判断）
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

// ✅ Transaction 在事务中执行自定义操作
// txRepo 是绑定了当前事务的 Repository 实例，闭包内必须用它
func (d *nodeDao) Transaction(fn func(txRepo NodeRepository) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &nodeDao{db: tx}
		return fn(txRepo)
	})
}