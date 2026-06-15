package service

import (
	"encoding/json"
	"fmt"
	"errors"
	"strconv"

	"github.com/tidwall/gjson"
	"gorm.io/gorm"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/form"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

var SettingService = newSettingService()

func newSettingService() *settingService {
	return &settingService{}
}

type settingService struct{}

func (s *settingService) Get(id int64) *model.Setting {
	return dao.SettingDao.Get(id)
}

func (s *settingService) Take(where ...interface{}) *model.Setting {
	return dao.SettingDao.Take(where...)
}

func (s *settingService) Find(cnd *querybuilder.QueryBuilder) []model.Setting {
	return dao.SettingDao.Find(cnd)
}

func (s *settingService) FindOne(cnd *querybuilder.QueryBuilder) *model.Setting {
	return dao.SettingDao.FindOne(cnd)
}

func (s *settingService) List(cnd *querybuilder.QueryBuilder) (list []model.Setting, paging *querybuilder.Paging) {
	return dao.SettingDao.List(cnd)
}

func (s *settingService) GetAll() []model.Setting {
	return dao.SettingDao.Find(querybuilder.NewQueryBuilder().Asc("id"))
}

// SetAll 批量设置配置
func (s *settingService) SetAll(configStr string) error {
	json := gjson.Parse(configStr)
	configs, ok := json.Value().(map[string]interface{})
	if !ok {
		return errors.New("配置数据格式错误")
	}

	return dao.DB().Transaction(func(tx *gorm.DB) error {
		for k := range configs {
			v := json.Get(k).String()
			if err := s.setSingle(tx, k, v, "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}

// SetAllFromStruct 接收强类型表单结构体并批量保存配置
func (s *settingService) SetAllFromStruct(req form.SettingsRequest) error {
	// ✅ 将 form 结构体序列化为标准 JSON 字符串
	bytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// ✅ 直接复用现有的 SetAll 逻辑（包含 gjson 解析、事务、setSingle）
	return s.SetAll(string(bytes))
}

// Set 设置单个配置，不存在则创建
func (s *settingService) Set(key, value, name, description string) error {
	return dao.DB().Transaction(func(tx *gorm.DB) error {
		return s.setSingle(tx, key, value, name, description)
	})
}

// setSingle 内部设置单个配置项
func (s *settingService) setSingle(tx *gorm.DB, key, value, name, description string) error {
	if len(key) == 0 {
		return errors.New("sys config key is null")
	}

	// ✅ 使用 tx 查询，而非全局 dao.SettingDao.GetByKey
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
		// ✅ 使用 tx 创建
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
		// ✅ 使用 tx 更新
		if err := tx.Model(&sysConfig).Updates(updates).Error; err != nil {
			return err
		}
	}

	cache.SettingCache.Invalidate(key)
	return nil
}

// GetSetting 获取站点基础配置
func (s *settingService) GetSetting() *model.ConfigData {
	var (
		siteTitle        = cache.SettingCache.GetValue(model.SettingSiteTitle)
		siteDescription  = cache.SettingCache.GetValue(model.SettingSiteDescription)
		defaultNodeIdStr = cache.SettingCache.GetValue(model.SettingDefaultNodeId)
		recommendTags    = cache.SettingCache.GetValue(model.SettingRecommendTags)
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