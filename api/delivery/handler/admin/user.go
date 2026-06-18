package admin

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
	"ultrathreads/dto"
	"ultrathreads/delivery/handler/base"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// UserHandler user controller
type UserHandler struct {
	base.BaseHandler
	userSvc   service.UserService
	userCache cache.UserCacheInterface
}

func NewUserHandler(userSvc service.UserService, userCache cache.UserCacheInterface) *UserHandler {
	return &UserHandler{userSvc: userSvc, userCache: userCache}
}

// Show show user
func (h *UserHandler) Show(ctx *gin.Context) {
	var req dto.SlugRequest
	if h.BindAndValidate(ctx, &req) {
		user := h.userSvc.GetBySlug(req.Slug)
		if user == nil {
			h.Fail(ctx, util.NewErrorMsg("User not found"))
			return
		}
		h.Success(ctx, h.buildUserItem(user))
	}
}

// Store 创建用户
func (h *UserHandler) Store(ctx *gin.Context) {
	h.Success(ctx, nil)
}

// Update 更新用户信息
func (h *UserHandler) Update(ctx *gin.Context) {
	var req dto.UpdateUserRequest
	if !h.BindAndValidate(ctx, &req) {
		return
	}
	user := h.userSvc.GetBySlug(req.Slug)
	if user == nil {
		h.Fail(ctx, util.NewErrorMsg("User not found"))
		return
	}

	err := h.userSvc.Update(req)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, user)
}

// Delete delete user
func (h *UserHandler) Delete(ctx *gin.Context) {
	var req dto.SlugRequest
	if !h.BindAndValidate(ctx, &req) {
		return
	}
	h.userSvc.Delete(req.Slug)
	h.Success(ctx, nil)
}

// List list users
func (h *UserHandler) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	id := ctx.Request.FormValue("id")
	nickname := ctx.Request.FormValue("nickname")
	username := ctx.Request.FormValue("username")

	conditions := querybuilder.NewQueryBuilder()
	if len(id) > 0 {
		conditions.Eq("id", id)
	}
	if len(username) > 0 {
		conditions.Eq("username", username)
	}
	if len(nickname) > 0 {
		conditions.Like("nickname", nickname)
	}
	list, paging := h.userSvc.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, user := range list {
		results = append(results, h.buildUserItem(&user))
	}

	h.Success(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}

func (h *UserHandler) buildUserItem(user *model.User) map[string]interface{} {
	score := h.userCache.GetScore(user.ID)

	result := make(map[string]interface{})
	result["id"] = user.ID
	result["status"] = user.Status
	result["level"] = user.Level
	result["username"] = user.Username.String
	result["nickname"] = user.Nickname
	result["avatar"] = user.Avatar
	result["email"] = user.Email.String
	result["score"] = score
	result["createTime"] = user.CreateTime
	result["updateTime"] = user.UpdateTime

	return result
}
