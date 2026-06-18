package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/tidwall/gjson"
	"gorm.io/gorm"

	"ultrathreads/cache"
	"ultrathreads/dto"
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
	SetAllFromStruct(req dto.SettingsRequest) error
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

	return s.repo.Transaction(func(tx *gorm.DB) error {
		for k := range configs {
			v := json.Get(k).String()
			if err := s.setSingle(tx, k, v, "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *settingService) SetAllFromStruct(req dto.SettingsRequest) error {
	bytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	return s.SetAll(string(bytes))
}

func (s *settingService) Set(key, value, name, description string) error {
	return s.repo.Transaction(func(tx *gorm.DB) error {
		return s.setSingle(tx, key, value, name, description)
	})
}

func (s *settingService) setSingle(tx *gorm.DB, key, value, name, description string) error {
	if len(key) == 0 {
		return errors.New("sys config key is null")
	}

	var sysConfig model.Setting
	err := tx.Where("`key` = ?", key).First(&sysConfig).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil && !notFound {
		return err
	}

	if notFound {
		sysConfig = model.Setting{
			Key:        key,
			Value:      value,
			CreateTime: util.NowTimestamp(),
			UpdateTime: util.NowTimestamp(),
		}
		if len(name) > 0 {
			sysConfig.Name = name
		}
		if len(description) > 0 {
			sysConfig.Description = description
		}
		if err := tx.Create(&sysConfig).Error; err != nil {
			return err
		}
	} else {
		updates := map[string]interface{}{
			"value":       value,
			"update_time": util.NowTimestamp(),
		}
		if len(name) > 0 {
			updates["name"] = name
		}
		if len(description) > 0 {
			updates["description"] = description
		}
		if err := tx.Model(&sysConfig).Updates(updates).Error; err != nil {
			return err
		}
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
