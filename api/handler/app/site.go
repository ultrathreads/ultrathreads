package app

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
	"ultrathreads/handler/base"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util/hashid"
)

type SiteHandler struct {
	base.BaseHandler
	settingSvc       service.SettingServicer
	appinfoSvc       service.AppinfoServicer
	userReadStateSvc service.UserReadStateServicer
}

func NewSiteHandler(settingSvc service.SettingServicer, appinfoSvc service.AppinfoServicer, userReadStateSvc service.UserReadStateServicer) *SiteHandler {
	return &SiteHandler{settingSvc: settingSvc, appinfoSvc: appinfoSvc, userReadStateSvc: userReadStateSvc}
}

func (h *SiteHandler) Debug(ctx *gin.Context) {
	userID := int64(1)
	states := h.userReadStateSvc.GetUserReadStates(userID)

	h.Success(ctx, states)
}

func (h *SiteHandler) Config(ctx *gin.Context) {
	data := map[string]interface{}{}
	data["setting"] = h.settingSvc.GetSetting()
	data["appinfo"] = h.appinfoSvc.GetAppinfo()

	h.Success(ctx, data)
}

func (h *SiteHandler) Stat(ctx *gin.Context) {
	data := make(map[string]interface{})
	data["userCount"] = cache.StatCache.GetUserCount()
	data["postCount"] = cache.StatCache.GetPostCount()

	h.Success(ctx, data)
}

// Ping 健康检查 - 通过事件总线异步输出 pong
func (h *SiteHandler) Ping(ctx *gin.Context) {
	slug := hashid.Id2Slug[model.Node](123)

	id := hashid.Slug2Id[model.Node]("4q7ZxEa5")
	h.Success(ctx, gin.H{
		"message": "pong",
		"slug":    slug,
		"id":      id,
	})
}
