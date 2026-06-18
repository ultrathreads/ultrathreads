package app

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/delivery/handler/base"
	"ultrathreads/service"
	"ultrathreads/util"
)

type FavoriteHandler struct {
	base.BaseHandler
	favoriteSvc service.FavoriteService
}

func NewFavoriteHandler(favoriteSvc service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteSvc: favoriteSvc}
}

// GetFavorited 是否收藏
func (h *FavoriteHandler) GetFavorited(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	entityType := util.FormStringDefault(ctx, "entityType", "")
	entityID := util.FormInt64Default(ctx, "entityId", 0)

	data := map[string]interface{}{}
	if user == nil || len(entityType) == 0 || entityID <= 0 {
		data["favorited"] = false
	} else {
		tmp := h.favoriteSvc.GetBy(user.ID, entityType, entityID)
		data["favorited"] = tmp != nil
	}
	h.Success(ctx, data)
}

// Delete 取消收藏
func (h *FavoriteHandler) Delete(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	if user == nil {
		h.Fail(ctx, util.ErrorNotLogin)
		return
	}

	entityType := util.FormStringDefault(ctx, "entityType", "")
	entityID := util.FormInt64Default(ctx, "entityId", 0)

	tmp := h.favoriteSvc.GetBy(user.ID, entityType, entityID)
	if tmp != nil {
		h.favoriteSvc.Delete(tmp.ID)
	}
	h.Success(ctx, nil)
}
