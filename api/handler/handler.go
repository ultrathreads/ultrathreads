package handler

import (
	"net/http"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"ultrathreads/bus"
	"ultrathreads/handler/admin"
	"ultrathreads/handler/app"
	"ultrathreads/middleware"
	"ultrathreads/service"
)

type Handler struct {
	services *service.Services
	mgr      *bus.Manager

	jwtAuth  *jwt.GinJWTMiddleware
	jwtOAuth *jwt.GinJWTMiddleware
}

func NewHandlers(services *service.Services, mgr *bus.Manager) *Handler {
	h := &Handler{
		services: services,
		mgr:      mgr,
	}

	h.jwtAuth = middleware.JwtAuth(middleware.LoginStandard, services.User, services.LoginSource)
	h.jwtOAuth = middleware.JwtAuth(middleware.LoginOAuth, services.User, services.LoginSource)

	return h
}

// Init 组装引擎、全局中间件，并注册所有路由
func (h *Handler) Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middleware.Cors())
	router.Use(func(c *gin.Context) {
		c.Set(bus.BusKey, h.mgr.Bus)
		c.Next()
	})

	router.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "UltraThreads API\n")
	})
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 创建 /api 路由组
	api := router.Group("/api")

	// 注册前台路由
	appHandler := app.NewHandler(h.services, h.jwtAuth, h.jwtOAuth)
	appHandler.Init(api)

	// 注册后台管理路由
	adminHandler := admin.NewHandler(h.services, h.jwtAuth)
	adminHandler.Init(api)

	return router
}
