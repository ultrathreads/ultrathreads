package service

import (
	"errors"

	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
)

// FavoriteService 收藏业务契约
type FavoriteService interface {
	Get(id int64) *model.Favorite
	Take(where ...interface{}) *model.Favorite
	Find(cnd *querybuilder.QueryBuilder) []model.Favorite
	FindOne(cnd *querybuilder.QueryBuilder) *model.Favorite
	List(cnd *querybuilder.QueryBuilder) ([]model.Favorite, *querybuilder.Paging)
	Create(t *model.Favorite) error
	Update(t *model.Favorite) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	GetBy(userId int64, entityType string, entityId int64) *model.Favorite
	AddArticleFavorite(userId, articleId int64) error
	AddPostFavorite(userId int64, postSlug string) error
}

func NewFavoriteService(repo repository.FavoriteRepository, articleRepo repository.ArticleRepository, postRepo repository.PostRepository) FavoriteService {
	return &favoriteService{repo: repo, articleRepo: articleRepo, postRepo: postRepo}
}

type favoriteService struct {
	repo        repository.FavoriteRepository
	articleRepo repository.ArticleRepository
	postRepo    repository.PostRepository
}

func (s *favoriteService) Get(id int64) *model.Favorite {
	return s.repo.Get(id)
}

func (s *favoriteService) Take(where ...interface{}) *model.Favorite {
	return s.repo.Take(where...)
}

func (s *favoriteService) Find(cnd *querybuilder.QueryBuilder) []model.Favorite {
	return s.repo.Find(cnd)
}

func (s *favoriteService) FindOne(cnd *querybuilder.QueryBuilder) *model.Favorite {
	return s.repo.FindOne(cnd)
}

func (s *favoriteService) List(cnd *querybuilder.QueryBuilder) ([]model.Favorite, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *favoriteService) Create(t *model.Favorite) error {
	return s.repo.Create(t)
}

func (s *favoriteService) Update(t *model.Favorite) error {
	return s.repo.Update(t)
}

func (s *favoriteService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *favoriteService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *favoriteService) Delete(id int64) {
	s.repo.Delete(id)
}

func (s *favoriteService) GetBy(userId int64, entityType string, entityId int64) *model.Favorite {
	return s.repo.Take("user_id = ? and entity_type = ? and entity_id = ?",
		userId, entityType, entityId)
}

func (s *favoriteService) AddArticleFavorite(userId, articleId int64) error {
	article := s.articleRepo.Get(articleId)
	if article == nil || article.Status != model.StatusOk {
		return errors.New("收藏的文章不存在")
	}
	return s.addFavorite(userId, model.EntityTypeArticle, articleId)
}

func (s *favoriteService) AddPostFavorite(userId int64, postSlug string) error {
	postId := hashid.Slug2Id[model.Post](postSlug)
	post := s.postRepo.Get(postId)
	if post == nil || post.Status != model.StatusOk {
		return errors.New("收藏的话题不存在")
	}
	return s.addFavorite(userId, model.EntityTypePost, postId)
}

func (s *favoriteService) addFavorite(userId int64, entityType string, entityId int64) error {
	temp := s.GetBy(userId, entityType, entityId)
	if temp != nil {
		return nil
	}
	return s.repo.Create(&model.Favorite{
		UserId:     userId,
		EntityType: entityType,
		EntityId:   entityId,
		CreateTime: util.NowTimestamp(),
	})
}
