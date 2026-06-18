package repository

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

type PostRepository interface {
	Get(id int64) *model.Post
	Take(where ...interface{}) *model.Post
	Find(cnd *querybuilder.QueryBuilder) []model.Post
	FindOne(cnd *querybuilder.QueryBuilder) *model.Post
	List(cnd *querybuilder.QueryBuilder) ([]model.Post, *querybuilder.Paging)
	Count(cnd *querybuilder.QueryBuilder) int64
	Create(t *model.Post) error
	Update(t *model.Post) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64) error
	GetRootPosts(limit int) ([]*model.Post, error)
	IncrViewCount(id int64) error
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepo{db: db}
}

type postRepo struct {
	db *gorm.DB
}

// Get 根据 ID 获取帖子，未找到返回 nil
func (r *postRepo) Get(id int64) *model.Post {
	ret := &model.Post{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

// Take 按条件获取单条记录（无排序保证），未找到返回 nil
func (r *postRepo) Take(where ...interface{}) *model.Post {
	ret := &model.Post{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *postRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.Post) {
	cnd.Find(r.db, &list)
	return
}

// FindOne 通过 QueryBuilder 查询单条记录
func (r *postRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.Post {
	ret := &model.Post{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *postRepo) List(cnd *querybuilder.QueryBuilder) (list []model.Post, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.Post{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

// Count 统计数量
func (r *postRepo) Count(cnd *querybuilder.QueryBuilder) int64 {
	return cnd.Count(r.db, &model.Post{})
}

func (r *postRepo) Create(t *model.Post) error {
	return r.db.Create(t).Error
}

func (r *postRepo) Update(t *model.Post) error {
	return r.db.Save(t).Error
}

func (r *postRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(columns).Error
}

func (r *postRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 根据 ID 删除
func (r *postRepo) Delete(id int64) error {
	return r.db.Delete(&model.Post{}, "id = ?", id).Error
}

// GetRootPosts 获取根帖子列表
func (r *postRepo) GetRootPosts(limit int) ([]*model.Post, error) {
	var posts []*model.Post
	err := r.db.Where("parent_id = ?", 0).
		Order("id DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// IncrViewCount 原子递增指定字段
func (r *postRepo) IncrViewCount(id int64) error {
	field := "view_count"
	return r.db.Model(&model.Post{}).
		Where("id = ?", id).
		UpdateColumn(field, gorm.Expr(field+" + ?", 1)).Error
}
