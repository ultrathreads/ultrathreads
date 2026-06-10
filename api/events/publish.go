// events/publish.go
package events

import "ultrathreads/util/log"

type Publisher interface {
	Publish(event string, args ...interface{})
}

// ✅ T 不再需要实现 EventPayload 接口
func PublishTyped[T any](pub Publisher, payload T) {
	if pub == nil {
		log.Error("PublishTyped failed: publisher is nil, event=%s", AutoEventName(payload))
		return
	}
	pub.Publish(AutoEventName(payload), payload)
}