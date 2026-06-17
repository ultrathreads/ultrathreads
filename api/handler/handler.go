package handler

import (
	"net/http"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"ultrathreads/bus"
	"ultrathreads/middleware"
	"ultrathreads/service"
)

type Handler struct {
	services *service.Services
	mgr      *bus.Manager
	
	// ✅ 将 JWT 中间件作为 Handler 的成员变量，保证全局单例
	jwtAuth  *jwt.GinJWTMiddleware
	jwtOAuth *jwt.GinJWTMiddleware
}

func NewHandlers(services *service.Services, mgr *bus.Manager) *Handler {
	h := &Handler{
		services: services,
		mgr:      mgr,
	}
	
	// ✅ 在构造阶段统一初始化，后续所有路由共享同一实例
	h.jwtAuth = middleware.JwtAuth(middleware.LoginStandard)
	h.jwtOAuth = middleware.JwtAuth(middleware.LoginOAuth)
	
	return h
}

// Init 组装引擎、全局中间件，并调用 setupApp 注册所有路由
func (h *Handler) Init() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(middleware.Cors())
	engine.Use(func(c *gin.Context) {
		c.Set(bus.BusKey, h.mgr.Bus)
		c.Next()
	})

	engine.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "UltraThreads API\n")
	})

	// ✅ 直接使用成员变量，不再传递局部参数
	h.setupApp(engine)
	h.setupAdmin(engine)

	return engine
}