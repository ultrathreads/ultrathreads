package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/converter"
	"ultrathreads/cache"
	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

type TagController struct {
	BaseController
}

// Show 标签详情
func (c *TagController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		tag := cache.TagCache.Get(gDto.ID)
		if tag == nil {
			c.Fail(ctx, util.ErrorTagNotFound)
			return
		}
		c.Success(ctx, converter.ToTag(tag))
	}
}

// List 标签列表
func (c *TagController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)

	tags, paging := service.TagService.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).
		Page(page, 200).Desc("id"))

	data := map[string]interface{}{}
	data["results"] = converter.ToTags(tags)
	data["page"] = paging
	c.Success(ctx, data)
}

// AutoComplete 标签自动完成
func (c *TagController) AutoComplete(ctx *gin.Context) {
	input := ctx.Request.FormValue("input")
	tags := service.TagService.AutoComplete(input)
	c.Success(ctx, tags)
}

// HotTags 标签自动完成
func (c *TagController) HotTags(ctx *gin.Context) {
	tags := cache.TagCache.GetHot()

	c.Success(ctx, converter.ToTags(tags))
}
