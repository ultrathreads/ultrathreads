// bus/register.go
package bus

import (
	"ultrathreads/bus/handler"
	"ultrathreads/cache"
	"ultrathreads/repository"
)

// RegisterHandlers 集中注册所有事件处理器
func RegisterHandlers(mgr *Manager, caches *cache.Caches, repos *repository.Repositories) {

	// 设置 bus handler 需要的依赖
	handler.SetTagCache(caches.Tag)
	handler.SetReadStateCache(caches.ReadState)
	handler.SetBusDaos(repos.Tag, repos.PostTag)
	handler.SetBusPostDao(repos.Post)
	handler.SetBusUserReadStateDao(repos.UserReadState)

	// post viewed handler
	handler.PostViewedHandler(mgr)
	handler.PostCreatedHandler(mgr)
	handler.PostUpdatedHandler(mgr)
}
