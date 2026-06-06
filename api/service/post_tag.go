package service

import (
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var PostTagService = newPostTagService()

func newPostTagService() *postTagService {
	return &postTagService{}
}

type postTagService struct {
}

func (s *postTagService) Get(id int64) *model.PostTag {
	return dao.PostTagDao.Get(id)
}

func (s *postTagService) Take(where ...interface{}) *model.PostTag {
	return dao.PostTagDao.Take(where...)
}

func (s *postTagService) Find(cnd *querybuilder.QueryBuilder) []model.PostTag {
	return dao.PostTagDao.Find(cnd)
}

func (s *postTagService) FindOne(cnd *querybuilder.QueryBuilder) *model.PostTag {
	return dao.PostTagDao.FindOne(cnd)
}

func (s *postTagService) List(cnd *querybuilder.QueryBuilder) (list []model.PostTag, paging *querybuilder.Paging) {
	return dao.PostTagDao.List(cnd)
}

func (s *postTagService) Create(t *model.PostTag) error {
	return dao.PostTagDao.Create(t)
}

func (s *postTagService) Update(t *model.PostTag) error {
	return dao.PostTagDao.Update(t)
}

func (s *postTagService) Updates(id int64, columns map[string]interface{}) error {
	return dao.PostTagDao.Updates(id, columns)
}

func (s *postTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.PostTagDao.UpdateColumn(id, name, value)
}

func (s *postTagService) DeleteByPostId(postId int64) {
	dao.DB().Model(model.PostTag{}).Where("post_id = ?", postId).UpdateColumn("status", model.StatusDeleted)
}

func (s *postTagService) UndeleteByPostId(postId int64) {
	dao.DB().Model(model.PostTag{}).Where("post_id = ?", postId).UpdateColumn("status", model.StatusOk)
}
