package controller

import (
	"github.com/gin-gonic/gin"
	"strings"

	"ultrathreads/render"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

type UserController struct {
	BaseController
}

// GetCurrent get current user
func (c *UserController) GetCurrent(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	userInfo := render.ToUser(user)
	userInfo.Permissions = service.RbacService.GetUserPermissions(user.ID)
	userInfo.Roles = service.RbacService.GetUserRoles(user.ID)

	c.Success(ctx, userInfo)
}

// 用户详情
func (c *UserController) Show(ctx *gin.Context) {
	var req dto.SlugRequest
	if c.BindAndValidate(ctx, &req) {
		user := service.Srv.User.GetBySlug(req.Slug)
		if user != nil && user.Status != model.StatusDeleted {
			c.Success(ctx, render.ToUser(user))
		} else {
			c.Fail(ctx, util.NewErrorMsg("用户不存在"))
		}
	}
}

// Update 用户资料编辑
func (c *UserController) Update(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, util.ErrorNotLogin)
		return
	}

	var req dto.UpdateUserRequest
	if c.BindAndValidate(ctx, &req) {
		if len(req.Website) > 0 && util.IsValidateUrl(req.Website) != nil {
			c.Fail(ctx, util.NewErrorMsg("个人主页地址错误"))
			return
		}
		err := service.Srv.User.Updates(user.ID, map[string]interface{}{
			"nickname":    req.Nickname,
			"avatar":      req.Avatar,
			"website":     req.Website,
			"description": req.Description,
		})
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, nil)
	}
}

// GetScoreRank 积分排行
func (c *UserController) GetScoreRank(ctx *gin.Context) {
	userScores := service.UserScoreService.Find(querybuilder.NewQueryBuilder().Desc("score").Limit(10))
	var results []*model.UserInfo
	for _, userScore := range userScores {
		results = append(results, render.ToDefaultUser(userScore.UserId))
	}
	c.Success(ctx, results)
}

// GetScorelogs 用户积分记录
func (c *UserController) GetScorelogs(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	user := c.GetCurrentUser(ctx)

	logs, paging := service.UserScoreLogService.List(querybuilder.NewQueryBuilder().
		Eq("user_id", user.ID).
		Page(page, 20).Desc("id"))

	c.Success(ctx, gin.H{
		"results": logs,
		"paging":  paging,
	})
}

