package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/render"
	"ultrathreads/cache"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

type TagController struct {
	BaseController
	tagSvc service.TagServicer
}

func NewTagController(tagSvc service.TagServicer) *TagController {
	return &TagController{tagSvc: tagSvc}
}

// Show 标签详情
func (c *TagController) Show(ctx *gin.Context) {
	var gDto dto.SlugRequest
	if c.BindAndValidate(ctx, &gDto) {
		tag := c.tagSvc.GetBySlug(gDto.Slug)
		if tag == nil {
			c.Fail(ctx, util.ErrorTagNotFound)
			return
		}
		c.Success(ctx, render.ToTag(tag))
	}
}

// List 标签列表
func (c *TagController) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)

	tags, paging := c.tagSvc.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).
		Page(page, 200).Desc("id"))

	data := map[string]interface{}{}
	data["results"] = render.ToTags(tags)
	data["page"] = paging
	c.Success(ctx, data)
}

// AutoComplete 标签自动完成
func (c *TagController) AutoComplete(ctx *gin.Context) {
	input := util.FormStringDefault(ctx, "input","")
	tags := c.tagSvc.AutoComplete(input)
	c.Success(ctx, tags)
}

// HotTags 热门标签
func (c *TagController) HotTags(ctx *gin.Context) {
	tags := cache.TagCache.GetHot()

	c.Success(ctx, render.ToTags(tags))
}