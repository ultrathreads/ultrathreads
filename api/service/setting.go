package service

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

var SettingService = newSettingService()

func newSettingService() *settingService {
	return &settingService{}
}

type settingService struct {
}

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

func (s *settingService) SetAll(configStr string) error {
	json := gjson.Parse(configStr)
	configs, ok := json.Value().(map[string]interface{})
	if !ok {
		return errors.New("配置数据格式错误")
	}
	return dao.Tx(dao.DB(), func(tx *gorm.DB) error {
		for k := range configs {
			v := json.Get(k).String()
			if err := s.setSingle(tx, k, v, "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}

// 设置配置，如果配置不存在，那么创建
func (s *settingService) Set(key, value, name, description string) error {
	return dao.Tx(dao.DB(), func(tx *gorm.DB) error {
		if err := s.setSingle(tx, key, value, name, description); err != nil {
			return err
		}
		return nil
	})
}

func (s *settingService) setSingle(db *gorm.DB, key, value, name, description string) error {
	if len(key) == 0 {
		return errors.New("sys config key is null")
	}
	sysConfig := dao.SettingDao.GetByKey(key)
	if sysConfig == nil {
		sysConfig = &model.Setting{
			CreateTime: util.NowTimestamp(),
		}
	}
	sysConfig.Key = key
	sysConfig.Value = value
	sysConfig.UpdateTime = util.NowTimestamp()

	if len(name) > 0 {
		sysConfig.Name = name
	}
	if len(description) > 0 {
		sysConfig.Description = description
	}

	var err error
	if sysConfig.ID > 0 {
		err = dao.SettingDao.Update(sysConfig)
	} else {
		err = dao.SettingDao.Create(sysConfig)
	}
	if err != nil {
		return err
	} else {
		cache.SettingCache.Invalidate(key)
		return nil
	}
}

func (s *settingService) GetSetting() *model.ConfigData {
	var (
		siteTitle        = cache.SettingCache.GetValue(model.SettingSiteTitle)
		siteDescription  = cache.SettingCache.GetValue(model.SettingSiteDescription)
		defaultNodeIdStr = cache.SettingCache.GetValue(model.SettingDefaultNodeId)
	)


	var defaultNodeId, _ = strconv.ParseInt(defaultNodeIdStr, 10, 64)

	return &model.ConfigData{
		SiteTitle:        siteTitle,
		SiteDescription:  siteDescription,
		DefaultNodeId:    defaultNodeId,
	}
}
