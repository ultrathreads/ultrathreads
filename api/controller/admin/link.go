package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"ultrathreads/controller"
	"ultrathreads/dto"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// LinkController link controller
type LinkController struct {
	controller.BaseController
	linkSvc service.LinkServicer
}

func NewLinkController(linkSvc service.LinkServicer) *LinkController {
	return &LinkController{linkSvc: linkSvc}
}

// Show show link
func (c *LinkController) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		link := c.linkSvc.Get(gDto.ID)
		if link == nil {
			c.Fail(ctx, util.NewErrorMsg("Link not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, link)
	}
}

// Store create a link
func (c *LinkController) Store(ctx *gin.Context) {
	var linkForm dto.LinkCreateForm
	if !c.BindAndValidate(ctx, &linkForm) {
		return
	}
	link, err := c.linkSvc.Create(linkForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, link)
}

// Update update a link
func (c *LinkController) Update(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	link := c.linkSvc.Get(gDto.ID)
	if link == nil {
		c.Fail(ctx, util.NewErrorMsg("Link not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var linkForm dto.LinkUpdateForm
	if !c.BindAndValidate(ctx, &linkForm) {
		return
	}
	linkForm.ID = gDto.ID
	err := c.linkSvc.Update(linkForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, link)
}

// Delete delete link
func (c *LinkController) Delete(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	c.linkSvc.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// List list links
func (c *LinkController) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	name := ctx.Request.FormValue("name")

	conditions := querybuilder.NewQueryBuilder()
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := c.linkSvc.List(conditions.Page(page, limit).Desc("id"))

	c.Success(ctx, &querybuilder.PageResult{Results: list, Page: paging})
}