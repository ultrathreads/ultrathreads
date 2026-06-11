// bus/event/payload.go
package event

// PostViewed 节点被浏览事件载荷
type PostViewed struct {
	UserID     int64 `json:"user_id"`
	NodeID     int64 `json:"node_id"`
	PostID     int64 `json:"post_id"`
	ViewedTime int64 `json:"viewed_time"`
}