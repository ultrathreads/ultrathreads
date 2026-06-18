package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ultrathreads/dto"
	"ultrathreads/delivery/handler/base"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// TagHandler tag controller
type TagHandler struct {
	base.BaseHandler
	tagSvc service.TagService
}

func NewTagHandler(tagSvc service.TagService) *TagHandler {
	return &TagHandler{tagSvc: tagSvc}
}

// Show show tag
func (h *TagHandler) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		tag := h.tagSvc.Get(gDto.ID)
		if tag == nil {
			h.Fail(ctx, util.NewErrorMsg("Tag not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		h.Success(ctx, tag)
	}
}

// Store create a tag
func (h *TagHandler) Store(ctx *gin.Context) {
	var tagForm dto.TagCreateForm
	if !h.BindAndValidate(ctx, &tagForm) {
		return
	}
	tag, err := h.tagSvc.Create(tagForm)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, tag)
}

// Update update a tag
func (h *TagHandler) Update(ctx *gin.Context) {
	var tagForm dto.TagUpdateForm
	if !h.BindAndValidate(ctx, &tagForm) {
		return
	}
	tag := h.tagSvc.Get(tagForm.ID)
	if tag == nil {
		h.Fail(ctx, util.NewErrorMsg("Tag not found, id="+strconv.FormatInt(tagForm.ID, 10)))
		return
	}

	err := h.tagSvc.Update(tagForm.ID, tagForm)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, tag)
}

// Delete delete tag
func (h *TagHandler) Delete(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	h.tagSvc.Delete(gDto.ID)
	h.Success(ctx, nil)
}

// List list tags
func (h *TagHandler) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	name := ctx.Request.FormValue("name")

	conditions := querybuilder.NewQueryBuilder()
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := h.tagSvc.List(conditions.Page(page, limit).Desc("id"))

	h.Success(ctx, &querybuilder.PageResult{Results: list, Page: paging})
}
