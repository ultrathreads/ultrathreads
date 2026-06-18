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

// UserScoreLogHandler user score controller
type UserScoreLogHandler struct {
	base.BaseHandler
	userScoreLogSvc service.UserScoreLogServicer
}

func NewUserScoreLogHandler(userScoreLogSvc service.UserScoreLogServicer) *UserScoreLogHandler {
	return &UserScoreLogHandler{userScoreLogSvc: userScoreLogSvc}
}

// Show 显示积分纪录
func (h *UserScoreLogHandler) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		userScoreLog := h.userScoreLogSvc.Get(gDto.ID)
		if userScoreLog == nil {
			h.Fail(ctx, util.NewErrorMsg("User score log not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		h.Success(ctx, userScoreLog)
	}
}

// List 显示积分列表
func (h *UserScoreLogHandler) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	userId := ctx.Request.FormValue("userId")
	sourceType := ctx.Request.FormValue("sourceType")
	sourceId := ctx.Request.FormValue("sourceId")
	ltype := ctx.Request.FormValue("type")

	conditions := querybuilder.NewQueryBuilder()
	if len(userId) > 0 {
		conditions.Eq("user_id", userId)
	}
	if len(sourceType) > 0 {
		conditions.Eq("source_type", sourceType)
	}
	if len(sourceId) > 0 {
		conditions.Eq("source_id", sourceId)
	}
	if len(ltype) > 0 {
		conditions.Eq("type", ltype)
	}

	list, paging := h.userScoreLogSvc.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, userScoreLog := range list {
		item := util.StructToMap(userScoreLog)
		item["user"] = render.ToDefaultUser(userScoreLog.UserId)
		results = append(results, item)
	}

	h.Success(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}
