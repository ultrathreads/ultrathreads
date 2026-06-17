package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/render"
	"ultrathreads/bus/event"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
)

type PostController struct {
	BaseController
}

// Show 话题详情
func (c *PostController) Show(ctx *gin.Context) {
	var req dto.SlugRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	post := service.Srv.Post.GetBySlug(req.Slug)
	if post == nil || post.Status != model.StatusOk {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}
	renderPost := render.ToPost(post)
	user := service.Srv.User.Get(post.UserId)
	renderPost.User = render.ToUser(user)
	log.Debug("renderPost.User=%v", renderPost.User)

	c.Success(ctx, renderPost)
}

// ListThreads 帖子列表（含扁平化回帖）
func (c *PostController) ListThreads(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	pageSize := util.FormIntDefault(ctx, "pageSize", 20)
	nodeSlug := util.ParamStringDefault(ctx, "slug", "")

	posts, paging := service.Srv.Post.GetNodeThreadsFull(page, pageSize, nodeSlug)

	var lastReadAtMap map[string]int64
	if nodeSlug != "" {
		lastReadAtMap = c.GetLastReadStates(ctx, nodeSlug)
	} else {
		lastReadAtMap = c.GetLastReadStates(ctx)
	}

	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts)

	rsp := model.PostListWithIncluded{
	    Data:     results,
	    Meta:     *paging,
	    Context:  model.Context{LastReadAtMap: lastReadAtMap},
	    Included: model.PostIncluded{
	        Users: incUsers,
	        Nodes: incNodes,
	        Tags:  incTags,
	    },
	}

	c.SuccessWithIncluded(ctx, rsp)
}

// ListTagThreads 标签帖子列表
func (c *PostController) ListTagThreads(ctx *gin.Context) {
	tagSlug := util.ParamStringDefault(ctx, "slug", "")
	page := util.FormIntDefault(ctx, "page", 1)

	posts, paging := service.Srv.Post.GetTagThreadsFull(tagSlug, page)

	lastReadAtMap := c.GetLastReadStates(ctx)

	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts)

	rsp := model.PostListWithIncluded{
	    Data:     results,
	    Meta:     *paging,
	    Context:  model.Context{LastReadAtMap: lastReadAtMap},
	    Included: model.PostIncluded{
	        Users: incUsers,
	        Nodes: incNodes,
	        Tags:  incTags,
	    },
	}

	c.SuccessWithIncluded(ctx, rsp)
}

