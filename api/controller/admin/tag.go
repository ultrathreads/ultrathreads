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

// TagController tag controller
type TagController struct {
	controller.BaseController
}

// Show show tag
func (c *TagController) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		tag := service.TagService.Get(gDto.ID)
		if tag == nil {
			c.Fail(ctx, util.NewErrorMsg("Tag not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, tag)
	}
}

// Store create a tag
func (c *TagController) Store(ctx *gin.Context) {
	var tagForm dto.TagCreateForm
	if !c.BindAndValidate(ctx, &tagForm) {
		return
	}
	tag, err := service.TagService.Create(tagForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, tag)
}

// Update update a tag
func (c *TagController) Update(ctx *gin.Context) {
	var tagForm dto.TagUpdateForm
	if !c.BindAndValidate(ctx, &tagForm) {
		return
	}
	tag := service.TagService.Get(tagForm.ID)
	if tag == nil {
		c.Fail(ctx, util.NewErrorMsg("Tag not found, id="+strconv.FormatInt(tagForm.ID, 10)))
		return
	}

	err := service.TagService.Update(tagForm.ID, tagForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, tag)
}

// Delete delete tag
func (c *TagController) Delete(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.TagService.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// List list tags
func (c *TagController) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	name := ctx.Request.FormValue("name")

	conditions := querybuilder.NewQueryBuilder()
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := service.TagService.List(conditions.Page(page, limit).Desc("id"))

	c.Success(ctx, &querybuilder.PageResult{Results: list, Page: paging})
}
