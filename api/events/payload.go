// events/payload.go
package events

type NodeViewedPayload struct {
	UserID     int64 `json:"user_id"`
	NodeID     int64 `json:"node_id"`
	PostID	   int64 `json:"post_id"`
	ViewedTime int64 `json:"viewed_time"`
}
