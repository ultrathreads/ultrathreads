package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ultrathreads/dto"
	"ultrathreads/handler/base"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// UserScoreHandler user score controller
type UserScoreHandler struct {
	base.BaseHandler
	userScoreSvc service.UserScoreServicer
}

func NewUserScoreHandler(userScoreSvc service.UserScoreServicer) *UserScoreHandler {
	return &UserScoreHandler{userScoreSvc: userScoreSvc}
}

// Show 显示积分
func (h *UserScoreHandler) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		userScore := h.userScoreSvc.Get(gDto.ID)
		if userScore == nil {
			h.Fail(ctx, util.NewErrorMsg("User score not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		h.Success(ctx, userScore)
	}
}

// List 显示积分列表
func (h *UserScoreHandler) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	userId := ctx.Request.FormValue("userId")

	conditions := querybuilder.NewQueryBuilder()
	if len(userId) > 0 {
		conditions.Eq("user_id", userId)
	}

	list, paging := h.userScoreSvc.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, userScore := range list {
		item := util.StructToMap(userScore)
		item["user"] = render.ToDefaultUser(userScore.UserId)
		results = append(results, item)
	}

	h.Success(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}
