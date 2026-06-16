package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

type PostRepository interface {
    Get(id int64) *model.Post
}

func NewPostDao(db *gorm.DB) *postDao {
	return &postDao{db: db}
}

type postDao struct {
	db *gorm.DB
}

// Get 根据 ID 获取帖子，未找到返回 nil
func (d *postDao) Get(id int64) *model.Post {
	ret := &model.Post{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

// Take 按条件获取单条记录（无排序保证），未找到返回 nil
func (d *postDao) Take(where ...interface{}) *model.Post {
	ret := &model.Post{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *postDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Post) {
	cnd.Find(d.db, &list)
	return
}

// FindOne 通过 QueryBuilder 查询单条记录
func (d *postDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Post {
	ret := &model.Post{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *postDao) List(cnd *querybuilder.QueryBuilder) (list []model.Post, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.Post{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

// Count 统计数量
func (d *postDao) Count(cnd *querybuilder.QueryBuilder) int64 {
	return cnd.Count(d.db, &model.Post{})
}

func (d *postDao) Create(t *model.Post) error {
	return d.db.Create(t).Error
}

func (d *postDao) Update(t *model.Post) error {
	return d.db.Save(t).Error
}

func (d *postDao) Updates(id int64, columns map[string]interface{}) error {
	return d.db.Model(&model.Post{}).Where("id = ?", id).Updates(columns).Error
}

func (d *postDao) UpdateColumn(id int64, name string, value interface{}) error {
	return d.db.Model(&model.Post{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 根据 ID 删除
func (d *postDao) Delete(id int64) error {
	return d.db.Delete(&model.Post{}, "id = ?", id).Error
}

// GetRootPosts 获取根帖子列表
func (d *postDao) GetRootPosts(limit int) ([]*model.Post, error) {
	var posts []*model.Post
	err := d.db.Where("parent_id = ?", 0).
		Order("id DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// IncrViewCount 原子递增指定字段
func (d *postDao) IncrViewCount(id int64) error {
	field := "view_count"
	return d.db.Model(&model.Post{}).
		Where("id = ?", id).
		UpdateColumn(field, gorm.Expr(field+" + ?", 1)).Error
}