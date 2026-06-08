package router

import (
	"net/http"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"ultrathreads/middleware"
)

var (
	jwtAuth  *jwt.GinJWTMiddleware
	jwtOAuth *jwt.GinJWTMiddleware
)

// Setup 初始化路由引擎
func Setup(e *gin.Engine) {
	e.Use(gin.Recovery())
	e.Use(middleware.Cors())

	e.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "UltraThreads API\n")
	})

	// 初始化 JWT 中间件（供 api.go 和 admin_api.go 使用）
	jwtAuth = middleware.JwtAuth(middleware.LoginStandard)
	jwtOAuth = middleware.JwtAuth(middleware.LoginOAuth)

	// 注册各模块路由
	setupApp(e)
	setupAdmin(e)
}