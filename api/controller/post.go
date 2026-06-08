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
	page := form.FormValueIntDefault(ctx, "page", 1)

	posts, paging := service.PostService.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = converter.ToSimplePosts(posts)
	data["page"] = paging
	c.Success(ctx, data)
}

// ListWithReplies 帖子列表（含扁平化回帖）
func (c *PostController) ListWithReplies(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	nodeId := form.FormValueIntDefault(ctx, "nodeId", 0)

	posts, paging := service.PostService.ListThreadsWithReplies(page, limit, nodeId)

	data := map[string]interface{}{
		"results": converter.ToSimplePosts(posts),
		"page":    paging,
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
	page := form.FormValueIntDefault(ctx, "page", 1)

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
	page := form.FormValueIntDefault(ctx, "page", 1)

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
	page := form.FormValueIntDefault(ctx, "page", 1)

	posts, paging := service.PostService.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).
		Eq("comment_count", 0).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = converter.ToSimplePosts(posts)
	data["page"] = paging
	c.Success(ctx, data)
}

// 节点帖子列表
func (c *PostController) GetNodePosts(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	nodeId := form.FormValueInt64Default(ctx, "nodeId", 0)

	posts, paging := service.PostService.List(querybuilder.NewQueryBuilder().
		Eq("node_id", nodeId).
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = converter.ToSimplePosts(posts)
	data["page"] = paging

	c.Success(ctx, data)
}

// 标签帖子列表
func (c *PostController) GetTagPosts(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	tagId, err := form.FormValueInt64(ctx, "tagId")
	if err != nil {
		c.Fail(ctx, util.ErrorTagNotFound)
		return
	}
	posts, paging := service.PostService.GetTagPosts(tagId, page)

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

// GetUserPosts 用户的帖子
func (c *PostController) GetUserPosts(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		posts, paging := service.PostService.List(querybuilder.NewQueryBuilder().
			Eq("user_id", gDto.ID).
			Eq("status", model.StatusOk).
			Page(page, 20).Desc("id"))

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
