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
}

// List list settings
func (c *SettingController) List(ctx *gin.Context) {
	c.Success(ctx, service.SettingService.GetSetting())
}

func (c *SettingController) Store(ctx *gin.Context) {
	var req dto.SettingsRequest
	// ✅ ShouldBindJSON 自动解析 + 校验
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Fail(ctx, util.NewError(400, "参数格式错误: "+err.Error()))
		return
	}

	// 将结构体传给 service 层，由 service 负责序列化/存储
	if err := service.SettingService.SetAllFromStruct(req); err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	c.Success(ctx, nil)
}