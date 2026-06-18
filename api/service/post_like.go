package service

import (
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// PostLikeServicer 点赞业务契约
type PostLikeServicer interface {
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
	Like(userID, postID int64) error
}

func NewPostLikeService(repo repository.PostLikeRepository) PostLikeServicer {
	return &postLikeService{repo: repo}
}

type postLikeService struct {
	repo repository.PostLikeRepository
}

func (s *postLikeService) Get(id int64) *model.PostLike {
	return s.repo.Get(id)
}

func (s *postLikeService) Take(where ...interface{}) *model.PostLike {
	return s.repo.Take(where...)
}

func (s *postLikeService) Find(cnd *querybuilder.QueryBuilder) []model.PostLike {
	return s.repo.Find(cnd)
}

func (s *postLikeService) FindOne(cnd *querybuilder.QueryBuilder) *model.PostLike {
	return s.repo.FindOne(cnd)
}

func (s *postLikeService) List(cnd *querybuilder.QueryBuilder) ([]model.PostLike, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *postLikeService) Count(cnd *querybuilder.QueryBuilder) int64 {
	return s.repo.Count(cnd)
}

func (s *postLikeService) Create(t *model.PostLike) error {
	return s.repo.Create(t)
}

func (s *postLikeService) Update(t *model.PostLike) error {
	return s.repo.Update(t)
}

func (s *postLikeService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *postLikeService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *postLikeService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Like 点赞或取消点赞
func (s *postLikeService) Like(userID, postID int64) error {
	tmp := s.repo.Take("user_id = ? and post_id = ?", userID, postID)
	if tmp != nil {
		return s.repo.Delete(tmp.ID)
	}
	return s.repo.Create(&model.PostLike{
		UserId:     userID,
		PostId:     postID,
		CreateTime: util.NowTimestamp(),
	})
}
