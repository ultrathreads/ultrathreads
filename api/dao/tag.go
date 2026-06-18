package dao

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// TagRepository 标签数据访问契约
type TagRepository interface {
	Get(id int64) *model.Tag
	Take(where ...interface{}) *model.Tag
	Find(cnd *querybuilder.QueryBuilder) []model.Tag
	FindOne(cnd *querybuilder.QueryBuilder) *model.Tag
	List(cnd *querybuilder.QueryBuilder) ([]model.Tag, *querybuilder.Paging)
	Create(t *model.Tag) error
	Update(t *model.Tag) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64) error
	FindByIds(ids []int64) []model.Tag
	FindTagIdsByPostIds(postIds []int64) map[int64][]int64
	GetByName(name string) *model.Tag
	GetOrCreate(name string) (*model.Tag, error)
	GetOrCreates(tags []string) []int64
}

func NewTagDao(db *gorm.DB) TagRepository {
	return &tagRepo{db: db}
}

type tagRepo struct {
	db *gorm.DB
}

func (r *tagRepo) Get(id int64) *model.Tag {
	ret := &model.Tag{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tagRepo) Take(where ...interface{}) *model.Tag {
	ret := &model.Tag{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tagRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.Tag) {
	cnd.Find(r.db, &list)
	return
}

func (r *tagRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.Tag {
	ret := &model.Tag{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		return nil
	}
	return ret
}

func (r *tagRepo) List(cnd *querybuilder.QueryBuilder) (list []model.Tag, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.Tag{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *tagRepo) Create(t *model.Tag) error {
	return r.db.Create(t).Error
}

func (r *tagRepo) Update(t *model.Tag) error {
	return r.db.Save(t).Error
}

func (r *tagRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.Tag{}).Where("id = ?", id).Updates(columns).Error
}

func (r *tagRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *tagRepo) Delete(id int64) error {
	return r.db.Delete(&model.Tag{}, "id = ?", id).Error
}

func (r *tagRepo) FindByIds(ids []int64) []model.Tag {
	if len(ids) == 0 {
		return nil
	}
	var tags []model.Tag
	r.db.Where("id IN ?", ids).Find(&tags)
	return tags
}

func (r *tagRepo) FindTagIdsByPostIds(postIds []int64) map[int64][]int64 {
	if len(postIds) == 0 {
		return nil
	}
	var rows []model.PostTag
	r.db.Select("post_id, tag_id").
		Where("post_id IN ? AND status = ?", postIds, model.StatusOk).
		Find(&rows)

	result := make(map[int64][]int64, len(rows))
	for _, row := range rows {
		result[row.PostId] = append(result[row.PostId], row.TagId)
	}
	return result
}

func (r *tagRepo) GetByName(name string) *model.Tag {
	if len(name) == 0 {
		return nil
	}
	return r.Take("name = ?", name)
}

func (r *tagRepo) GetOrCreate(name string) (*model.Tag, error) {
	if len(name) == 0 {
		return nil, errors.New("标签为空")
	}
	tag := r.GetByName(name)
	if tag != nil {
		return tag, nil
	}
	tag = &model.Tag{
		Name:   name,
		Status: model.StatusOk,
	}
	err := r.Create(tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *tagRepo) GetOrCreates(tags []string) (tagIDs []int64) {
	for _, tagName := range tags {
		tagName = strings.TrimSpace(tagName)
		tag, err := r.GetOrCreate(tagName)
		if err == nil {
			tagIDs = append(tagIDs, tag.ID)
		}
	}
	return
}
