package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var PostDao = newPostDao()

func newPostDao() *postDao {
	return &postDao{}
}

type postDao struct{}

// Get 根据 ID 获取帖子，未找到返回 nil
func (d *postDao) Get(id int64) *model.Post {
	ret := &model.Post{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		// 其他数据库错误可选择记录日志
		return nil
	}
	return ret
}

// Take 按条件获取单条记录（无排序保证），未找到返回 nil
func (d *postDao) Take(where ...interface{}) *model.Post {
	ret := &model.Post{}
	if err := db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *postDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Post) {
	cnd.Find(db, &list)
	return
}

// FindOne 通过 QueryBuilder 查询单条记录
func (d *postDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Post {
	ret := &model.Post{}
	if err := cnd.FindOne(db, ret); err != nil { // ✅ 传入 ret 而非 &ret，避免 **model.Post
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *postDao) List(cnd *querybuilder.QueryBuilder) (list []model.Post, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Post{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total: count,
	}
	return
}

// Count 统计数量
func (d *postDao) Count(cnd *querybuilder.QueryBuilder) int64 { // ✅ 改为 int64
	return cnd.Count(db, &model.Post{})
}

func (d *postDao) Create(t *model.Post) error {
	return db.Create(t).Error
}

func (d *postDao) Update(t *model.Post) error {
	return db.Save(t).Error
}

func (d *postDao) Updates(id int64, columns map[string]interface{}) error {
	return db.Model(&model.Post{}).Where("id = ?", id).Updates(columns).Error
}

func (d *postDao) UpdateColumn(id int64, name string, value interface{}) error {
	return db.Model(&model.Post{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 根据 ID 删除
func (d *postDao) Delete(id int64) error { // ✅ 补充 error 返回值
	return db.Delete(&model.Post{}, "id = ?", id).Error
}

// GetRootPosts 获取根帖子列表
func (d *postDao) GetRootPosts(limit int) ([]*model.Post, error) {
	var posts []*model.Post
	err := db.Where("parent_id = ?", 0).
		Order("id DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// IncrViewCount 原子递增/递减指定字段（Service 层解耦 gorm 的关键方法）
func (d *postDao) IncrViewCount(id int64) error {
	field := "view_count"
	return db.Model(&model.Post{}).
		Where("id = ?", id).
		UpdateColumn(field, gorm.Expr(field+" + ?", 1)).Error
}
