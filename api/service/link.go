package service

import (
	"errors"
	"strings"
	"time"

	"ultrathreads/domain"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util/querybuilder"
)

// LinkService 链接业务契约
type LinkService interface {
	Get(id int64) *model.Link
	Find(cnd *querybuilder.QueryBuilder) []model.Link
	List(cnd *querybuilder.QueryBuilder) ([]model.Link, *querybuilder.Paging)
	Create(cmd domain.CreateLinkCommand) (*model.Link, error)
	Update(cmd domain.UpdateLinkCommand) error
	Delete(id int64) error
	Submit(url, title, summary, logo string) (*model.Link, error)
}

func NewLinkService(repo repository.LinkRepository) LinkService {
	return &linkService{repo: repo}
}

type linkService struct {
	repo repository.LinkRepository
}

func (s *linkService) Get(id int64) *model.Link {
	return s.repo.Get(id)
}

func (s *linkService) Find(cnd *querybuilder.QueryBuilder) []model.Link {
	return s.repo.Find(cnd)
}

func (s *linkService) List(cnd *querybuilder.QueryBuilder) ([]model.Link, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *linkService) Create(cmd domain.CreateLinkCommand) (*model.Link, error) {
	link := &model.Link{
		Title:     cmd.Title,
		Url:       cmd.URL,
		Summary:   cmd.Summary,
		Logo:      cmd.Logo,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *linkService) Update(cmd domain.UpdateLinkCommand) error {
	return s.repo.Updates(cmd.ID, map[string]interface{}{
		"title":      cmd.Title,
		"url":        cmd.URL,
		"summary":    cmd.Summary,
		"logo":       cmd.Logo,
		"status":     cmd.Status,
		"updated_at": time.Now(),
	})
}

func (s *linkService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *linkService) Submit(url, title, summary, logo string) (*model.Link, error) {
	url = strings.TrimSpace(url)
	title = strings.TrimSpace(title)
	summary = strings.TrimSpace(summary)
	logo = strings.TrimSpace(logo)

	if len(url) == 0 {
		return nil, errors.New("网址不能为空")
	}
	if len(title) == 0 {
		return nil, errors.New("标题不能为空")
	}

	link := &model.Link{
		Url:       url,
		Title:     title,
		Summary:   summary,
		Logo:      logo,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(link); err != nil {
		return nil, err
	}
	return link, nil
}
