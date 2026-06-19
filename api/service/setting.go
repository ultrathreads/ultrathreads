package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/tidwall/gjson"

	"ultrathreads/cache"
	"ultrathreads/domain"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// SettingService 系统设置业务契约
type SettingService interface {
	Get(id int64) *model.Setting
	Take(where ...interface{}) *model.Setting
	Find(cnd *querybuilder.QueryBuilder) []model.Setting
	FindOne(cnd *querybuilder.QueryBuilder) *model.Setting
	List(cnd *querybuilder.QueryBuilder) ([]model.Setting, *querybuilder.Paging)
	GetAll() []model.Setting
	SetAll(configStr string) error
	SetAllFromStruct(cmd domain.UpdateSettingsCommand) error
	Set(key, value, name, description string) error
	GetSetting() *model.ConfigData
}

func NewSettingService(repo repository.SettingRepository, settingCache cache.SettingCacheInterface) SettingService {
	return &settingService{repo: repo, settingCache: settingCache}
}

type settingService struct {
	repo         repository.SettingRepository
	settingCache cache.SettingCacheInterface
}

func (s *settingService) Get(id int64) *model.Setting {
	return s.repo.Get(id)
}

func (s *settingService) Take(where ...interface{}) *model.Setting {
	return s.repo.Take(where...)
}

func (s *settingService) Find(cnd *querybuilder.QueryBuilder) []model.Setting {
	return s.repo.Find(cnd)
}

func (s *settingService) FindOne(cnd *querybuilder.QueryBuilder) *model.Setting {
	return s.repo.FindOne(cnd)
}

func (s *settingService) List(cnd *querybuilder.QueryBuilder) ([]model.Setting, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *settingService) GetAll() []model.Setting {
	return s.repo.Find(querybuilder.NewQueryBuilder().Asc("id"))
}

func (s *settingService) SetAll(configStr string) error {
	json := gjson.Parse(configStr)
	configs, ok := json.Value().(map[string]interface{})
	if !ok {
		return errors.New("配置数据格式错误")
	}

	if err := s.repo.UpsertAll(configs); err != nil {
		return err
	}

	// 失效缓存
	for k := range configs {
		s.settingCache.Invalidate(k)
	}
	return nil
}

func (s *settingService) SetAllFromStruct(cmd domain.UpdateSettingsCommand) error {
	bytes, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	return s.SetAll(string(bytes))
}

func (s *settingService) Set(key, value, name, description string) error {
	if err := s.repo.UpsertByKey(key, value, name, description); err != nil {
		return err
	}
	s.settingCache.Invalidate(key)
	return nil
}

func (s *settingService) GetSetting() *model.ConfigData {
	var (
		siteTitle        = s.settingCache.GetValue(model.SettingSiteTitle)
		siteDescription  = s.settingCache.GetValue(model.SettingSiteDescription)
		defaultNodeIdStr = s.settingCache.GetValue(model.SettingDefaultNodeId)
		recommendTags    = s.settingCache.GetValue(model.SettingRecommendTags)
	)

	var recommendTagsArr []string
	if err := util.ParseJson(recommendTags, &recommendTagsArr); err != nil {
		recommendTagsArr = []string{}
	}

	defaultNodeId, _ := strconv.ParseInt(defaultNodeIdStr, 10, 64)

	return &model.ConfigData{
		SiteTitle:       siteTitle,
		SiteDescription: siteDescription,
		DefaultNodeId:   defaultNodeId,
		RecommendTags:   recommendTagsArr,
	}
}
