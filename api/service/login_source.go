package service

import (
	"database/sql"
	"strconv"
	"strings"

	"ultrathreads/model"
	"ultrathreads/oauth/gitee"
	"ultrathreads/oauth/github"
	"ultrathreads/oauth/qq"
	"ultrathreads/repository"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// LoginSourceServicer 登录来源业务契约
type LoginSourceServicer interface {
	Get(id int64) *model.LoginSource
	List(cnd *querybuilder.QueryBuilder) ([]model.LoginSource, *querybuilder.Paging)
	Create(t *model.LoginSource) error
	Update(t *model.LoginSource) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	GetLoginSource(targetType string, targetID string) *model.LoginSource
	GetOrCreate(provider, code, state string) (*model.LoginSource, error)
	GetOrCreateByGithub(code, state string) (*model.LoginSource, error)
	GetOrCreateByGitee(code, state string) (*model.LoginSource, error)
	GetOrCreateByQQ(code, state string) (*model.LoginSource, error)
}

func NewLoginSourceService(repo repository.LoginSourceRepository) LoginSourceServicer {
	return &loginSourceService{repo: repo}
}

type loginSourceService struct {
	repo repository.LoginSourceRepository
}

func (s *loginSourceService) Get(id int64) *model.LoginSource {
	return s.repo.Get(id)
}

func (s *loginSourceService) List(cnd *querybuilder.QueryBuilder) ([]model.LoginSource, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *loginSourceService) Create(t *model.LoginSource) error {
	return s.repo.Create(t)
}

func (s *loginSourceService) Update(t *model.LoginSource) error {
	return s.repo.Update(t)
}

func (s *loginSourceService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *loginSourceService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *loginSourceService) Delete(id int64) {
	s.repo.Delete(id)
}

func (s *loginSourceService) GetLoginSource(targetType string, targetID string) *model.LoginSource {
	return s.repo.Take("target_type = ? and target_id = ?", targetType, targetID)
}

func (s *loginSourceService) GetOrCreate(provider, code, state string) (*model.LoginSource, error) {
	if provider == "github" {
		return s.GetOrCreateByGithub(code, state)
	} else if provider == "gitee" {
		return s.GetOrCreateByGitee(code, state)
	}

	return s.GetOrCreateByQQ(code, state)
}

func (s *loginSourceService) GetOrCreateByGithub(code, state string) (*model.LoginSource, error) {
	userInfo, err := github.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetLoginSource(model.LoginSourceTypeGithub, strconv.FormatInt(userInfo.Id, 10))
	if account != nil {
		return account, nil
	}

	nickname := userInfo.Login
	if len(userInfo.Name) > 0 {
		nickname = strings.TrimSpace(userInfo.Name)
	}

	userInfoJson, _ := util.FormatJson(userInfo)
	account = &model.LoginSource{
		UserID:     sql.NullInt64{},
		Avatar:     userInfo.AvatarUrl,
		Nickname:   nickname,
		TargetType: model.LoginSourceTypeGithub,
		TargetID:   strconv.FormatInt(userInfo.Id, 10),
		ExtraData:  userInfoJson,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *loginSourceService) GetOrCreateByGitee(code, state string) (*model.LoginSource, error) {
	userInfo, err := gitee.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetLoginSource(model.LoginSourceTypeGitee, strconv.FormatInt(int64(userInfo.Id), 10))
	if account != nil {
		return account, nil
	}

	userInfoJson, _ := util.FormatJson(userInfo)
	account = &model.LoginSource{
		UserID:     sql.NullInt64{},
		Avatar:     userInfo.AvatarUrl,
		Nickname:   userInfo.Login,
		TargetType: model.LoginSourceTypeGitee,
		TargetID:   strconv.FormatInt(int64(userInfo.Id), 10),
		ExtraData:  userInfoJson,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *loginSourceService) GetOrCreateByQQ(code, state string) (*model.LoginSource, error) {
	userInfo, err := qq.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetLoginSource(model.LoginSourceTypeQQ, userInfo.Openid)
	if account != nil {
		return account, nil
	}

	userInfoJson, _ := util.FormatJson(userInfo)
	account = &model.LoginSource{
		UserID:     sql.NullInt64{},
		Avatar:     userInfo.FigureurlQQ1,
		Nickname:   userInfo.Nickname,
		TargetType: model.LoginSourceTypeQQ,
		TargetID:   userInfo.Openid,
		ExtraData:  userInfoJson,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
