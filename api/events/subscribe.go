// events/subscribe.go
package events

import "ultrathreads/util/log"

type Subscriber interface {
	SubscribeSafe(event string, handler func(args ...interface{})) error
}

// ✅ T 不再需要实现 EventPayload 接口
func SubscribeTyped[T any](sub Subscriber, handler func(payload T)) {
	var zero T
	eventName := AutoEventName(zero) // 用零值获取事件名，无副作用

	if err := sub.SubscribeSafe(eventName, func(args ...interface{}) {
		if len(args) == 0 {
			log.Error("SubscribeTyped: received empty args for event=%s", eventName)
			return
		}
		payload, ok := args[0].(T)
		if !ok {
			log.Error("SubscribeTyped: type assertion failed for event=%s, got %T",
				eventName, args[0])
			return
		}
		handler(payload)
	}); err != nil {
		log.Error("SubscribeTyped: subscribe failed for event=%s, err=%v", eventName, err)
	}
}