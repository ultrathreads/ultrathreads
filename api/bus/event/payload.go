// bus/event/payload.go
package event

import "fmt"

// PostViewed 节点被浏览事件载荷
type PostViewed struct {
	UserID     int64 `json:"user_id"`
	NodeID     int64 `json:"node_id"`
	PostID     int64 `json:"post_id"`
	ViewedTime int64 `json:"viewed_time"`
}

func (p PostViewed) String() string {
    return fmt.Sprintf("PostViewed(PostID=%d, UserID=%d, NodeID=%d, ViewedTime=%d)", p.PostID, p.UserID, p.NodeID, p.ViewedTime)
}

// PostCreated 帖子创建成功事件载荷
type PostCreated struct {
	UserID int64 `json:"user_id"`
	PostID int64 `json:"post_id"`
	IsRoot bool `json:"is_root"`
}

func (p PostCreated) String() string {
    return fmt.Sprintf("PostCreated(PostID=%d, UserID=%d, IsRoot)", p.PostID, p.UserID, p.IsRoot)
}