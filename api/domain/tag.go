package domain

import "time"

// Tag 标签领域模型
type Tag struct {
	ID          int64
	Name        string
	Description string
	Status      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
