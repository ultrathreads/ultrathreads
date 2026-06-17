package service

import (
	"strings"
	"errors"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/dto"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/hashid"
)

type ScanTagCallback func(tags []model.Tag) bool

var TagService = newTagService()

func newTagService() *tagService {
	return &tagService{}
}

type tagService struct {
}

func (s *tagService) Get(id int64) *model.Tag {
	return dao.TagDao.Get(id)
}

func (s *tagService) GetBySlug(slug string) *model.Tag {
	id := hashid.Slug2Id[model.Tag](slug)
	return cache.TagCache.Get(id)
}

func (s *tagService) Take(where ...interface{}) *model.Tag {
	return dao.TagDao.Take(where...)
}

func (s *tagService) Find(cnd *querybuilder.QueryBuilder) []model.Tag {
	return dao.TagDao.Find(cnd)
}

func (s *tagService) FindOne(cnd *querybuilder.QueryBuilder) *model.Tag {
	return dao.TagDao.FindOne(cnd)
}

func (s *tagService) List(cnd *querybuilder.QueryBuilder) (list []model.Tag, paging *querybuilder.Paging) {
	return dao.TagDao.List(cnd)
}

func (s *tagService) Create(req dto.TagCreateForm) (*model.Tag, error) {
	tag := &model.Tag{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	if err := dao.TagDao.Create(tag); err != nil {
		return nil, errors.New("创建标签失败")
	}
	return tag, nil
}

func (s *tagService) Update(int int64, req dto.TagUpdateForm) error {
	err := dao.TagDao.Updates(req.ID, map[string]interface{}{
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
	if err := dao.TagDao.Delete(id); err != nil {
		return errors.New("删除标签失败")
	}
	return nil
}

// 自动完成
func (s *tagService) AutoComplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.TagDao.Find(querybuilder.NewQueryBuilder().Where("status = ? and name like ?",
		model.StatusOk, "%"+input+"%").Limit(6))
}

func (s *tagService) GetOrCreate(name string) (*model.Tag, error) {
	return dao.TagDao.GetOrCreate(name)
}

func (s *tagService) GetByName(name string) *model.Tag {
	return dao.TagDao.GetByName(name)
}

func (s *tagService) GetTags() []model.TagResponse {
	list := dao.TagDao.Find(querybuilder.NewQueryBuilder().Where("status = ?", model.StatusOk))

	var tags []model.TagResponse
	for _, tag := range list {
		tags = append(tags, model.TagResponse{TagName: tag.Name})
	}
	return tags
}

func (s *tagService) FindByIds(tagIds []int64) []model.Tag {
	return dao.TagDao.FindByIds(tagIds)
}

// 扫描
func (s *tagService) Scan(cb ScanTagCallback) {
	var cursor int64
	for {
		list := dao.TagDao.Find(querybuilder.NewQueryBuilder().Where("id > ?", cursor).Asc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		if !cb(list) {
			break
		}
	}
}
