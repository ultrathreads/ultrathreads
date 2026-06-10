package events

import (
	"ultrathreads/util/log" // ⚠️ 替换为你项目实际的 log 包路径
)

// Publisher 抽象底层事件总线接口
// ⚠️ 签名必须与你项目中 EventBus.Bus 的实际发布方法完全一致
type Publisher interface {
	Publish(event string, args ...interface{})
}

// PublishTyped 类型安全的事件发布函数
// T 被约束为 EventPayload，事件名由 payload 自身提供，杜绝魔法字符串
func PublishTyped[T EventPayload](pub Publisher, payload T) {
	if pub == nil {
		log.Error("PublishTyped failed: publisher is nil, event=%s", payload.EventName())
		return
	}

	pub.Publish(payload.EventName(), payload)
}