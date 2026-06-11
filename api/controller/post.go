package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/converter"
	"ultrathreads/cache"
	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

type PostController struct {
	BaseController
}

// Show 话题详情
func (c *PostController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		post := service.PostService.Get(gDto.ID)
		if post == nil || post.Status != model.StatusOk {
			c.Fail(ctx, util.ErrorPostNotFound)
			return
		}
		service.PostService.IncrViewCount(post.ID) // 增加浏览量
		c.Success(ctx, converter.ToPost(post))
	}
}

// List 帖子列表
func (c *PostController) List(ctx *gin.Context) {
	page := util.FormValueIntDefault(ctx, "page", 1)

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
	page := util.FormValueIntDefault(ctx, "page", 1)
	limit := util.FormValueIntDefault(ctx, "limit", 20)
	nodeId := util.FormValueIntDefault(ctx, "nodeId", 0)

	posts, paging := service.PostService.GetNodeThreadsFull(page, limit, nodeId)

	var lastReadAtMap map[int64]int64
	if nodeId > 0 {
	    lastReadAtMap = c.GetLastReadStates(ctx, int64(nodeId))
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
	page := util.FormValueIntDefault(ctx, "page", 1)
	tagId, err := util.FormValueInt64(ctx, "tagId")
	if err != nil {
		c.Fail(ctx, util.ErrorTagNotFound)
		return
	}
	posts, paging := service.PostService.GetTagThreadsFull(tagId, page)

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
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	post, replies, err := service.PostService.GetPostWithThread(gDto.ID)
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
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	posts, err := service.PostService.GetPostsByThreadId(gDto.ID)
	if err != nil {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	data := map[string]interface{}{
		"posts": converter.ToPosts(posts),
	}
	c.Success(ctx, data)
}

// Store 发表帖子
func (c *PostController) Store(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var postForm form.PostCreateForm
	if c.BindAndValidate(ctx, &postForm) {
		postForm.UserID = user.ID
		post, err := service.PostService.Create(postForm)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, converter.ToSimplePost(post))
	}
}

// Edit 为编辑话题准备数据
func (c *PostController) Edit(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		post := service.PostService.Get(gDto.ID)
		if post == nil || post.Status != model.StatusOk {
			c.Fail(ctx, util.NewErrorMsg("话题不存在或已被删除"))
			return
		}
		if post.UserId != user.ID {
			c.Fail(ctx, util.NewErrorMsg("无权限"))
			return
		}

		tags := service.PostService.GetPostTags(post.ID)
		var tagNames []string
		if len(tags) > 0 {
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}
		}

		c.Success(ctx, gin.H{
			"postId": post.ID,
			"nodeId":  post.NodeId,
			"title":   post.Title,
			"content": post.Content,
			"tags":    tagNames,
		})
	}
}

// Update 更新话题
func (c *PostController) Update(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	post := service.PostService.Get(gDto.ID)
	if post == nil || post.Status == model.StatusDeleted {
		c.Fail(ctx, util.ErrorPostNotFound)
		return
	}

	if post.UserId != user.ID {
		c.Fail(ctx, util.NewErrorMsg("无权限"))
		return
	}

	var postForm form.PostUpdateForm
	if c.BindAndValidate(ctx, &postForm) {
		postForm.ID = post.ID
		err := service.PostService.Update(postForm)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, converter.ToSimplePost(post))
	}
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
	page := util.FormValueIntDefault(ctx, "page", 1)

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
	page := util.FormValueIntDefault(ctx, "page", 1)

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
	page := util.FormValueIntDefault(ctx, "page", 1)

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

// GetUserPosts 用户的帖子（支持通过 type 参数区分 root 和 reply）
func (c *PostController) GetUserPosts(ctx *gin.Context) {
	// 1. 获取分页参数
	page := util.FormValueIntDefault(ctx, "page", 1)
	
	// 2. 获取并验证基础参数（如 user_id）
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return // 注意：BindAndValidate 失败时通常框架会自动返回错误，这里做防御性处理
	}

	// 3. 获取 type 参数（默认为 root，兼容旧逻辑）
	postType := ctx.DefaultQuery("type", "root")

	// 4. 构建基础查询条件
	qb := querybuilder.NewQueryBuilder().
		Eq("user_id", gDto.ID).
		Eq("status", model.StatusOk).
		Page(page, 20).
		Desc("id")

	// 5. 根据 type 动态追加过滤条件
	switch postType {
	case "reply":
		// 假设回帖的 parent_id 大于 0，或者不等于 0
		// 具体字段名和逻辑请根据你的数据库表结构进行调整
		qb.NotEq("parent_id", 0) 
	case "root":
		fallthrough // 默认情况，只查询根帖
	default:
		// 假设根帖的 parent_id 为 0
		qb.Eq("parent_id", 0) 
	}

	// 6. 执行查询并返回结果
	if c.BindAndValidate(ctx, &gDto) { // 保持你原有的执行逻辑
		posts, paging := service.PostService.List(qb)

		c.Success(ctx, gin.H{
			"results": converter.ToSimplePosts(posts),
			"page":    paging,
		})
	}
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
    var gDto form.GeneralGetDto
    if !c.BindAndValidateUri(ctx, &gDto) {
        return
    }

    if err := service.FavoriteService.AddPostFavorite(user.ID, gDto.ID); err != nil {
        c.Fail(ctx, util.FromError(err))
        return
    }

    c.Success(ctx, nil)
}
