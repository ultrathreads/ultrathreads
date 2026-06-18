package service

import (
	"errors"
	"strings"

	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// LinkServicer 链接业务契约
type LinkServicer interface {
	Get(id int64) *model.Link
	Find(cnd *querybuilder.QueryBuilder) []model.Link
	List(cnd *querybuilder.QueryBuilder) ([]model.Link, *querybuilder.Paging)
	Create(req dto.LinkCreateForm) (*model.Link, error)
	Update(req dto.LinkUpdateForm) error
	Delete(id int64) error
	Submit(url, title, summary, logo string) (*model.Link, error)
}

func NewLinkService(repo repository.LinkRepository) LinkServicer {
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

func (s *linkService) Create(req dto.LinkCreateForm) (*model.Link, error) {
	link := &model.Link{
		Title:      req.Title,
		Url:        req.URL,
		Summary:    req.Summary,
		Logo:       req.Logo,
		CreateTime: util.NowTimestamp(),
	}
	if err := s.repo.Create(link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *linkService) Update(req dto.LinkUpdateForm) error {
	return s.repo.Updates(req.ID, map[string]interface{}{
		"title":       req.Title,
		"url":         req.URL,
		"summary":     req.Summary,
		"logo":        req.Logo,
		"status":      req.Status,
		"update_time": util.NowTimestamp(),
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
		Url:        url,
		Title:      title,
		Summary:    summary,
		Logo:       logo,
		Status:     model.StatusPending,
		CreateTime: util.NowTimestamp(),
	}

	if err := s.repo.Create(link); err != nil {
		return nil, err
	}
	return link, nil
}
