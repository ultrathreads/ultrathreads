package router

import (
	"net/http"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"ultrathreads/bus"
	"ultrathreads/middleware"
	"ultrathreads/service"
)

var (
	jwtAuth  *jwt.GinJWTMiddleware
	jwtOAuth *jwt.GinJWTMiddleware
)

// Setup 初始化路由引擎
func Setup(e *gin.Engine, mgr *bus.Manager, srv *service.Services) {
	e.Use(gin.Recovery())
	e.Use(middleware.Cors())

	e.Use(func(c *gin.Context) {
		c.Set(bus.BusKey, mgr.Bus)
		c.Next()
	})

	e.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "UltraThreads API\n")
	})

	// 初始化 JWT 中间件（供 api.go 和 admin_api.go 使用）
	jwtAuth = middleware.JwtAuth(middleware.LoginStandard)
	jwtOAuth = middleware.JwtAuth(middleware.LoginOAuth)

	// 注册各模块路由
	setupApp(e, srv)
	setupAdmin(e)
}