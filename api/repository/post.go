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
	// 以下方法封装事务/多表操作，避免 service 层依赖 *gorm.DB
	FindByIds(ids []int64) []model.Post
	FindPostIdsByTagId(tagId int64) []int64
	CreateRootPost(post *model.Post) error
	CreateReply(post *model.Post, parentID int64) error
	OnComment(postId, lastCommentUserId, lastCommentTime int64) error
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

// FindByIds 根据 ID 列表批量查询
func (r *postRepo) FindByIds(ids []int64) []model.Post {
	if len(ids) == 0 {
		return nil
	}
	var posts []model.Post
	r.db.Where("id IN (?)", ids).Find(&posts)
	return posts
}

// FindPostIdsByTagId 查询指定标签下的所有帖子ID
func (r *postRepo) FindPostIdsByTagId(tagId int64) []int64 {
	var postIds []int64
	r.db.Model(&model.PostTag{}).
		Where("tag_id = ? AND status = ?", tagId, model.StatusOk).
		Pluck("post_id", &postIds)
	return postIds
}

// CreateRootPost 创建根帖（事务：创建帖子 + 更新 thread_id = 自身 ID）
func (r *postRepo) CreateRootPost(post *model.Post) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		if err := tx.Model(post).UpdateColumn("thread_id", post.ID).Error; err != nil {
			return err
		}
		post.ThreadId = post.ID
		return nil
	})
}

// CreateReply 创建回复（事务：校验父帖 + 创建回复 + 更新根帖最后回复时间）
func (r *postRepo) CreateReply(post *model.Post, parentID int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var parentPost model.Post
		if err := tx.Where("id = ? AND status = ?", parentID, model.StatusOk).First(&parentPost).Error; err != nil {
			return err
		}

		threadId := parentPost.ThreadId
		if threadId == 0 {
			threadId = parentPost.ID
		}
		post.ParentId = parentID
		post.ThreadId = threadId
		post.NodeId = parentPost.NodeId

		if err := tx.Create(post).Error; err != nil {
			return err
		}

		return tx.Model(&model.Post{}).
			Where("id = ?", threadId).
			UpdateColumn("last_replied_at", post.CreatedAt).Error
	})
}

// OnComment 评论时更新帖子统计（事务：更新帖子 + 更新帖子标签）
func (r *postRepo) OnComment(postId, lastCommentUserId, lastCommentTime int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where("id = ?", postId).Updates(map[string]interface{}{
			"comment_count":      gorm.Expr("comment_count + ?", 1),
			"last_reply_user_id": lastCommentUserId,
			"last_replied_at":    lastCommentTime,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&model.PostTag{}).Where("post_id = ?", postId).Updates(map[string]interface{}{
			"last_replied_at": lastCommentTime,
		}).Error
	})
}
