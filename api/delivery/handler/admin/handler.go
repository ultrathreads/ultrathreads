package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt/v2"

	"ultrathreads/cache"
	"ultrathreads/service"
)

// Handler 后台管理 API 路由处理器
type Handler struct {
	services *service.Services
	caches   *cache.Caches
	jwtAuth  *jwt.GinJWTMiddleware
}

// NewHandler 创建后台管理路由处理器
func NewHandler(services *service.Services, caches *cache.Caches, jwtAuth *jwt.GinJWTMiddleware) *Handler {
	return &Handler{
		services: services,
		caches:   caches,
		jwtAuth:  jwtAuth,
	}
}

// Init 初始化Admin路由
func (h *Handler) Init(api *gin.RouterGroup) {
	h.initAdminRoutes(api)
}

