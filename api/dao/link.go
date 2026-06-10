package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var LinkDao = newLinkDao()

func newLinkDao() *linkDao {
	return &linkDao{}
}

type linkDao struct{}

// Get 根据 ID 获取链接，未找到返回 nil
func (d *linkDao) Get(id int64) *model.Link {
	ret := &model.Link{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

// Take 按条件获取单条记录，未找到返回 nil
func (d *linkDao) Take(where ...interface{}) *model.Link {
	ret := &model.Link{}
	if err := db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *linkDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Link) {
	cnd.Find(db, &list)
	return
}

// FindOne 通过 QueryBuilder 查询单条记录
func (d *linkDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Link {
	ret := &model.Link{}
	// ✅ 修复：传入 ret 而非 &ret，避免 **model.Link 导致 v2 扫描失败
	if err := cnd.FindOne(db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
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
		Total: count, // ✅ int64，与升级后的 Paging.Total 类型一致
	}
	return
}

// Count 统计数量
func (d *linkDao) Count(cnd *querybuilder.QueryBuilder) int64 { // ✅ 改为 int64
	return cnd.Count(db, &model.Link{})
}

func (d *linkDao) Create(t *model.Link) error {
	return db.Create(t).Error
}

func (d *linkDao) Update(t *model.Link) error {
	return db.Save(t).Error
}

func (d *linkDao) Updates(id int64, columns map[string]interface{}) error {
	return db.Model(&model.Link{}).Where("id = ?", id).Updates(columns).Error
}

func (d *linkDao) UpdateColumn(id int64, name string, value interface{}) error {
	return db.Model(&model.Link{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 根据 ID 删除
func (d *linkDao) Delete(id int64) error { // ✅ 补充 error 返回值
	return db.Delete(&model.Link{}, "id = ?", id).Error
}