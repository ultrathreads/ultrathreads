package repository

import (
	"time"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// PostTagRepository 帖子标签关联数据访问契约
type PostTagRepository interface {
	Get(id int64) *model.PostTag
	Take(where ...interface{}) *model.PostTag
	Find(cnd *querybuilder.QueryBuilder) []model.PostTag
	FindOne(cnd *querybuilder.QueryBuilder) *model.PostTag
	List(cnd *querybuilder.QueryBuilder) ([]model.PostTag, *querybuilder.Paging)
	Create(t *model.PostTag) error
	Update(t *model.PostTag) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	AddPostTags(postId int64, tagIds []int64)
	DeletePostTags(postId int64)
}

type postTagRepo struct {
	db *gorm.DB
}

func NewPostTagRepository(db *gorm.DB) PostTagRepository {
	return &postTagRepo{db: db}
}

func (r *postTagRepo) Get(id int64) *model.PostTag {
	ret := &model.PostTag{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *postTagRepo) Take(where ...interface{}) *model.PostTag {
	ret := &model.PostTag{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *postTagRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.PostTag) {
	cnd.Find(r.db, &list)
	return
}

func (r *postTagRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.PostTag {
	ret := &model.PostTag{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		return nil
	}
	return ret
}

func (r *postTagRepo) List(cnd *querybuilder.QueryBuilder) (list []model.PostTag, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.PostTag{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *postTagRepo) Create(t *model.PostTag) error {
	return r.db.Create(t).Error
}

func (r *postTagRepo) Update(t *model.PostTag) error {
	return r.db.Save(t).Error
}

func (r *postTagRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.PostTag{}).Where("id = ?", id).Updates(columns).Error
}

func (r *postTagRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.PostTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *postTagRepo) Delete(id int64) {
	r.db.Delete(&model.PostTag{}, "id = ?", id)
}

func (r *postTagRepo) AddPostTags(postId int64, tagIds []int64) {
	if postId <= 0 || len(tagIds) == 0 {
		return
	}
	for _, tagId := range tagIds {
		_ = r.Create(&model.PostTag{
			PostId:    postId,
			TagId:     tagId,
			CreatedAt: time.Now(),
		})
	}
}

func (r *postTagRepo) DeletePostTags(postId int64) {
	if postId <= 0 {
		return
	}
	r.db.Where("post_id = ?", postId).Delete(&model.PostTag{})
}
