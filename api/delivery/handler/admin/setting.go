package admin

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/dto"
	"ultrathreads/delivery/handler/base"
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

	if err := h.settingSvc.SetAllFromStruct(req); err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.Success(ctx, nil)
}
