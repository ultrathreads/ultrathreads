package service

import (
	"errors"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
)

// UserWatchService 用户关注业务契约
type UserWatchService interface {
	Get(id int64) *model.UserWatch
	Take(where ...interface{}) *model.UserWatch
	Find(cnd *querybuilder.QueryBuilder) []model.UserWatch
	FindOne(cnd *querybuilder.QueryBuilder) *model.UserWatch
	List(cnd *querybuilder.QueryBuilder) ([]model.UserWatch, *querybuilder.Paging)
	Create(t *model.UserWatch) error
	Update(t *model.UserWatch) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	Watch(userSlug string, userID int64) error
	GetBy(userID, watchedUserID int64) *model.UserWatch
}

func NewUserWatchService(repo repository.UserWatchRepository) UserWatchService {
	return &userWatchService{repo: repo}
}

type userWatchService struct {
	repo repository.UserWatchRepository
}

func (s *userWatchService) Get(id int64) *model.UserWatch {
	return s.repo.Get(id)
}

func (s *userWatchService) Take(where ...interface{}) *model.UserWatch {
	return s.repo.Take(where...)
}

func (s *userWatchService) Find(cnd *querybuilder.QueryBuilder) []model.UserWatch {
	return s.repo.Find(cnd)
}

func (s *userWatchService) FindOne(cnd *querybuilder.QueryBuilder) *model.UserWatch {
	return s.repo.FindOne(cnd)
}

func (s *userWatchService) List(cnd *querybuilder.QueryBuilder) ([]model.UserWatch, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *userWatchService) Create(t *model.UserWatch) error {
	return s.repo.Create(t)
}

func (s *userWatchService) Update(t *model.UserWatch) error {
	return s.repo.Update(t)
}

func (s *userWatchService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *userWatchService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *userWatchService) Delete(id int64) {
	s.repo.Delete(id)
}

// Watch 关注用户
func (s *userWatchService) Watch(userSlug string, userID int64) error {
	watchedUserID := hashid.Slug2Id[model.User](userSlug)
	if watchedUserID <= 0 {
		return errors.New("用户不存在")
	}
	if watchedUserID == userID {
		return errors.New("不能关注自己")
	}
	tmp := s.repo.Take("user_id = ? AND watcher_id = ?", userID, watchedUserID)
	if tmp != nil {
		return nil
	}
	return s.repo.Create(&model.UserWatch{
		UserID:     userID,
		WatcherID:  watchedUserID,
		CreateTime: util.NowTimestamp(),
	})
}

// GetBy 查询关注关系
func (s *userWatchService) GetBy(userID, watchedUserID int64) *model.UserWatch {
	return s.repo.Take("user_id = ? AND watcher_id = ?", userID, watchedUserID)
}
