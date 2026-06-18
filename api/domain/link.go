package domain

import "time"

// Link 友链领域模型
type Link struct {
	ID        int64
	Url       string
	Title     string
	Summary   string
	Logo      string
	Status    int
	CreatedAt time.Time
}
