package domain

import "time"

// Notification 消息领域模型
type Notification struct {
	ID           int64
	FromId       int64
	UserId       int64
	Content      string
	QuoteContent string
	Type         int
	ExtraData    string
	Status       int
	CreatedAt    time.Time
}
