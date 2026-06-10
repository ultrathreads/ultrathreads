package events

// EventPayload 所有事件载荷必须实现的契约接口
type EventPayload interface {
	EventName() string
}

// ==================== 节点相关事件 ====================

// NodeViewedPayload 节点浏览事件
type NodeViewedPayload struct {
	UserID     int64 `json:"user_id"`
	NodeID     int64 `json:"node_id"`
	ViewedTime int64 `json:"viewed_time"`
}

func (NodeViewedPayload) EventName() string { return "node.viewed" }

// ==================== 用户相关事件（示例） ====================

// UserCreatedPayload 用户注册事件
type UserCreatedPayload struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
}

func (UserCreatedPayload) EventName() string { return "user.created" }