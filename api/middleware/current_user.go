package middleware

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/service"
)

func CurrentUser(ctx *gin.Context) {
	ctx.Set("CurrentUser", service.UserService.GetCurrent(ctx))
	ctx.Next()
}
