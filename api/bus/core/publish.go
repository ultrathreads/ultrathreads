// bus/core/publish.go
package core

import "ultrathreads/util/log"

type Publisher interface {
	Publish(event string, args ...interface{})
}

func PublishTyped[T any](pub Publisher, payload T) {
	if pub == nil {
		log.Error("PublishTyped failed: publisher is nil, event=%s", AutoEventName(payload))
		return
	}
	pub.Publish(AutoEventName(payload), payload)
}