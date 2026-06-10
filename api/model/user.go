package model

import (
	"database/sql"
)

// User 用户模型
type User struct {
	Model
	Username     sql.NullString `gorm:"size:32;uniqueIndex" json:"username" form:"username"`                       // 用户名
	Email        sql.NullString `gorm:"size:128;uniqueIndex" json:"email" form:"email"`                            // 邮箱
	Nickname     string         `gorm:"size:16" json:"nickname" form:"nickname"`                                   // 昵称
	Avatar       string         `gorm:"type:text" json:"avatar" form:"avatar"`                                     // 头像
	Password     string         `gorm:"size:512" json:"-" form:"password"`                                         // 密码（✅ 禁止 JSON 序列化）
	Website      string         `gorm:"size:1024" json:"website" form:"website"`                                   // 个人主页
	Description  string         `gorm:"type:text" json:"description" form:"description"`                           // 个人描述
	Status       int            `gorm:"index:idx_user_status;not null" json:"status" form:"status"`                // 状态
	TopicCount   int64          `gorm:"column:post_count;not null" json:"topicCount" form:"topicCount"`            // ✅ int→int64 + 列名映射
	CommentCount int64          `gorm:"not null" json:"commentCount" form:"commentCount"`                          // ✅ int→int64
	Level        int            `gorm:"not null" json:"level" form:"level"`                                        // 用户等级
	CreateTime   int64          `json:"createTime" form:"createTime"`                                              // 创建时间
	UpdateTime   int64          `json:"updateTime" form:"updateTime"`                                              // 更新时间
}

// UserScore 用户积分
type UserScore struct {
	Model
	UserId     int64 `gorm:"uniqueIndex;not null" json:"userId" form:"userId"` // ✅ unique → uniqueIndex
	Score      int   `gorm:"not null" json:"score" form:"score"`               // 积分
	CreateTime int64 `json:"createTime" form:"createTime"`                     // 创建时间
	UpdateTime int64 `json:"updateTime" form:"updateTime"`                     // 更新时间
}

// UserWatch 用户关注
type UserWatch struct {
	Model
	UserID     int64 `gorm:"not null;index:idx_user_watch_user_id" json:"userId" form:"userId"`          // 用户
	WatcherID  int64 `gorm:"not null;index:idx_user_watch_watcher_id" json:"watcherId" form:"watcherId"` // 关注者编号
	CreateTime int64 `json:"createTime" form:"createTime"`                                               // 创建时间
}

// UserScoreLog 用户积分流水
type UserScoreLog struct {
	Model
	UserId      int64  `gorm:"not null;index:idx_user_score_log_user_id" json:"userId" form:"userId"`             // 用户编号
	SourceType  string `gorm:"not null;index:idx_user_score_source,priority:1" json:"sourceType" form:"sourceType"` // ✅ 复合索引 priority
	SourceId    string `gorm:"not null;index:idx_user_score_source,priority:2" json:"sourceId" form:"sourceId"`     // ✅ 复合索引 priority
	Description string `json:"description" form:"description"`                                                     // 描述
	Type        int    `json:"type" form:"type"`                                                                   // 类型(增加、减少)
	Score       int    `json:"score" form:"score"`                                                                 // 积分
	CreateTime  int64  `json:"createTime" form:"createTime"`                                                       // 创建时间
}