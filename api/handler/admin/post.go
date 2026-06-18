package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ultrathreads/domain"
	"ultrathreads/dto"
	"ultrathreads/handler/base"
	"ultrathreads/model"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/markdown"
	"ultrathreads/util/querybuilder"
)

// PostHandler post controller
type PostHandler struct {
	base.BaseHandler
	postSvc service.PostServicer
	nodeSvc service.NodeServicer
}

func NewPostHandler(postSvc service.PostServicer, nodeSvc service.NodeServicer) *PostHandler {
	return &PostHandler{postSvc: postSvc, nodeSvc: nodeSvc}
}

// Show show post
func (h *PostHandler) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		post := h.postSvc.Get(gDto.ID)
		if post == nil {
			h.Fail(ctx, util.NewErrorMsg("Post not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		h.Success(ctx, post)
	}
}

// Update update a post
func (h *PostHandler) Update(ctx *gin.Context) {
	var gDto dto.SlugRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	postID := hashid.Slug2Id[model.Post](gDto.Slug)
	post := h.postSvc.Get(postID)
	if post == nil {
		h.Fail(ctx, util.NewErrorMsg("Post not found, id="+strconv.FormatInt(postID, 10)))
		return
	}

	var req dto.UpdateRootPostRequest
	if !h.BindAndValidate(ctx, &req) {
		return
	}
	req.Slug = gDto.Slug

	cmd := domain.UpdatePostCommand{
		Slug:     req.Slug,
		Title:    req.Title,
		Content:  req.Content,
		NodeSlug: req.NodeSlug,
		Tags:     req.Tags,
	}

	err := h.postSvc.UpdateRootPost(cmd)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, post)
}

// Delete delete post
func (h *PostHandler) Delete(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	h.postSvc.Delete(gDto.ID)
	h.Success(ctx, nil)
}

// Undelete undelete post
func (h *PostHandler) Undelete(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	h.postSvc.Undelete(gDto.ID)
	h.Success(ctx, nil)
}

// List list posts
func (h *PostHandler) List(ctx *gin.Context) {
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

	list, paging := h.postSvc.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, post := range list {
		result := util.StructToMap(post, "content")
		result["user"] = render.ToDefaultUser(post.UserId)
		result["node"] = h.nodeSvc.Get(post.NodeId)
		result["tags"] = render.ToTags(h.postSvc.GetPostTags(post.ID))
		mr := markdown.NewMd().Run(post.Content)
		result["summary"] = mr.SummaryText

		results = append(results, result)
	}

	h.Success(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}

// Recommend 推荐
func (h *PostHandler) Recommend(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	err := h.postSvc.SetRecommend(gDto.ID, true)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, nil)
}

// Unrecommend 取消推荐
func (h *PostHandler) Unrecommend(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	err := h.postSvc.SetRecommend(gDto.ID, false)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, nil)
}
