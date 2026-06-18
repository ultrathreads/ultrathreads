package service

import (
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
)

// UserReadStateServicer 用户阅读状态业务契约
type UserReadStateServicer interface {
	Get(id int64) *model.UserReadState
	Take(where ...interface{}) *model.UserReadState
	Find(cnd *querybuilder.QueryBuilder) []model.UserReadState
	FindOne(cnd *querybuilder.QueryBuilder) *model.UserReadState
	List(cnd *querybuilder.QueryBuilder) ([]model.UserReadState, *querybuilder.Paging)
	Create(t *model.UserReadState) error
	Update(t *model.UserReadState) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	GetUserReadStates(userID int64) map[int64]int64
	MarkAsRead(userID int64, nodeSlug string, now int64) error
}

func NewUserReadStateService(repo dao.UserReadStateRepository) UserReadStateServicer {
	return &userReadStateService{repo: repo}
}

type userReadStateService struct {
	repo dao.UserReadStateRepository
}

func (s *userReadStateService) Get(id int64) *model.UserReadState {
	return s.repo.Get(id)
}

func (s *userReadStateService) Take(where ...interface{}) *model.UserReadState {
	return s.repo.Take(where...)
}

func (s *userReadStateService) Find(cnd *querybuilder.QueryBuilder) []model.UserReadState {
	return s.repo.Find(cnd)
}

func (s *userReadStateService) FindOne(cnd *querybuilder.QueryBuilder) *model.UserReadState {
	return s.repo.FindOne(cnd)
}

func (s *userReadStateService) List(cnd *querybuilder.QueryBuilder) ([]model.UserReadState, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *userReadStateService) Create(t *model.UserReadState) error {
	return s.repo.Create(t)
}

func (s *userReadStateService) Update(t *model.UserReadState) error {
	return s.repo.Update(t)
}

func (s *userReadStateService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *userReadStateService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *userReadStateService) Delete(id int64) {
	s.repo.Delete(id)
}

// GetUserReadStates 获取用户所有节点的已读状态
func (s *userReadStateService) GetUserReadStates(userID int64) map[int64]int64 {
	return s.repo.GetAllReadStates(userID)
}

// MarkAsRead 标记用户在指定节点已读
func (s *userReadStateService) MarkAsRead(userID int64, nodeSlug string, now int64) error {
	nodeID := hashid.Slug2Id[model.Node](nodeSlug)
	if nodeID <= 0 {
		return nil
	}
	return s.repo.Upsert(userID, nodeID, now)
}
