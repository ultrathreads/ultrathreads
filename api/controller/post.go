package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/converter"
	"ultrathreads/cache"
	"ultrathreads/bus/event"
	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	//"ultrathreads/util/log"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
)

type PostController struct {
	BaseController
}

// Show 话题详情
func (c *PostController) Show(ctx *gin.Context) {
	var gDto form.IdentifierDto
	if c.BindAndValidate(ctx, &gDto) {
		post := service.PostService.GetBySlug(gDto.Slug)
		if post == nil || post.Status != model.StatusOk {
			c.Fail(ctx, util.ErrorPostNotFound)
			return
		}
		c.Success(ctx, converter.ToPost(post))
	}
}

// List 帖子列表
func (c *PostController) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)

	posts, paging := service.PostService.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = converter.ToSimplePosts(posts)
	data["page"] = paging
	c.Success(ctx, data)
}

// ListThreads 帖子列表（含扁平化回帖）
func (c *PostController) ListThreads(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	nodeSlug := util.ParamStringDefault(ctx, "slug", "")

	posts, paging := service.PostService.GetNodeThreadsFull(page, limit, nodeSlug)

	var lastReadAtMap map[string]int64
	if nodeSlug != "" {
		lastReadAtMap = c.GetLastReadStates(ctx, nodeSlug)
	} else {
		lastReadAtMap = c.GetLastReadStates(ctx)
	}

	data := map[string]interface{}{
		"results": 		 converter.ToSimplePosts(posts),
		"page":    		 paging,
		"lastReadAtMap": lastReadAtMap,
	}
	c.Success(ctx, data)
}

// ListTagThreads 标签帖子列表
func (c *PostController) ListTagThreads(ctx *gin.Context) {
	tagSlug := util.ParamStringDefault(ctx, "slug", "")
	page := util.FormIntDefault(ctx, "page", 1)

	posts, paging := service.PostService.GetTagThreadsFull(tagSlug, page)

	lastReadAtMap := c.GetLastReadStates(ctx)

	data := map[string]interface{}{
		"results":       converter.ToSimplePosts(posts),
		"page":          paging,
		"lastReadAtMap": lastReadAtMap,
	}
	c.Success(ctx, data)
}

