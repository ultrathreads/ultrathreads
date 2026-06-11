// bus/handler/post_created.go
package handler

import (
	"ultrathreads/bus/core" 
	"ultrathreads/bus/event"
	"ultrathreads/util/log"
)

func PostCreatedHandler(sub core.SafeSubscriber) {
	core.SubscribeTyped(sub, func(payload event.PostCreated) {
		log.Debug("payload=%v", payload)
	})
}