// GetPostTree 帖子详情（含扁平化回帖）
func (c *PostController) GetPostTree(ctx *gin.Context) {
	var req dto.SlugRequest
	if !c.BindAndValidate(ctx, &req) {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	currentPost, posts, err := service.Srv.Post.GetPostTree(req.Slug)
	if err != nil {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	currentPostRender := render.ToPost(currentPost)
	user := service.Srv.User.Get(currentPost.UserId)
	currentPostRender.User = render.ToUser(user)

	//把render之后的currentPost压入extra，虽然有点怪异，但也算的上是巧思。
	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts)

	rsp := model.PostListWithIncluded{
	    Data:     results,
	    Included: model.PostIncluded{
	        Users: incUsers,
	        Nodes: incNodes,
	        Tags:  incTags,
	    },
	    Extra: currentPostRender,
	}
	c.SuccessWithIncluded(ctx, rsp)
}


// GetPostFlat 帖子详情（含扁平化回帖）
func (c *PostController) GetPostFlat(ctx *gin.Context) {
	var req dto.SlugRequest
	if !c.BindAndValidate(ctx, &req) {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	posts, err := service.Srv.Post.GetPostsByThreadId(req.Slug)
	if err != nil {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts, render.WithContent(), render.WithViewCount())
	rsp := model.PostListWithIncluded{
	    Data:     results,
	    Included: model.PostIncluded{
	        Users: incUsers,
	        Nodes: incNodes,
	        Tags:  incTags,
	    },
	}
	c.SuccessWithIncluded(ctx, rsp)
}

func (c *PostController) GetUserPosts(ctx *gin.Context) {
	// 1. 获取并验证基础参数（如 user_id）
	var req dto.SlugRequest
	if !c.BindAndValidate(ctx, &req) {
		return 
	}

	// 2. 获取分页和类型参数
	page := util.FormIntDefault(ctx, "page", 1)
	postType := ctx.DefaultQuery("type", "root")

	// 3. 调用 Service 层获取数据
	posts, paging := service.Srv.Post.GetUserPosts(req.Slug, postType, page, 20)

	lastReadAtMap := c.GetLastReadStates(ctx)

	results, incUsers, incNodes, incTags := render.ToSimplePostsWithIncluded(posts)

	rsp := model.PostListWithIncluded{
	    Data:     results,
	    Meta:     *paging,
	    Context:  model.Context{LastReadAtMap: lastReadAtMap},
	    Included: model.PostIncluded{
	        Users: incUsers,
	        Nodes: incNodes,
	        Tags:  incTags,
	    },
	}

	c.SuccessWithIncluded(ctx, rsp)
}

// StoreRootPost 发表根帖
func (c *PostController) StoreRootPost(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var req dto.CreateRootPostRequest

	if !c.BindAndValidate(ctx, &req) {
		return // BindAndValidate 内部已写回错误响应
	}

	post, err := service.Srv.Post.CreateRootPost(user.ID, req)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	// IsRoot 恒为 true，无需运行时判断
	c.PublishEvent(ctx, event.PostCreated{
		UserID: user.ID,
		PostID: post.ID,
		IsRoot: true,
		Tags:   req.Tags,
	})

	c.RespondOK(ctx, render.ToSimplePost(post))
}

// Update 更新主贴
func (c *PostController) UpdateRootPost(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var req dto.UpdateRootPostRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	post := service.Srv.Post.GetBySlug(req.Slug)
	if post == nil || post.Status == model.StatusDeleted {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	if post.UserId != user.ID {
		c.Fail(ctx, util.NewErrorMsg("无权限"))
		return
	}

	err := service.Srv.Post.UpdateRootPost(req)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	
	c.PublishEvent(ctx, event.PostUpdated{
		UserID: user.ID,
		PostID: post.ID,
		Tags:   req.Tags,
		IsRoot: post.IsRoot(),
	})

	c.RespondOK(ctx, render.ToSimplePost(post))
}

// StoreReply 发表回复
func (c *PostController) StoreReply(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	var req dto.CreateReplyRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	post, err := service.Srv.Post.CreateReply(user.ID, req)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	// IsRoot 恒为 false
	c.PublishEvent(ctx, event.PostCreated{
		UserID: user.ID,
		PostID: post.ID,
		IsRoot: false,
	})

	c.Success(ctx, render.ToSimplePost(post))
}

func (c *PostController) UpdateReply(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var req dto.UpdateReplyRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	post := service.Srv.Post.GetBySlug(req.Slug)
	if post == nil || post.Status == model.StatusDeleted {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	if post.UserId != user.ID {
		c.Fail(ctx, util.NewErrorMsg("无权限"))
		return
	}

	err := service.Srv.Post.UpdateReply(req)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	
	c.PublishEvent(ctx, event.PostUpdated{
		UserID: user.ID,
		PostID: post.ID,
		IsRoot: post.IsRoot(),
	})

	c.Success(ctx, render.ToSimplePost(post))
}

// GetUserRecent 用户最近的帖子
func (c *PostController) GetUserRecent(ctx *gin.Context) {
	var req dto.SlugRequest
	id := hashid.Slug2Id[model.User](req.Slug)
	if c.BindAndValidate(ctx, &req) {
		posts := service.Srv.Post.Find(querybuilder.NewQueryBuilder().Where("user_id = ? and status = ?",
			id, model.StatusOk).Desc("id").Limit(10))
		c.Success(ctx, render.ToSimplePosts(posts))
	}
}

// Like 点赞
func (c *PostController) Like(ctx *gin.Context) {
    user := c.GetCurrentUser(ctx)
    var req dto.SlugRequest
    if !c.BindAndValidateUri(ctx, &req) {
        return
    }

    id := hashid.Slug2Id[model.Post](req.Slug)
    if err := service.PostLikeService.Like(user.ID, id); err != nil {
        c.Fail(ctx, util.FromError(err))
        return
    }

    c.Success(ctx, nil)
}

// Favorite 收藏话题
func (c *PostController) Favorite(ctx *gin.Context) {
    user := c.GetCurrentUser(ctx)
    var req dto.SlugRequest
    if !c.BindAndValidateUri(ctx, &req) {
        return
    }

    if err := service.FavoriteService.AddPostFavorite(user.ID, req.Slug); err != nil {
        c.Fail(ctx, util.FromError(err))
        return
    }

    c.Success(ctx, nil)
}

func (c *PostController) ViewPost(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var req dto.SlugRequest
    if !c.BindAndValidateUri(ctx, &req) {
        return
    }

	nodeSlug := util.QueryStringDefault(ctx, "nodeSlug","")

    c.PublishEvent(ctx, event.PostViewed{
        UserID:     user.ID,
        PostSlug:   req.Slug,
        NodeSlug:   nodeSlug,
        ViewedTime: util.NowTimestamp(),
    })

	c.Success(ctx, nil)
}
