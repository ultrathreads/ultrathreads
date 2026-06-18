package service

import (
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util/querybuilder"
)

// PostTagService 帖子标签关联业务契约
type PostTagService interface {
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
	DeleteByPostId(postId int64)
	UndeleteByPostId(postId int64)
}

func NewPostTagService(repo repository.PostTagRepository) PostTagService {
	return &postTagService{repo: repo}
}

type postTagService struct {
	repo repository.PostTagRepository
}

func (s *postTagService) Get(id int64) *model.PostTag {
	return s.repo.Get(id)
}

func (s *postTagService) Take(where ...interface{}) *model.PostTag {
	return s.repo.Take(where...)
}

func (s *postTagService) Find(cnd *querybuilder.QueryBuilder) []model.PostTag {
	return s.repo.Find(cnd)
}

func (s *postTagService) FindOne(cnd *querybuilder.QueryBuilder) *model.PostTag {
	return s.repo.FindOne(cnd)
}

func (s *postTagService) List(cnd *querybuilder.QueryBuilder) ([]model.PostTag, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *postTagService) Create(t *model.PostTag) error {
	return s.repo.Create(t)
}

func (s *postTagService) Update(t *model.PostTag) error {
	return s.repo.Update(t)
}

func (s *postTagService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *postTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *postTagService) Delete(id int64) {
	s.repo.Delete(id)
}

func (s *postTagService) AddPostTags(postId int64, tagIds []int64) {
	s.repo.AddPostTags(postId, tagIds)
}

func (s *postTagService) DeletePostTags(postId int64) {
	s.repo.DeletePostTags(postId)
}

func (s *postTagService) DeleteByPostId(postId int64) {
	s.repo.Updates(postId, map[string]interface{}{"status": model.StatusDeleted})
}

func (s *postTagService) UndeleteByPostId(postId int64) {
	s.repo.Updates(postId, map[string]interface{}{"status": model.StatusOk})
}
