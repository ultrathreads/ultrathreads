package service

import (
	"errors"
	"strings"

	"ultrathreads/dao"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

var LinkService = newLinkService()

func newLinkService() *linkService {
	return &linkService{}
}

type linkService struct{}

func (s *linkService) Get(id int64) *model.Link {
	return dao.LinkDao.Get(id)
}

func (s *linkService) Find(cnd *querybuilder.QueryBuilder) []model.Link {
	return dao.LinkDao.Find(cnd)
}

func (s *linkService) List(cnd *querybuilder.QueryBuilder) (list []model.Link, paging *querybuilder.Paging) {
	return dao.LinkDao.List(cnd)
}

// Create 创建链接
// ✅ 移除无效事务包装：单条 Create 无需事务，且原代码 DAO 内部用全局 db 导致 tx 未生效
func (s *linkService) Create(req dto.LinkCreateForm) (*model.Link, error) {
	link := &model.Link{
		Title:      req.Title,
		Url:        req.URL,
		Summary:    req.Summary,
		Logo:       req.Logo,
		CreateTime: util.NowTimestamp(),
	}
	if err := dao.LinkDao.Create(link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *linkService) Update(req dto.LinkUpdateForm) error {
	return dao.LinkDao.Updates(req.ID, map[string]interface{}{
		"title":       req.Title,
		"url":         req.URL,
		"summary":     req.Summary,
		"logo":        req.Logo,
		"status":      req.Status,
		"update_time": util.NowTimestamp(),
	})
}

// Delete 删除链接
func (s *linkService) Delete(id int64) error {
	return dao.LinkDao.Delete(id)
}

// Submit 提交友情链接
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

	if err := dao.LinkDao.Create(link); err != nil {
		return nil, err
	}
	return link, nil
}