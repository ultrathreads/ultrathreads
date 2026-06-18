package app

import (
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
	"ultrathreads/service"
)

// Handler 前台 API 路由处理器
type Handler struct {
	services *service.Services
	caches   *cache.Caches
	jwtAuth  *jwt.GinJWTMiddleware
	jwtOAuth *jwt.GinJWTMiddleware
}

// NewHandler 创建前台路由处理器
func NewHandler(services *service.Services, caches *cache.Caches, jwtAuth, jwtOAuth *jwt.GinJWTMiddleware) *Handler {
	return &Handler{
		services: services,
		caches:   caches,
		jwtAuth:  jwtAuth,
		jwtOAuth: jwtOAuth,
	}
}

// Init 初始化App路由
func (h *Handler) Init(api *gin.RouterGroup) {
	h.initAppRoutes(api)
}
