package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"ultrathreads/cache"
	"ultrathreads/model"
)

// GetCurrent 获取当前登录用户
func GetCurrent(ctx *gin.Context) *model.User {
	userTmp, _ := ctx.Get(viper.GetString("jwt.identity_key"))
	if userTmp == nil {
		return nil
	}

	user := cache.UserCache.Get(userTmp.(model.UserClaims).ID)
	if user == nil || user.Status != model.StatusOk {
		return nil
	}
	return user
}
