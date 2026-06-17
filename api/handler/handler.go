package handler

import (
	"github.com/gin-gonic/gin"
	"ultrathreads/bus"
	"ultrathreads/router"
	"ultrathreads/service"
)

type Handler struct {
	services *service.Services
	mgr      *bus.Manager
}

func NewHandlers(services *service.Services, mgr *bus.Manager) *Handler {
	return &Handler{
		services: services,
		mgr:      mgr,
	}
}

// Init 组装引擎并注册所有路由
func (h *Handler) Init() *gin.Engine {
	engine := gin.Default()

	// ✅ 路由注册内聚到 Handler 内部，直接使用 h.services 和 h.mgr
	router.Setup(engine, h.mgr, h.services)

	return engine
}