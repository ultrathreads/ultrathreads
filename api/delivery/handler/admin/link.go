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

// LinkHandler link controller
type LinkHandler struct {
	base.BaseHandler
	linkSvc service.LinkService
}

func NewLinkHandler(linkSvc service.LinkService) *LinkHandler {
	return &LinkHandler{linkSvc: linkSvc}
}

// Show show link
func (h *LinkHandler) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		link := h.linkSvc.Get(gDto.ID)
		if link == nil {
			h.Fail(ctx, util.NewErrorMsg("Link not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		h.Success(ctx, link)
	}
}

// Store create a link
func (h *LinkHandler) Store(ctx *gin.Context) {
	var linkForm dto.LinkCreateForm
	if !h.BindAndValidate(ctx, &linkForm) {
		return
	}
	link, err := h.linkSvc.Create(linkForm)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, link)
}

// Update update a link
func (h *LinkHandler) Update(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	link := h.linkSvc.Get(gDto.ID)
	if link == nil {
		h.Fail(ctx, util.NewErrorMsg("Link not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var linkForm dto.LinkUpdateForm
	if !h.BindAndValidate(ctx, &linkForm) {
		return
	}
	linkForm.ID = gDto.ID
	err := h.linkSvc.Update(linkForm)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, link)
}

// Delete delete link
func (h *LinkHandler) Delete(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	h.linkSvc.Delete(gDto.ID)
	h.Success(ctx, nil)
}

// List list links
func (h *LinkHandler) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	name := ctx.Request.FormValue("name")

	conditions := querybuilder.NewQueryBuilder()
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := h.linkSvc.List(conditions.Page(page, limit).Desc("id"))

	h.Success(ctx, &querybuilder.PageResult{Results: list, Page: paging})
}
