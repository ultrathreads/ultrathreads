// bus/core/subscribe.go
package core

import "ultrathreads/util/log"

type Subscriber interface {
	SubscribeSafe(event string, handler func(args ...interface{})) error
}

// SafeSubscriber 定义了安全订阅事件的能力
// 这样 handler 包就不需要再 import 顶层的 bus 包了
type SafeSubscriber interface {
	SubscribeSafe(topic string, handler func(...interface{})) error
}

func SubscribeTyped[T any](sub Subscriber, handler func(payload T)) {
	var zero T
	eventName := AutoEventName(zero)

	if err := sub.SubscribeSafe(eventName, func(args ...interface{}) {
		if len(args) == 0 {
			log.Error("SubscribeTyped: received empty args for event=%s", eventName)
			return
		}
		payload, ok := args[0].(T)
		if !ok {
			log.Error("SubscribeTyped: type assertion failed for event=%s, got %T", eventName, args[0])
			return
		}
		handler(payload)
	}); err != nil {
		log.Error("SubscribeTyped: subscribe failed for event=%s, err=%v", eventName, err)
	}
}