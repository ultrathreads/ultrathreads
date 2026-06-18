package middleware

import (
	"github.com/gin-gonic/gin"
)

func CurrentUser(ctx *gin.Context) {
	ctx.Set("CurrentUser", GetCurrent(ctx))
	ctx.Next()
}
