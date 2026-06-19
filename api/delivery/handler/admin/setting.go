package admin

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/delivery/handler/base"
	"ultrathreads/domain"
	"ultrathreads/dto"
	"ultrathreads/service"
	"ultrathreads/util"
)

// SettingHandler setting controller
type SettingHandler struct {
	base.BaseHandler
	settingSvc service.SettingService
}

func NewSettingHandler(settingSvc service.SettingService) *SettingHandler {
	return &SettingHandler{settingSvc: settingSvc}
}

// List list settings
func (h *SettingHandler) List(ctx *gin.Context) {
	h.Success(ctx, h.settingSvc.GetSetting())
}

func (h *SettingHandler) Store(ctx *gin.Context) {
	var req dto.SettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Fail(ctx, util.NewError(400, "参数格式错误: "+err.Error()))
		return
	}

	cmd := domain.UpdateSettingsCommand{
		SiteTitle:       req.SiteTitle,
		SiteDescription: req.SiteDescription,
		SiteKeywords:    req.SiteKeywords,
		SiteNavs:        req.SiteNavs,
		DefaultNodeId:   req.DefaultNodeId,
		RecommendTags:   req.RecommendTags,
	}
	if err := h.settingSvc.SetAllFromStruct(cmd); err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.Success(ctx, nil)
}
