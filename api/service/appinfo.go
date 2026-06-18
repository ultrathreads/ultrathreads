package service

import (
	"ultrathreads/config"
	"ultrathreads/model"
)

// AppinfoServicer 应用信息业务契约
type AppinfoServicer interface {
	GetAppinfo() *model.AppData
}

func NewAppinfoService() AppinfoServicer {
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