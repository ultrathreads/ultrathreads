package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var PostLikeDao = newPostLikeDao()

func newPostLikeDao() *postLikeDao {
	return &postLikeDao{}
}

type postLikeDao struct{}

// Get 根据 ID 获取点赞记录，未找到返回 nil
func (d *postLikeDao) Get(id int64) *model.PostLike {
	ret := &model.PostLike{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

// Take 按条件获取单条记录，未找到返回 nil
func (d *postLikeDao) Take(where ...interface{}) *model.PostLike {
	ret := &model.PostLike{}
	if err := db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *postLikeDao) Find(cnd *querybuilder.QueryBuilder) (list []model.PostLike) {
	cnd.Find(db, &list)
	return
}

// FindOne 通过 QueryBuilder 查询单条记录
func (d *postLikeDao) FindOne(cnd *querybuilder.QueryBuilder) *model.PostLike {
	ret := &model.PostLike{}
	// ✅ 修复：传入 ret 而非 &ret，避免 **model.PostLike 导致 v2 扫描失败
	if err := cnd.FindOne(db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
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
		Total: count, // ✅ int64，与升级后的 Paging.Total 类型一致
	}
	return
}

// Count 统计数量
func (d *postLikeDao) Count(cnd *querybuilder.QueryBuilder) int64 { // ✅ 新增方法，返回 int64
	return cnd.Count(db, &model.PostLike{})
}

func (d *postLikeDao) Create(t *model.PostLike) error {
	return db.Create(t).Error
}

func (d *postLikeDao) Update(t *model.PostLike) error {
	return db.Save(t).Error
}

func (d *postLikeDao) Updates(id int64, columns map[string]interface{}) error {
	return db.Model(&model.PostLike{}).Where("id = ?", id).Updates(columns).Error
}

func (d *postLikeDao) UpdateColumn(id int64, name string, value interface{}) error {
	return db.Model(&model.PostLike{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 根据 ID 删除
func (d *postLikeDao) Delete(id int64) error { // ✅ 补充 error 返回值
	return db.Delete(&model.PostLike{}, "id = ?", id).Error
}