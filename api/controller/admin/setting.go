package admin

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/controller"
	"ultrathreads/util"
	"ultrathreads/dto"
	"ultrathreads/service"
)

// SettingController setting controller
type SettingController struct {
	controller.BaseController
	settingSvc service.SettingServicer
}

func NewSettingController(settingSvc service.SettingServicer) *SettingController {
	return &SettingController{settingSvc: settingSvc}
}

// List list settings
func (c *SettingController) List(ctx *gin.Context) {
	c.Success(ctx, c.settingSvc.GetSetting())
}

func (c *SettingController) Store(ctx *gin.Context) {
	var req dto.SettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Fail(ctx, util.NewError(400, "参数格式错误: "+err.Error()))
		return
	}

	if err := c.settingSvc.SetAllFromStruct(req); err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	c.Success(ctx, nil)
}