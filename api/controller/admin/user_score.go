package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"ultrathreads/render"
	"ultrathreads/controller"
	"ultrathreads/dto"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// UserScoreController user score controller
type UserScoreController struct {
	controller.BaseController
	userScoreSvc service.UserScoreServicer
}

func NewUserScoreController(userScoreSvc service.UserScoreServicer) *UserScoreController {
	return &UserScoreController{userScoreSvc: userScoreSvc}
}

// Show 显示积分
func (c *UserScoreController) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		userScore := c.userScoreSvc.Get(gDto.ID)
		if userScore == nil {
			c.Fail(ctx, util.NewErrorMsg("User score not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, userScore)
	}
}

// List 显示积分列表
func (c *UserScoreController) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	userId := ctx.Request.FormValue("userId")

	conditions := querybuilder.NewQueryBuilder()
	if len(userId) > 0 {
		conditions.Eq("user_id", userId)
	}

	list, paging := c.userScoreSvc.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, userScore := range list {
		item := util.StructToMap(userScore)
		item["user"] = render.ToDefaultUser(userScore.UserId)
		results = append(results, item)
	}

	c.Success(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}