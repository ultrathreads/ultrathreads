package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
	"ultrathreads/service"
	"ultrathreads/model"
	"ultrathreads/util/hashid"
)

type SiteController struct {
	BaseController
	settingSvc        service.SettingServicer
	appinfoSvc        service.AppinfoServicer
	userReadStateSvc  service.UserReadStateServicer
}

func NewSiteController(settingSvc service.SettingServicer, appinfoSvc service.AppinfoServicer, userReadStateSvc service.UserReadStateServicer) *SiteController {
	return &SiteController{settingSvc: settingSvc, appinfoSvc: appinfoSvc, userReadStateSvc: userReadStateSvc}
}

func (c *SiteController) Debug(ctx *gin.Context) {
	userID := int64(1)
	states := c.userReadStateSvc.GetUserReadStates(userID)

	c.Success(ctx, states)
}

func (c *SiteController) Config(ctx *gin.Context) {
	data := map[string]interface{}{}
	data["setting"] = c.settingSvc.GetSetting()
	data["appinfo"] = c.appinfoSvc.GetAppinfo()

	c.Success(ctx, data)
}

func (c *SiteController) Stat(ctx *gin.Context) {
	data := make(map[string]interface{})
	data["userCount"] = cache.StatCache.GetUserCount()
	data["postCount"] = cache.StatCache.GetPostCount()

	c.Success(ctx, data)
}

// Ping 健康检查 - 通过事件总线异步输出 pong
func (c *SiteController) Ping(ctx *gin.Context) {
	slug := hashid.Id2Slug[model.Node](123)

	id := hashid.Slug2Id[model.Node]("4q7ZxEa5")
	c.Success(ctx, gin.H{
		"message": "pong",
		"slug": slug,
		"id": id,
	})
}