// GetNotificationsRecent 获取最近3条未读消息
func (c *UserController) GetNotificationsRecent(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var count int64 = 0
	var notifications []model.Notification
	if user != nil {
		count = service.NotificationService.GetUnReadCount(user.ID)
		notifications = service.NotificationService.Find(querybuilder.NewQueryBuilder().Eq("user_id", user.ID).Eq("status", model.NotificationStatusUnread).Limit(3).Desc("id"))
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["notifications"] = render.ToNotifications(notifications)
	c.Success(ctx, data)
}

// GetNotifications 用户通知
func (c *UserController) GetNotifications(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	page := util.FormIntDefault(ctx, "page", 1)

	messages, paging := service.NotificationService.List(querybuilder.NewQueryBuilder().
		Eq("user_id", user.ID).
		Page(page, 20).Desc("id"))

	// 全部标记为已读
	service.NotificationService.MarkRead(user.ID)

	c.Success(ctx, gin.H{
		"results": render.ToNotifications(messages),
		"paging":  paging,
	})
}

// GetFavorites get favorites
func (c *UserController) GetFavorites(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	page := util.FormIntDefault(ctx, "page", 1)

	// 1. 查询收藏列表
	qb := querybuilder.NewQueryBuilder().
		Eq("user_id", user.ID).
		Page(page, 20).
		Desc("id")
	favorites, paging := service.FavoriteService.List(qb)

	// 2. 收集需要预加载的实体 ID
	var articleIDs, postIDs []int64
	for _, fav := range favorites {
		switch fav.EntityType {
		case model.EntityTypeArticle:
			articleIDs = append(articleIDs, fav.EntityId)
		case model.EntityTypePost:
			postIDs = append(postIDs, fav.EntityId)
		}
	}

	// 3. 批量查询文章（返回切片，需手动转指针Map）
	articles := service.ArticleService.GetArticleInIds(articleIDs)
	articleMap := make(map[int64]*model.Article, len(articles))
	for i := range articles {
		articleMap[articles[i].ID] = &articles[i]
	}

	// 4. ✅ 批量查询帖子（返回值类型Map，需转换为指针Map）
	rawPostMap := service.Srv.Post.GetPostInIds(postIDs)
	postMap := make(map[int64]*model.Post, len(rawPostMap))
	for id, pst := range rawPostMap {
		// ⚠️ 关键：必须用临时变量接收值再取地址，不能在 range value 上直接 &pst
		tmp := pst 
		postMap[id] = &tmp
	}

	// 5. 提取 Author ID，使用现有 Find + QueryBuilder.In 批量查用户
	var userIDs []int64
	for _, art := range articleMap {
		userIDs = append(userIDs, art.UserId)
	}
	for _, pst := range postMap {
		userIDs = append(userIDs, pst.UserId)
	}

	userMap := make(map[int64]*model.User)
	if len(userIDs) > 0 {
		users := service.Srv.User.Find(
			querybuilder.NewQueryBuilder().In("id", userIDs),
		)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}
	}

	// 6. 组装 Context 切片
	contexts := make([]*render.FavoriteContext, len(favorites))
	for i, fav := range favorites {
		favCtx := &render.FavoriteContext{}
		switch fav.EntityType {
		case model.EntityTypeArticle:
			if art, ok := articleMap[fav.EntityId]; ok {
				favCtx.Article = art
				favCtx.User = userMap[art.UserId]
			}
		case model.EntityTypePost:
			if pst, ok := postMap[fav.EntityId]; ok {
				favCtx.Post = pst
				favCtx.User = userMap[pst.UserId]
			}
		}
		contexts[i] = favCtx
	}

	// 7. 纯渲染
	c.Success(ctx, gin.H{
		"results": render.ToFavorites(favorites, contexts),
		"page":    paging,
	})
}

// Watch 关注
func (c *UserController) Watch(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var req dto.SlugRequest
	if c.BindAndValidate(ctx, &req) {
		err := service.UserWatchService.Watch(req.Slug, user.ID)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, nil)
	}
}

// GetWatched 是否关注了
func (c *UserController) GetWatched(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	userID := util.FormInt64Default(ctx, "userId", 0)

	data := map[string]interface{}{}
	if user == nil || userID <= 0 {
		data["watched"] = false
	} else {
		tmp := service.UserWatchService.GetBy(userID, user.ID)
		data["watched"] = tmp != nil
	}
	c.Success(ctx, data)
}

// Delete 取消收藏
func (c *UserController) WatchDelete(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	userID := util.FormInt64Default(ctx, "userId", 0)

	tmp := service.UserWatchService.GetBy(userID, user.ID)
	if tmp != nil {
		service.UserWatchService.Delete(tmp.ID)
	}
	c.Success(ctx, nil)
}

// UpdateAvatar 修改头像
func (c *UserController) UpdateAvatar(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	avatar := strings.TrimSpace(ctx.Request.FormValue("avatar"))

	err := service.Srv.User.UpdateAvatar(user.ID, avatar)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// SetUsername 设置用户名
func (c *UserController) SetUsername(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	username := strings.TrimSpace(ctx.Request.FormValue("username"))

	err := service.Srv.User.SetUsername(user.ID, username)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// SetEmail 设置邮箱
func (c *UserController) SetEmail(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	email := strings.TrimSpace(ctx.Request.FormValue("email"))

	err := service.Srv.User.SetEmail(user.ID, email)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// SetPassword 设置密码
func (c *UserController) SetPassword(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	var (
		password   = strings.TrimSpace(ctx.Request.FormValue("password"))
		rePassword = strings.TrimSpace(ctx.Request.FormValue("rePassword"))
	)

	err := service.Srv.User.SetPassword(user.ID, password, rePassword)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// ChangePassword 更改密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var (
		oldPassword = ctx.Request.FormValue("oldPassword")
		password    = ctx.Request.FormValue("password")
		rePassword  = ctx.Request.FormValue("rePassword")
	)
	err := service.Srv.User.UpdatePassword(user.ID, oldPassword, password, rePassword)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}
