package service

import (
	"ultrathreads/config"
	"ultrathreads/model"
)

// AppinfoService 应用信息业务契约
type AppinfoService interface {
	GetAppinfo() *model.AppData
}

func NewAppinfoService() AppinfoService {
	return &appinfoService{}
}

type appinfoService struct {
}

func (s *appinfoService) GetAppinfo() *model.AppData {
	return &model.AppData{
		Name:           config.AppName,
		Version:        config.AppVersion,
		UserLevelAdmin: model.UserLevelAdmin,
	}
}