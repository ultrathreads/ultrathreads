package model

import "time"

// 站点地图
type Sitemap struct {
	Model
	Loc       string    `gorm:"not null;size:1024" json:"loc" form:"loc"`              // loc
	Lastmod   int64     `gorm:"not null" json:"lastmod" form:"lastmod"`                // 最后更新时间
	LocName   string    `gorm:"not null;size:32;unique" json:"locName" form:"locName"` // loc的md5
	CreatedAt time.Time `gorm:"not null" json:"createdAt" form:"createdAt"`            // 创建时间
}
