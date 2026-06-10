// events/payload.go
package events

// ==================== 节点相关事件 ====================
type NodeViewedPayload struct {
	UserID     int64 `json:"user_id"`
	NodeID     int64 `json:"node_id"`
	ViewedTime int64 `json:"viewed_time"`
}

// ==================== 用户相关事件 ====================
type UserCreatedPayload struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
}