package app

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
	"ultrathreads/dto"
	"ultrathreads/handler/base"
	"ultrathreads/model"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

type TagHandler struct {
	base.BaseHandler
	tagSvc   service.TagServicer
	tagCache cache.TagCacheInterface
}

func NewTagHandler(tagSvc service.TagServicer, tagCache cache.TagCacheInterface) *TagHandler {
	return &TagHandler{tagSvc: tagSvc, tagCache: tagCache}
}

// Show 标签详情
func (h *TagHandler) Show(ctx *gin.Context) {
	var gDto dto.SlugRequest
	if h.BindAndValidate(ctx, &gDto) {
		tag := h.tagSvc.GetBySlug(gDto.Slug)
		if tag == nil {
			h.Fail(ctx, util.ErrorTagNotFound)
			return
		}
		h.Success(ctx, render.ToTag(tag))
	}
}

// List 标签列表
func (h *TagHandler) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)

	tags, paging := h.tagSvc.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).
		Page(page, 200).Desc("id"))

	data := map[string]interface{}{}
	data["results"] = render.ToTags(tags)
	data["page"] = paging
	h.Success(ctx, data)
}

// AutoComplete 标签自动完成
func (h *TagHandler) AutoComplete(ctx *gin.Context) {
	input := util.FormStringDefault(ctx, "input", "")
	tags := h.tagSvc.AutoComplete(input)
	h.Success(ctx, tags)
}

// HotTags 热门标签
func (h *TagHandler) HotTags(ctx *gin.Context) {
	tags := h.tagCache.GetHot()

	h.Success(ctx, render.ToTags(tags))
}