// GetPostWithThread 帖子详情（含扁平化回帖）
func (c *PostController) GetPostWithThread(ctx *gin.Context) {

	var gDto form.IdentifierDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	post, replies, err := service.PostService.GetPostWithThread(gDto.Slug)
	if err != nil {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	data := map[string]interface{}{
		"post":    converter.ToPost(post),
		"replies": converter.ToSimplePosts(replies),
	}
	c.Success(ctx, data)
}


// GetPostWithFlat 帖子详情（含扁平化回帖）
func (c *PostController) GetPostsFlat(ctx *gin.Context) {
	var gDto form.IdentifierDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	posts, err := service.PostService.GetPostsByThreadId(gDto.Slug)
	if err != nil {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	data := map[string]interface{}{
		"posts": converter.ToPosts(posts),
	}
	c.Success(ctx, data)
}

// StoreRootPost 发表根帖
func (c *PostController) StoreRootPost(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var postForm form.RootPostCreateForm

	if !c.BindAndValidate(ctx, &postForm) {
		return // BindAndValidate 内部已写回错误响应
	}

	postForm.UserSlug = hashid.Id2Slug[model.User](user.ID)

	post, err := service.PostService.CreateRootPost(postForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	// ✅ IsRoot 恒为 true，无需运行时判断
	c.PublishEvent(ctx, event.PostCreated{
		UserID: user.ID,
		PostID: post.ID,
		IsRoot: true,
		Tags:   postForm.Tags,
	})

	c.Success(ctx, converter.ToSimplePost(post))
}

// StoreReply 发表回复
func (c *PostController) StoreReply(ctx *gin.Context) {
	var gDto form.IdentifierDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	user := c.GetCurrentUser(ctx)

	var replyForm form.ReplyCreateForm
	if !c.BindAndValidate(ctx, &replyForm) {
		return
	}

	replyForm.UserSlug = hashid.Id2Slug[model.User](user.ID)
	replyForm.ParentSlug = gDto.Slug
	replyForm.Title = util.ExtractReplyTitle(replyForm.Content, 20) // 从内容提取前20字符

	post, err := service.PostService.CreateReply(replyForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	// ✅ IsRoot 恒为 false
	c.PublishEvent(ctx, event.PostCreated{
		UserID: user.ID,
		PostID: post.ID,
		IsRoot: false,
	})

	c.Success(ctx, converter.ToSimplePost(post))
}

// Update 更新话题
func (c *PostController) Update(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var postForm form.PostUpdateForm
	if !c.BindAndValidate(ctx, &postForm) {
		return
	}

	post := service.PostService.GetBySlug(postForm.Slug)
	if post == nil || post.Status == model.StatusDeleted {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	if post.UserId != user.ID {
		c.Fail(ctx, util.NewErrorMsg("无权限"))
		return
	}

	err := service.PostService.Update(postForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	
	c.PublishEvent(ctx, event.PostUpdated{
		UserID: user.ID,
		PostID: post.ID,
		Tags:   postForm.Tags,
		IsRoot: post.IsRoot(),
	})

	c.Success(ctx, converter.ToSimplePost(post))
}

// GetRecentLikes 点赞用户
func (c *PostController) GetRecentLikes(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		postLikes := service.PostLikeService.Recent(gDto.ID, 10)
		var users []model.UserInfo
		for _, postLike := range postLikes {
			userInfo := converter.ToUserById(postLike.UserId)
			if userInfo != nil {
				users = append(users, *userInfo)
			}
		}
		c.Success(ctx, users)
	}
}

// 精华帖子
func (c *PostController) GetPostsExcellent(ctx *gin.Context) {
	posts := cache.PostCache.GetRecommendPosts()

	var odd, even []model.Post
	for i, post := range posts {
		if i%2 == 1 {
			odd = append(odd, post)
		} else {
			even = append(even, post)
		}
	}

	data := make(map[string]interface{})
	data["odd"] = converter.ToSimplePosts(odd)
	data["even"] = converter.ToSimplePosts(even)

	c.Success(ctx, data)
}

// 推荐帖子
func (c *PostController) GetPostsRecommend(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)

	posts, paging := service.PostService.List(querybuilder.NewQueryBuilder().
		Eq("recommend", true).
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = converter.ToSimplePosts(posts)
	data["page"] = paging
	c.Success(ctx, data)
}

// 最新发布帖子列表
func (c *PostController) GetPostsLast(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)

	posts, paging := service.PostService.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("id"))

	data := map[string]interface{}{}
	data["results"] = converter.ToSimplePosts(posts)
	data["page"] = paging
	c.Success(ctx, data)
}

// 无人问津帖子列表
func (c *PostController) GetPostsNoreply(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)

	posts, paging := service.PostService.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).
		Eq("comment_count", 0).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = converter.ToSimplePosts(posts)
	data["page"] = paging
	c.Success(ctx, data)
}

// GetUserRecent 用户最近的帖子
func (c *PostController) GetUserRecent(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		posts := service.PostService.Find(querybuilder.NewQueryBuilder().Where("user_id = ? and status = ?",
			gDto.ID, model.StatusOk).Desc("id").Limit(10))
		c.Success(ctx, converter.ToSimplePosts(posts))
	}
}

func (c *PostController) GetUserPosts(ctx *gin.Context) {
	// 1. 获取并验证基础参数（如 user_id）
	var gDto form.IdentifierDto
	if !c.BindAndValidate(ctx, &gDto) {
		return 
	}

	// 2. 获取分页和类型参数
	page := util.FormIntDefault(ctx, "page", 1)
	postType := ctx.DefaultQuery("type", "root")

	// 3. 调用 Service 层获取数据
	posts, paging := service.PostService.GetUserPosts(gDto.Slug, postType, page, 20)

	data := map[string]interface{}{}
	data["results"] = converter.ToSimplePosts(posts)
	data["page"] = paging
	//data["lastReadAtMap"] = c.GetLastReadStates(ctx)
	// 4. 格式化并返回结果
	c.Success(ctx, data)
}

// Like 点赞
func (c *PostController) Like(ctx *gin.Context) {
    user := c.GetCurrentUser(ctx)
    var gDto form.GeneralGetDto
    if !c.BindAndValidateUri(ctx, &gDto) {
        return
    }

    if err := service.PostLikeService.Like(user.ID, gDto.ID); err != nil {
        c.Fail(ctx, util.FromError(err))
        return
    }

    c.Success(ctx, nil)
}

// Favorite 收藏话题
func (c *PostController) Favorite(ctx *gin.Context) {
    user := c.GetCurrentUser(ctx)
    var gDto form.IdentifierDto
    if !c.BindAndValidateUri(ctx, &gDto) {
        return
    }

    if err := service.FavoriteService.AddPostFavorite(user.ID, gDto.Slug); err != nil {
        c.Fail(ctx, util.FromError(err))
        return
    }

    c.Success(ctx, nil)
}

func (c *PostController) ViewPost(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.IdentifierDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

    c.PublishEvent(ctx, event.PostViewed{
        UserID:     user.ID,
        PostSlug:   gDto.Slug,
        ViewedTime: util.NowTimestamp(),
    })

	c.Success(ctx, nil)
}
