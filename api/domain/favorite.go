package domain

import "time"

// Favorite 收藏领域模型
type Favorite struct {
	ID         int64
	UserId     int64
	EntityType string
	EntityId   int64
	CreatedAt  time.Time
}
