package app

import (
	"strings"

	"github.com/gin-gonic/gin"

	"ultrathreads/dto"
	"ultrathreads/handler/base"
	"ultrathreads/model"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

type UserHandler struct {
	base.BaseHandler
	userSvc         service.UserServicer
	postSvc         service.PostServicer
	userScoreSvc    service.UserScoreServicer
	userScoreLogSvc service.UserScoreLogServicer
	notificationSvc service.NotificationServicer
	favoriteSvc     service.FavoriteServicer
	articleSvc      service.ArticleServicer
	userWatchSvc    service.UserWatchServicer
	rbacSvc         service.RbacServicer
}

func NewUserHandler(
	userSvc service.UserServicer,
	postSvc service.PostServicer,
	userScoreSvc service.UserScoreServicer,
	userScoreLogSvc service.UserScoreLogServicer,
	notificationSvc service.NotificationServicer,
	favoriteSvc service.FavoriteServicer,
	articleSvc service.ArticleServicer,
	userWatchSvc service.UserWatchServicer,
	rbacSvc service.RbacServicer,
) *UserHandler {
	return &UserHandler{
		userSvc:         userSvc,
		postSvc:         postSvc,
		userScoreSvc:    userScoreSvc,
		userScoreLogSvc: userScoreLogSvc,
		notificationSvc: notificationSvc,
		favoriteSvc:     favoriteSvc,
		articleSvc:      articleSvc,
		userWatchSvc:    userWatchSvc,
		rbacSvc:         rbacSvc,
	}
}

// GetCurrent get current user
func (h *UserHandler) GetCurrent(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)

	userInfo := render.ToUser(user)
	userInfo.Permissions = h.rbacSvc.GetUserPermissions(user.ID)
	userInfo.Roles = h.rbacSvc.GetUserRoles(user.ID)

	h.Success(ctx, userInfo)
}

// 用户详情
func (h *UserHandler) Show(ctx *gin.Context) {
	var req dto.SlugRequest
	if h.BindAndValidate(ctx, &req) {
		user := h.userSvc.GetBySlug(req.Slug)
		if user != nil && user.Status != model.StatusDeleted {
			h.Success(ctx, render.ToUser(user))
		} else {
			h.Fail(ctx, util.NewErrorMsg("用户不存在"))
		}
	}
}

// Update 用户资料编辑
func (h *UserHandler) Update(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	if user == nil {
		h.Fail(ctx, util.ErrorNotLogin)
		return
	}

	var req dto.UpdateUserRequest
	if h.BindAndValidate(ctx, &req) {
		if len(req.Website) > 0 && util.IsValidateUrl(req.Website) != nil {
			h.Fail(ctx, util.NewErrorMsg("个人主页地址错误"))
			return
		}
		err := h.userSvc.Updates(user.ID, map[string]interface{}{
			"nickname":    req.Nickname,
			"avatar":      req.Avatar,
			"website":     req.Website,
			"description": req.Description,
		})
		if err != nil {
			h.Fail(ctx, util.FromError(err))
			return
		}
		h.Success(ctx, nil)
	}
}

// GetScoreRank 积分排行
func (h *UserHandler) GetScoreRank(ctx *gin.Context) {
	userScores := h.userScoreSvc.Find(querybuilder.NewQueryBuilder().Desc("score").Limit(10))
	var results []*model.UserInfo
	for _, userScore := range userScores {
		results = append(results, render.ToDefaultUser(userScore.UserId))
	}
	h.Success(ctx, results)
}

// GetScorelogs 用户积分记录
func (h *UserHandler) GetScorelogs(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	user := h.GetCurrentUser(ctx)

	logs, paging := h.userScoreLogSvc.List(querybuilder.NewQueryBuilder().
		Eq("user_id", user.ID).
		Page(page, 20).Desc("id"))

	h.Success(ctx, gin.H{
		"results": logs,
		"paging":  paging,
	})
}

