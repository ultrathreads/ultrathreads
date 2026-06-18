package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// PostLikeRepository 点赞数据访问契约
type PostLikeRepository interface {
	Get(id int64) *model.PostLike
	Take(where ...interface{}) *model.PostLike
	Find(cnd *querybuilder.QueryBuilder) []model.PostLike
	FindOne(cnd *querybuilder.QueryBuilder) *model.PostLike
	List(cnd *querybuilder.QueryBuilder) ([]model.PostLike, *querybuilder.Paging)
	Count(cnd *querybuilder.QueryBuilder) int64
	Create(t *model.PostLike) error
	Update(t *model.PostLike) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64) error
}

type postLikeRepo struct {
	db *gorm.DB
}

func NewPostLikeDao(db *gorm.DB) PostLikeRepository {
	return &postLikeRepo{db: db}
}

func (r *postLikeRepo) Get(id int64) *model.PostLike {
	ret := &model.PostLike{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *postLikeRepo) Take(where ...interface{}) *model.PostLike {
	ret := &model.PostLike{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *postLikeRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.PostLike) {
	cnd.Find(r.db, &list)
	return
}

func (r *postLikeRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.PostLike {
	ret := &model.PostLike{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *postLikeRepo) List(cnd *querybuilder.QueryBuilder) (list []model.PostLike, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.PostLike{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *postLikeRepo) Count(cnd *querybuilder.QueryBuilder) int64 {
	return cnd.Count(r.db, &model.PostLike{})
}

func (r *postLikeRepo) Create(t *model.PostLike) error {
	return r.db.Create(t).Error
}

func (r *postLikeRepo) Update(t *model.PostLike) error {
	return r.db.Save(t).Error
}

func (r *postLikeRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.PostLike{}).Where("id = ?", id).Updates(columns).Error
}

func (r *postLikeRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.PostLike{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *postLikeRepo) Delete(id int64) error {
	return r.db.Delete(&model.PostLike{}, "id = ?", id).Error
}