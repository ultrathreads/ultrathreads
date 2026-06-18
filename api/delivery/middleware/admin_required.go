package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ultrathreads/service"
	"ultrathreads/util"
)

// AdminRequired 校验用户是否具有超级管理员权限或后台管理准入资格
func AdminRequired(rbacSvc service.RbacService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := GetCurrent(ctx)
		if user == nil {
			err := util.ErrorNotLogin
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    err.Code,
				"message": err.Message,
			})
			return
		}

		if !rbacSvc.CanAccessAdminPanel(user.ID) {
			err := util.ErrorPermissionDenied
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    err.Code,
				"message": err.Message,
			})
			return
		}

		ctx.Next()
	}
}