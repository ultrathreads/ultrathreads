package admin

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
	"ultrathreads/controller"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// UserController user controller
type UserController struct {
	controller.BaseController
	userSvc service.UserServicer
}

func NewUserController(userSvc service.UserServicer) *UserController {
	return &UserController{userSvc: userSvc}
}

// Show show user
func (c *UserController) Show(ctx *gin.Context) {
	var req dto.SlugRequest
	if c.BindAndValidate(ctx, &req) {
		user := c.userSvc.GetBySlug(req.Slug)
		if user == nil {
			c.Fail(ctx, util.NewErrorMsg("User not found"))
			return
		}
		c.Success(ctx, c.buildUserItem(user))
	}
}

// Store 创建用户
func (c *UserController) Store(ctx *gin.Context) {
	c.Success(ctx, nil)
}

// Update 更新用户信息
func (c *UserController) Update(ctx *gin.Context) {
	var req dto.UpdateUserRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	user := c.userSvc.GetBySlug(req.Slug)
	if user == nil {
		c.Fail(ctx, util.NewErrorMsg("User not found"))
		return
	}

	err := c.userSvc.Update(req)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, user)
}

// Delete delete user
func (c *UserController) Delete(ctx *gin.Context) {
	var req dto.SlugRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	c.userSvc.Delete(req.Slug)
	c.Success(ctx, nil)
}

// List list users
func (c *UserController) List(ctx *gin.Context) {
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
	list, paging := c.userSvc.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, user := range list {
		results = append(results, c.buildUserItem(&user))
	}

	c.Success(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}

func (c *UserController) buildUserItem(user *model.User) map[string]interface{} {
	score := cache.UserCache.GetScore(user.ID)

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