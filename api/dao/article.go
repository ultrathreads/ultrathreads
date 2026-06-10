package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var ArticleDao = newArticleDao()

func newArticleDao() *articleDao {
	return &articleDao{}
}

type articleDao struct{}

// Get 根据 ID 获取文章，未找到时返回 nil（兼容原逻辑）
func (d *articleDao) Get(id int64) *model.Article {
	ret := &model.Article{}
	// ⚠️ GORM v2 中 First 未找到记录会返回 gorm.ErrRecordNotFound
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		// 其他数据库错误也可选择返回 nil 或记录日志
		return nil
	}
	return ret
}

func (d *articleDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Article) {
	// ⚠️ 确保 querybuilder 内部已适配 GORM v2
	cnd.Find(db, &list)
	return
}

func (d *articleDao) List(cnd *querybuilder.QueryBuilder) (list []model.Article, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Article{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *articleDao) Create(t *model.Article) error {
	return db.Create(t).Error
}

// Update 全量更新（包含零值）
func (d *articleDao) Update(t *model.Article) error {
	return db.Save(t).Error
}

// Updates 按 map 更新指定字段（支持零值）
func (d *articleDao) Updates(id int64, columns map[string]interface{}) error {
	return db.Model(&model.Article{}).Where("id = ?", id).Updates(columns).Error
}

// UpdateColumn 更新单个列（跳过 Hook 和更新时间）
func (d *articleDao) UpdateColumn(id int64, name string, value interface{}) error {
	return db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 根据 ID 删除
func (d *articleDao) Delete(id int64) error {
	// ✅ v2 建议返回 error，便于上层感知删除是否成功
	return db.Delete(&model.Article{}, "id = ?", id).Error
}