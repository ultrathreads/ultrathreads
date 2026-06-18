package app

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/bus/event"
	"ultrathreads/delivery/handler/base"
	"ultrathreads/domain"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
)

type PostHandler struct {
	base.BaseHandler
	postSvc     service.PostService
	userSvc     service.UserService
	postLikeSvc service.PostLikeService
	favoriteSvc service.FavoriteService
}

func NewPostHandler(postSvc service.PostService, userSvc service.UserService, postLikeSvc service.PostLikeService, favoriteSvc service.FavoriteService) *PostHandler {
	return &PostHandler{
		postSvc:     postSvc,
		userSvc:     userSvc,
		postLikeSvc: postLikeSvc,
		favoriteSvc: favoriteSvc,
	}
}

// Show 话题详情
func (h *PostHandler) Show(ctx *gin.Context) {
	var req dto.SlugRequest
	if !h.BindAndValidate(ctx, &req) {
		return
	}

	post := h.postSvc.GetBySlug(req.Slug)
	if post == nil || post.Status != model.StatusOk {
		h.Fail(ctx, util.ErrorPostNotFound)
		return
	}
	renderPost := render.ToPost(post)
	user := h.userSvc.Get(post.UserId)
	renderPost.User = render.ToUser(user)
	log.Debug("renderPost.User=%v", renderPost.User)

	h.Success(ctx, renderPost)
}

// ListThreads 帖子列表（含扁平化回帖）
func (h *PostHandler) ListThreads(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	pageSize := util.FormIntDefault(ctx, "pageSize", 20)
	nodeSlug := util.ParamStringDefault(ctx, "slug", "")

	posts, paging := h.postSvc.GetNodeThreadsFull(page, pageSize, nodeSlug)

	var lastReadAtMap map[string]int64
	if nodeSlug != "" {
		lastReadAtMap = h.GetLastReadStates(ctx, nodeSlug)
	} else {
		lastReadAtMap = h.GetLastReadStates(ctx)
	}

	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts)

	rsp := model.PostListWithIncluded{
		Data:    results,
		Meta:    *paging,
		Context: model.Context{LastReadAtMap: lastReadAtMap},
		Included: model.PostIncluded{
			Users: incUsers,
			Nodes: incNodes,
			Tags:  incTags,
		},
	}

	h.SuccessWithIncluded(ctx, rsp)
}

// ListTagThreads 标签帖子列表
func (h *PostHandler) ListTagThreads(ctx *gin.Context) {
	tagSlug := util.ParamStringDefault(ctx, "slug", "")
	page := util.FormIntDefault(ctx, "page", 1)

	posts, paging := h.postSvc.GetTagThreadsFull(tagSlug, page)

	lastReadAtMap := h.GetLastReadStates(ctx)

	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts)

	rsp := model.PostListWithIncluded{
		Data:    results,
		Meta:    *paging,
		Context: model.Context{LastReadAtMap: lastReadAtMap},
		Included: model.PostIncluded{
			Users: incUsers,
			Nodes: incNodes,
			Tags:  incTags,
		},
	}

	h.SuccessWithIncluded(ctx, rsp)
}

