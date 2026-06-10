package events

import (
	"ultrathreads/util/log" // ⚠️ 确认替换为你项目实际的 log 包路径
)

// Subscriber 抽象底层事件总线的订阅接口
// ✅ 已补上 error 返回值，与 Manager.SubscribeSafe 签名完全对齐
type Subscriber interface {
	SubscribeSafe(event string, handler func(args ...interface{})) error
}

// SubscribeTyped 类型安全的事件订阅适配器
func SubscribeTyped[T EventPayload](sub Subscriber, handler func(payload T)) {
	var zero T

	if err := sub.SubscribeSafe(zero.EventName(), func(args ...interface{}) {
		if len(args) == 0 {
			log.Error("SubscribeTyped: received empty args for event=%s", zero.EventName())
			return
		}

		payload, ok := args[0].(T)
		if !ok {
			log.Error("SubscribeTyped: type assertion failed for event=%s, got %T",
				zero.EventName(), args[0])
			return
		}

		handler(payload)
	}); err != nil {
		// ✅ 顺便处理了原来被忽略的订阅失败错误
		log.Error("SubscribeTyped: subscribe failed for event=%s, err=%v", zero.EventName(), err)
	}
}