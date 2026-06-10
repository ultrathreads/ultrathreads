package events

import (
	"context"
	"fmt"
	"sync"

	evbus "github.com/asaskevich/EventBus"
)

// BusKey 是存储在 Gin Context 中的 key
const BusKey = "event_bus"

// Manager 管理事件总线的生命周期
type Manager struct {
	Bus  evbus.Bus
	once sync.Once
	done chan struct{}
}

// NewManager 创建一个带优雅停机支持的事件管理器
func NewManager() *Manager {
	return &Manager{
		Bus:  evbus.New(),
		done: make(chan struct{}),
	}
}

// SubscribeSafe 安全异步订阅，shutdown 期间自动丢弃新事件
func (m *Manager) SubscribeSafe(topic string, handler func(...interface{})) error {
	return m.Bus.SubscribeAsync(topic, func(args ...interface{}) {
		select {
		case <-m.done:
			fmt.Printf("[events] ⚠️ dropped '%s' during shutdown\n", topic)
		default:
			handler(args...)
		}
	}, true) // transactional=true: 同 topic handler 串行执行
}

// Shutdown 等待所有异步 handler 完成
func (m *Manager) Shutdown(ctx context.Context) error {
	m.once.Do(func() { close(m.done) })

	waitCh := make(chan struct{})
	go func() {
		m.Bus.WaitAsync()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		fmt.Println("[events] ✅ all handlers completed")
		return nil
	case <-ctx.Done():
		fmt.Println("[events] ❌ shutdown timeout")
		return ctx.Err()
	}
}