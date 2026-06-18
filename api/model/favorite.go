package model

import "time"

// 收藏
type Favorite struct {
	Model
	UserId     int64     `gorm:"index:idx_favorite_user_id;not null" json:"userId" form:"userId"`                     // 用户编号
	EntityType string    `gorm:"index:idx_favorite_entity_type;size:32;not null" json:"entityType" form:"entityType"` // 收藏实体类型
	EntityId   int64     `gorm:"index:idx_favorite_entity_id;not null" json:"entityId" form:"entityId"`               // 收藏实体编号
	CreatedAt  time.Time `json:"createdAt" form:"createdAt"`                                                          // 创建时间
}
