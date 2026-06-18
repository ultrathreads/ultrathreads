package repository

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

type nodeRepo struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) NodeRepository {
	return &nodeRepo{db: db}
}

func (r *nodeRepo) Get(id int64) *model.Node {
	ret := &model.Node{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *nodeRepo) Take(where ...interface{}) *model.Node {
	ret := &model.Node{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *nodeRepo) Find(cnd *querybuilder.QueryBuilder) []model.Node {
	var list []model.Node
	cnd.Find(r.db, &list)
	return list
}

func (r *nodeRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.Node {
	ret := &model.Node{}
	if err := cnd.FindOne(r.db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *nodeRepo) FindByIds(ids []int64) []model.Node {
	if len(ids) == 0 {
		return nil
	}
	qb := querybuilder.NewQueryBuilder().In("id", ids)
	return r.Find(qb)
}

func (r *nodeRepo) List(cnd *querybuilder.QueryBuilder) ([]model.Node, *querybuilder.Paging) {
	var list []model.Node
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.Node{})

	paging := &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return list, paging
}

func (r *nodeRepo) Create(t *model.Node) error {
	return r.db.Create(t).Error
}

func (r *nodeRepo) Update(t *model.Node) error {
	return r.db.Save(t).Error
}

func (r *nodeRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.Node{}).Where("id = ?", id).Updates(columns).Error
}

func (r *nodeRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.Node{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// IncrField 原子递增/递减指定字段
func (r *nodeRepo) IncrField(id int64, field string, delta int) error {
	return r.db.Model(&model.Node{}).
		Where("id = ?", id).
		UpdateColumn(field, gorm.Expr(field+" + ?", delta)).Error
}

// Delete 软删除或硬删除（根据 model 是否包含 gorm.DeletedAt 自动判断）
func (r *nodeRepo) Delete(id int64) error {
	result := r.db.Delete(&model.Node{}, "id = ?", id)
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
func (r *nodeRepo) Transaction(fn func(txRepo NodeRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &nodeRepo{db: tx}
		return fn(txRepo)
	})
}
