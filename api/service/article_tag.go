package service

import (
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util/querybuilder"
)

// ArticleTagService 文章标签关联业务契约
type ArticleTagService interface {
	Get(id int64) *model.ArticleTag
	Take(where ...interface{}) *model.ArticleTag
	Find(cnd *querybuilder.QueryBuilder) []model.ArticleTag
	List(cnd *querybuilder.QueryBuilder) ([]model.ArticleTag, *querybuilder.Paging)
	Create(t *model.ArticleTag) error
	Update(t *model.ArticleTag) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	DeleteByArticleId(postId int64)
}

func NewArticleTagService(repo repository.ArticleTagRepository) ArticleTagService {
	return &articleTagService{repo: repo}
}

type articleTagService struct {
	repo repository.ArticleTagRepository
}

func (s *articleTagService) Get(id int64) *model.ArticleTag {
	return s.repo.Get(id)
}

func (s *articleTagService) Take(where ...interface{}) *model.ArticleTag {
	return s.repo.Take(where...)
}

func (s *articleTagService) Find(cnd *querybuilder.QueryBuilder) []model.ArticleTag {
	return s.repo.Find(cnd)
}

func (s *articleTagService) List(cnd *querybuilder.QueryBuilder) ([]model.ArticleTag, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *articleTagService) Create(t *model.ArticleTag) error {
	return s.repo.Create(t)
}

func (s *articleTagService) Update(t *model.ArticleTag) error {
	return s.repo.Update(t)
}

func (s *articleTagService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *articleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *articleTagService) DeleteByArticleId(postId int64) {
	s.repo.Updates(postId, map[string]interface{}{"status": model.StatusDeleted})
}
