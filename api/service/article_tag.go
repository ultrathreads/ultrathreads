package service

import (
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var ArticleTagService = newArticleTagService()

func newArticleTagService() *articleTagService {
	return &articleTagService{}
}

type articleTagService struct {
}

func (s *articleTagService) Get(id int64) *model.ArticleTag {
	return dao.ArticleTagDao.Get(id)
}

func (s *articleTagService) Take(where ...interface{}) *model.ArticleTag {
	return dao.ArticleTagDao.Take(where...)
}

func (s *articleTagService) Find(cnd *querybuilder.QueryBuilder) []model.ArticleTag {
	return dao.ArticleTagDao.Find(cnd)
}

func (s *articleTagService) List(cnd *querybuilder.QueryBuilder) (list []model.ArticleTag, paging *querybuilder.Paging) {
	return dao.ArticleTagDao.List(cnd)
}

func (s *articleTagService) Create(t *model.ArticleTag) error {
	return dao.ArticleTagDao.Create(t)
}

func (s *articleTagService) Update(t *model.ArticleTag) error {
	return dao.ArticleTagDao.Update(t)
}

func (s *articleTagService) Updates(id int64, columns map[string]interface{}) error {
	return dao.ArticleTagDao.Updates(id, columns)
}

func (s *articleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.ArticleTagDao.UpdateColumn(id, name, value)
}

func (s *articleTagService) DeleteByArticleId(postId int64) {
	dao.DB().Model(model.ArticleTag{}).Where("article_id = ?", postId).UpdateColumn("status", model.StatusDeleted)
}
