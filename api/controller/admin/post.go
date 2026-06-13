package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"ultrathreads/converter"
	"ultrathreads/controller"
	"ultrathreads/form"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/model"
	"ultrathreads/util/hashid"
	"ultrathreads/util/markdown"
	"ultrathreads/util/querybuilder"
)

// PostController post controller
type PostController struct {
	controller.BaseController
}

// Show show post
func (c *PostController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		post := service.PostService.Get(gDto.ID)
		if post == nil {
			c.Fail(ctx, util.NewErrorMsg("Post not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, post)
	}
}

// Update update a post
func (c *PostController) Update(ctx *gin.Context) {
	var gDto form.IdentifierDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	postID := hashid.Slug2Id[model.Post](gDto.Slug)
	post := service.PostService.Get(postID)
	if post == nil {
		c.Fail(ctx, util.NewErrorMsg("Post not found, id="+strconv.FormatInt(postID, 10)))
		return
	}

	var postForm form.RootPostUpdateForm
	if !c.BindAndValidate(ctx, &postForm) {
		return
	}
	postForm.Slug = gDto.Slug
	err := service.PostService.UpdateRootPost(postForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, post)
}

// Delete delete post
func (c *PostController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.PostService.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// Undelete delete post
func (c *PostController) Undelete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.PostService.Undelete(gDto.ID)
	c.Success(ctx, nil)
}

// List list posts
func (c *PostController) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	id := ctx.Request.FormValue("id")
	userID := ctx.Request.FormValue("user_id")
	status := ctx.Request.FormValue("status")
	recommend := ctx.Request.FormValue("recommend")
	title := ctx.Request.FormValue("title")

	conditions := querybuilder.NewQueryBuilder()
	if len(id) > 0 {
		conditions.Eq("id", id)
	}
	if len(userID) > 0 {
		conditions.Eq("user_id", userID)
	}
	if len(status) > 0 {
		conditions.Eq("status", status)
	}
	if len(recommend) > 0 {
		conditions.Eq("recommend", recommend)
	}
	if len(title) > 0 {
		conditions.Like("title", title)
	}

	list, paging := service.PostService.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, post := range list {
		result := util.StructToMap(post, "content")
		result["user"] = converter.ToUserDefaultIfNull(post.UserId)
		result["node"] = service.NodeService.Get(post.NodeId)
		result["tags"] = converter.ToTags(service.PostService.GetPostTags(post.ID))
		// 简介
		mr := markdown.NewMd().Run(post.Content)
		result["summary"] = mr.SummaryText

		results = append(results, result)
	}

	c.Success(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}

// Recommend 推荐
func (c *PostController) Recommend(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	err := service.PostService.SetRecommend(gDto.ID, true)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// Unrecommend 取消推荐
func (c *PostController) Unrecommend(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	err := service.PostService.SetRecommend(gDto.ID, false)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}
