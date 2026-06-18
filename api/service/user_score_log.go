package service

import (
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// UserScoreLogServicer 用户积分日志业务契约
type UserScoreLogServicer interface {
	Get(id int64) *model.UserScoreLog
	Take(where ...interface{}) *model.UserScoreLog
	Find(cnd *querybuilder.QueryBuilder) []model.UserScoreLog
	FindOne(cnd *querybuilder.QueryBuilder) *model.UserScoreLog
	List(cnd *querybuilder.QueryBuilder) ([]model.UserScoreLog, *querybuilder.Paging)
	Create(t *model.UserScoreLog) error
	Update(t *model.UserScoreLog) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
}

func NewUserScoreLogService(repo dao.UserScoreLogRepository) UserScoreLogServicer {
	return &userScoreLogService{repo: repo}
}

type userScoreLogService struct {
	repo dao.UserScoreLogRepository
}

func (s *userScoreLogService) Get(id int64) *model.UserScoreLog {
	return s.repo.Get(id)
}

func (s *userScoreLogService) Take(where ...interface{}) *model.UserScoreLog {
	return s.repo.Take(where...)
}

func (s *userScoreLogService) Find(cnd *querybuilder.QueryBuilder) []model.UserScoreLog {
	return s.repo.Find(cnd)
}

func (s *userScoreLogService) FindOne(cnd *querybuilder.QueryBuilder) *model.UserScoreLog {
	return s.repo.FindOne(cnd)
}

func (s *userScoreLogService) List(cnd *querybuilder.QueryBuilder) ([]model.UserScoreLog, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *userScoreLogService) Create(t *model.UserScoreLog) error {
	return s.repo.Create(t)
}

func (s *userScoreLogService) Update(t *model.UserScoreLog) error {
	return s.repo.Update(t)
}

func (s *userScoreLogService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *userScoreLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *userScoreLogService) Delete(id int64) {
	s.repo.Delete(id)
}