// GetPostTree 帖子详情（含扁平化回帖）
func (h *PostHandler) GetPostTree(ctx *gin.Context) {
	var req dto.SlugRequest
	if !h.BindAndValidate(ctx, &req) {
		h.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	currentPost, posts, err := h.postSvc.GetPostTree(req.Slug)
	if err != nil {
		h.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	currentPostRender := render.ToPost(currentPost)
	user := h.userSvc.Get(currentPost.UserId)
	currentPostRender.User = render.ToUser(user)

	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts)

	rsp := model.PostListWithIncluded{
		Data: results,
		Included: model.PostIncluded{
			Users: incUsers,
			Nodes: incNodes,
			Tags:  incTags,
		},
		Extra: currentPostRender,
	}
	h.SuccessWithIncluded(ctx, rsp)
}

// GetPostFlat 帖子详情（含扁平化回帖）
func (h *PostHandler) GetPostFlat(ctx *gin.Context) {
	var req dto.SlugRequest
	if !h.BindAndValidate(ctx, &req) {
		h.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	posts, err := h.postSvc.GetPostsByThreadId(req.Slug)
	if err != nil {
		h.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts, render.WithContent(), render.WithViewCount())
	rsp := model.PostListWithIncluded{
		Data: results,
		Included: model.PostIncluded{
			Users: incUsers,
			Nodes: incNodes,
			Tags:  incTags,
		},
	}
	h.SuccessWithIncluded(ctx, rsp)
}

func (h *PostHandler) GetUserPosts(ctx *gin.Context) {
	var req dto.SlugRequest
	if !h.BindAndValidate(ctx, &req) {
		return
	}

	page := util.FormIntDefault(ctx, "page", 1)
	postType := ctx.DefaultQuery("type", "root")

	posts, paging := h.postSvc.GetUserPosts(req.Slug, postType, page, 20)

	lastReadAtMap := h.GetLastReadStates(ctx)

	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts)

	rsp := model.PostListWithIncluded{
		Data:    results,
		Meta:    *paging,
		Context: model.Context{LastReadAtMap: lastReadAtMap},
		Included: model.PostIncluded{
			Users: incUsers,
			Nodes: incNodes,
			Tags:  incTags,
		},
	}

	h.SuccessWithIncluded(ctx, rsp)
}

// StoreRootPost 发表根帖
func (h *PostHandler) StoreRootPost(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var req dto.CreateRootPostRequest

	if !h.BindAndValidate(ctx, &req) {
		return
	}

	cmd := domain.CreatePostCommand{
		NodeSlug:  req.NodeSlug,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		ImageList: req.ImageList,
	}

	post, err := h.postSvc.CreateRootPost(user.ID, cmd)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.PublishEvent(ctx, event.PostCreated{
		UserID: user.ID,
		PostID: post.ID,
		IsRoot: true,
		Tags:   req.Tags,
	})

	h.RespondOK(ctx, render.ToSimplePost(post))
}

// UpdateRootPost 更新主贴
func (h *PostHandler) UpdateRootPost(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var req dto.UpdateRootPostRequest
	if !h.BindAndValidate(ctx, &req) {
		return
	}

	post := h.postSvc.GetBySlug(req.Slug)
	if post == nil || post.Status == model.StatusDeleted {
		h.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	if post.UserId != user.ID {
		h.Fail(ctx, util.NewErrorMsg("无权操作"))
		return
	}

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

	h.PublishEvent(ctx, event.PostUpdated{
		UserID: user.ID,
		PostID: post.ID,
		Tags:   req.Tags,
		IsRoot: post.IsRoot(),
	})

	h.RespondOK(ctx, render.ToSimplePost(post))
}

// StoreReply 发表回复
func (h *PostHandler) StoreReply(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)

	var req dto.CreateReplyRequest
	if !h.BindAndValidate(ctx, &req) {
		return
	}

	cmd := domain.CreateReplyCommand{
		Slug:       req.Slug,
		Content:    req.Content,
		ImageList:  req.ImageList,
		ParentSlug: req.ParentSlug,
	}

	post, err := h.postSvc.CreateReply(user.ID, cmd)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.PublishEvent(ctx, event.PostCreated{
		UserID: user.ID,
		PostID: post.ID,
		IsRoot: false,
	})

	h.Success(ctx, render.ToSimplePost(post))
}

func (h *PostHandler) UpdateReply(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var req dto.UpdateReplyRequest
	if !h.BindAndValidate(ctx, &req) {
		return
	}

	post := h.postSvc.GetBySlug(req.Slug)
	if post == nil || post.Status == model.StatusDeleted {
		h.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	if post.UserId != user.ID {
		h.Fail(ctx, util.NewErrorMsg("无权操作"))
		return
	}

	cmd := domain.UpdateReplyCommand{
		Slug:    req.Slug,
		Content: req.Content,
	}

	err := h.postSvc.UpdateReply(cmd)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.PublishEvent(ctx, event.PostUpdated{
		UserID: user.ID,
		PostID: post.ID,
		IsRoot: post.IsRoot(),
	})

	h.Success(ctx, render.ToSimplePost(post))
}

// GetUserRecent 用户最近的帖子
func (h *PostHandler) GetUserRecent(ctx *gin.Context) {
	var req dto.SlugRequest
	id := hashid.Slug2Id[model.User](req.Slug)
	if h.BindAndValidate(ctx, &req) {
		posts := h.postSvc.Find(querybuilder.NewQueryBuilder().Where("user_id = ? and status = ?",
			id, model.StatusOk).Desc("id").Limit(10))
		h.Success(ctx, render.ToSimplePosts(posts))
	}
}

// Like 点赞
func (h *PostHandler) Like(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var req dto.SlugRequest
	if !h.BindAndValidateUri(ctx, &req) {
		return
	}

	id := hashid.Slug2Id[model.Post](req.Slug)
	if err := h.postLikeSvc.Like(user.ID, id); err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.Success(ctx, nil)
}

// Favorite 收藏话题
func (h *PostHandler) Favorite(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var req dto.SlugRequest
	if !h.BindAndValidateUri(ctx, &req) {
		return
	}

	if err := h.favoriteSvc.AddPostFavorite(user.ID, req.Slug); err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.Success(ctx, nil)
}

func (h *PostHandler) ViewPost(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var req dto.SlugRequest
	if !h.BindAndValidateUri(ctx, &req) {
		return
	}

	nodeSlug := util.QueryStringDefault(ctx, "nodeSlug", "")

	h.PublishEvent(ctx, event.PostViewed{
		UserID:     user.ID,
		PostSlug:   req.Slug,
		NodeSlug:   nodeSlug,
		ViewedTime: util.NowTimestamp(),
	})

	h.Success(ctx, nil)
}
