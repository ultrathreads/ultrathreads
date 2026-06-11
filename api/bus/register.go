// bus/register.go
package bus

import "ultrathreads/bus/handler"

// RegisterHandlers 集中注册所有事件处理器
func RegisterHandlers(mgr *Manager) {

	// post viewed handler
	handler.PostViewedHandler(mgr)
	handler.PostCreatedHandler(mgr) 
}