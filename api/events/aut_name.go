// events/auto_name.go
package events

import (
	"reflect"
	"strings"
	"sync"
)

// CustomNamer 可选接口：当结构体命名不符合约定时，手动覆盖事件名
type CustomNamer interface {
	CustomEventName() string
}

var eventCache sync.Map

// AutoEventName 根据结构体类型自动生成事件名
// NodeViewedPayload -> node.viewed
// UserCreatedPayload -> user.created
func AutoEventName(payload any) string {
	// 1. 优先检查是否实现了自定义命名
	if cn, ok := payload.(CustomNamer); ok {
		return cn.CustomEventName()
	}

	t := reflect.TypeOf(payload)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	name := t.Name()

	// 2. 查缓存
	if cached, ok := eventCache.Load(name); ok {
		return cached.(string)
	}

	// 3. 去掉 Payload 后缀
	eventName := strings.TrimSuffix(name, "Payload")

	// 4. 大驼峰转 dot.lower：NodeViewed -> node.viewed
	var b strings.Builder
	for i, c := range eventName {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				b.WriteByte('.')
			}
			b.WriteByte(byte(c + 'a' - 'A'))
		} else {
			b.WriteByte(byte(c))
		}
	}

	result := b.String()
	eventCache.Store(name, result)
	return result
}