// bus/event/payload.go
package event

import "fmt"

// PostViewed 节点被浏览事件载荷
type PostViewed struct {
	UserID     int64 `json:"userId"`
	PostSlug   string `json:"postSlug"`
	NodeSlug   string `json:"nodeSlug"`
	ViewedTime int64  `json:"viewedTime"`
}

func (p PostViewed) String() string {
    return fmt.Sprintf("PostViewed(PostSlug=%s, NodeSlug=%s, UserID=%d, ViewedTime=%d)", p.PostSlug, p.NodeSlug, p.UserID, p.ViewedTime)
}

// PostCreated 帖子创建成功事件载荷
type PostCreated struct {
	UserID int64  `json:"user_id"`
	PostID int64  `json:"post_id"`
	IsRoot bool   `json:"is_root"`
	Tags []string `json:"tags"`
}

func (p PostCreated) String() string {
    return fmt.Sprintf("PostCreated(PostID=%d, UserID=%d, , IsRoot=%t, Tags=%v)", p.PostID, p.UserID, p.IsRoot, p.Tags)
}

// PostUpdated 帖子更新成功事件载荷
type PostUpdated struct {
	UserID int64  `json:"user_id"`
	PostID int64  `json:"post_id"`
	IsRoot bool   `json:"is_root"`
	Tags []string `json:"tags"`
}

func (p PostUpdated) String() string {
	return fmt.Sprintf("PostUpdated(PostID=%d, UserID=%d, IsRoot=%t, Tags=%v)", p.PostID, p.UserID, p.IsRoot, p.Tags)
}
