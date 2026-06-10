package events

// RegisterHandlers 集中注册所有事件处理器
func RegisterHandlers(mgr *Manager) {
	NodeViewedHandler(mgr)
}