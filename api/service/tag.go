package service

import (
	"errors"
	"strings"

	"ultrathreads/cache"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
)

type ScanTagCallback func(tags []model.Tag) bool

// TagServicer 标签业务契约
type TagServicer interface {
	Get(id int64) *model.Tag
	GetBySlug(slug string) *model.Tag
	Take(where ...interface{}) *model.Tag
	Find(cnd *querybuilder.QueryBuilder) []model.Tag
	FindOne(cnd *querybuilder.QueryBuilder) *model.Tag
	List(cnd *querybuilder.QueryBuilder) ([]model.Tag, *querybuilder.Paging)
	Create(req dto.TagCreateForm) (*model.Tag, error)
	Update(int int64, req dto.TagUpdateForm) error
	Delete(id int64) error
	AutoComplete(input string) []model.Tag
	GetOrCreate(name string) (*model.Tag, error)
	GetByName(name string) *model.Tag
	GetTags() []model.TagResponse
	FindByIds(tagIds []int64) []model.Tag
	Scan(cb ScanTagCallback)
}

func NewTagService(repo repository.TagRepository, tagCache cache.TagCacheInterface) TagServicer {
	return &tagService{repo: repo, tagCache: tagCache}
}

type tagService struct {
	repo     repository.TagRepository
	tagCache cache.TagCacheInterface
}

func (s *tagService) Get(id int64) *model.Tag {
	return s.repo.Get(id)
}

func (s *tagService) GetBySlug(slug string) *model.Tag {
	id := hashid.Slug2Id[model.Tag](slug)
	return s.tagCache.Get(id)
}

func (s *tagService) Take(where ...interface{}) *model.Tag {
	return s.repo.Take(where...)
}

func (s *tagService) Find(cnd *querybuilder.QueryBuilder) []model.Tag {
	return s.repo.Find(cnd)
}

func (s *tagService) FindOne(cnd *querybuilder.QueryBuilder) *model.Tag {
	return s.repo.FindOne(cnd)
}

func (s *tagService) List(cnd *querybuilder.QueryBuilder) ([]model.Tag, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *tagService) Create(req dto.TagCreateForm) (*model.Tag, error) {
	tag := &model.Tag{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	if err := s.repo.Create(tag); err != nil {
		return nil, errors.New("创建标签失败")
	}
	return tag, nil
}

func (s *tagService) Update(int int64, req dto.TagUpdateForm) error {
	err := s.repo.Updates(req.ID, map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"status":      req.Status,
	})
	if err != nil {
		return errors.New("更新标签失败")
	}
	return nil
}

func (s *tagService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return errors.New("删除标签失败")
	}
	return nil
}

func (s *tagService) AutoComplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return s.repo.Find(querybuilder.NewQueryBuilder().Where("status = ? and name like ?",
		model.StatusOk, "%"+input+"%").Limit(6))
}

func (s *tagService) GetOrCreate(name string) (*model.Tag, error) {
	return s.repo.GetOrCreate(name)
}

func (s *tagService) GetByName(name string) *model.Tag {
	return s.repo.GetByName(name)
}

func (s *tagService) GetTags() []model.TagResponse {
	list := s.repo.Find(querybuilder.NewQueryBuilder().Where("status = ?", model.StatusOk))

	var tags []model.TagResponse
	for _, tag := range list {
		tags = append(tags, model.TagResponse{TagName: tag.Name})
	}
	return tags
}

func (s *tagService) FindByIds(tagIds []int64) []model.Tag {
	return s.repo.FindByIds(tagIds)
}

func (s *tagService) Scan(cb ScanTagCallback) {
	var cursor int64
	for {
		list := s.repo.Find(querybuilder.NewQueryBuilder().Where("id > ?", cursor).Asc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		if !cb(list) {
			break
		}
	}
}