// GetNotificationsRecent 获取最近3条未读消息
func (h *UserHandler) GetNotificationsRecent(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var count int64 = 0
	var notifications []model.Notification
	if user != nil {
		count = h.notificationSvc.GetUnReadCount(user.ID)
		notifications = h.notificationSvc.Find(querybuilder.NewQueryBuilder().Eq("user_id", user.ID).Eq("status", model.NotificationStatusUnread).Limit(3).Desc("id"))
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["notifications"] = render.ToNotifications(notifications)
	h.Success(ctx, data)
}

// GetNotifications 用户通知
func (h *UserHandler) GetNotifications(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	page := util.FormIntDefault(ctx, "page", 1)

	messages, paging := h.notificationSvc.List(querybuilder.NewQueryBuilder().
		Eq("user_id", user.ID).
		Page(page, 20).Desc("id"))

	// 全部标记为已读
	h.notificationSvc.MarkRead(user.ID)

	h.Success(ctx, gin.H{
		"results": render.ToNotifications(messages),
		"paging":  paging,
	})
}

// GetFavorites get favorites
func (h *UserHandler) GetFavorites(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	page := util.FormIntDefault(ctx, "page", 1)

	// 1. 查询收藏列表
	qb := querybuilder.NewQueryBuilder().
		Eq("user_id", user.ID).
		Page(page, 20).
		Desc("id")
	favorites, paging := h.favoriteSvc.List(qb)

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

	// 3. 批量查询文章
	articles := h.articleSvc.GetArticleInIds(articleIDs)
	articleMap := make(map[int64]*model.Article, len(articles))
	for i := range articles {
		articleMap[articles[i].ID] = &articles[i]
	}

	// 4. 批量查询帖子
	rawPostMap := h.postSvc.GetPostInIds(postIDs)
	postMap := make(map[int64]*model.Post, len(rawPostMap))
	for id, pst := range rawPostMap {
		tmp := pst
		postMap[id] = &tmp
	}

	// 5. 提取 Author ID，批量查用户
	var userIDs []int64
	for _, art := range articleMap {
		userIDs = append(userIDs, art.UserId)
	}
	for _, pst := range postMap {
		userIDs = append(userIDs, pst.UserId)
	}

	userMap := make(map[int64]*model.User)
	if len(userIDs) > 0 {
		users := h.userSvc.Find(
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
	h.Success(ctx, gin.H{
		"results": render.ToFavorites(favorites, contexts),
		"page":    paging,
	})
}

// Watch 关注
func (h *UserHandler) Watch(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var req dto.SlugRequest
	if h.BindAndValidate(ctx, &req) {
		err := h.userWatchSvc.Watch(req.Slug, user.ID)
		if err != nil {
			h.Fail(ctx, util.FromError(err))
			return
		}
		h.Success(ctx, nil)
	}
}

// GetWatched 是否关注了
func (h *UserHandler) GetWatched(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)

	userID := util.FormInt64Default(ctx, "userId", 0)

	data := map[string]interface{}{}
	if user == nil || userID <= 0 {
		data["watched"] = false
	} else {
		tmp := h.userWatchSvc.GetBy(userID, user.ID)
		data["watched"] = tmp != nil
	}
	h.Success(ctx, data)
}

// WatchDelete 取消关注
func (h *UserHandler) WatchDelete(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)

	userID := util.FormInt64Default(ctx, "userId", 0)

	tmp := h.userWatchSvc.GetBy(userID, user.ID)
	if tmp != nil {
		h.userWatchSvc.Delete(tmp.ID)
	}
	h.Success(ctx, nil)
}

// UpdateAvatar 修改头像
func (h *UserHandler) UpdateAvatar(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	avatar := strings.TrimSpace(ctx.Request.FormValue("avatar"))

	err := h.userSvc.UpdateAvatar(user.ID, avatar)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, nil)
}

// SetUsername 设置用户名
func (h *UserHandler) SetUsername(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	username := strings.TrimSpace(ctx.Request.FormValue("username"))

	err := h.userSvc.SetUsername(user.ID, username)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, nil)
}

// SetEmail 设置邮箱
func (h *UserHandler) SetEmail(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	email := strings.TrimSpace(ctx.Request.FormValue("email"))

	err := h.userSvc.SetEmail(user.ID, email)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, nil)
}

// SetPassword 设置密码
func (h *UserHandler) SetPassword(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)

	var (
		password   = strings.TrimSpace(ctx.Request.FormValue("password"))
		rePassword = strings.TrimSpace(ctx.Request.FormValue("rePassword"))
	)

	err := h.userSvc.SetPassword(user.ID, password, rePassword)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, nil)
}

// ChangePassword 更改密码
func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var (
		oldPassword = ctx.Request.FormValue("oldPassword")
		password    = ctx.Request.FormValue("password")
		rePassword  = ctx.Request.FormValue("rePassword")
	)
	err := h.userSvc.UpdatePassword(user.ID, oldPassword, password, rePassword)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, nil)
}
