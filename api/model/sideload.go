package model

import (
	"html/template"
	"ultrathreads/util/querybuilder"
)

// 替换原 PostSimpleResponse
type PostItem struct {
	Slug     string `json:"slug"`
	ThreadSlug string `json:"threadSlug"`
	ParentSlug string `json:"parentSlug"`
	NodeSlug string `json:"nodeSlug"`
	UserSlug string `json:"userSlug"`
	TagSlugs []string `json:"tagSlugs,omitempty"`

	Title    string `json:"title"`
	Content    template.HTML `json:"content,omitempty"`

	IsPinned      bool        `json:"isPinned"`
	IsRoot        bool        `json:"isRoot"`

	CreateTime      int64     `json:"createTime"`
	LastCommentTime int64     `json:"lastCommentTime"`

	ViewCount        int64    `json:"viewCount,omitempty"`
	//先考虑兼容
	LastCommentUser *UserInfo `json:"lastCommentUser"`
}

type Context struct {
    LastReadAtMap map[string]int64 `json:"lastReadAtMap,omitempty"`
}

// 侧载-用户（主key：slug）
type UserIncluded struct {
	Slug     string `json:"slug"` // 对外唯一key
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// 侧载-板块（主key：slug）
type NodeIncluded struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// 侧载-标签（主key：slug）
type TagIncluded struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	// 如有其他字段如 Color, Icon 等可在此补充
}

type PostIncluded struct {
    Users []UserIncluded `json:"users"`
    Nodes []NodeIncluded `json:"nodes"`
    Tags  []TagIncluded  `json:"tags"`
}

// 最终返回体
type PostListWithIncluded struct {
	Data  []PostItem              `json:"data,omitempty"`
	Meta     querybuilder.Paging  `json:"meta,omitempty"`
	Context  Context `json:"context,omitempty"`
	LastRead map[string]int64     `json:"lastReadAtMap,omitempty"`
	Included PostIncluded `json:"included,omitempty"`
	Extra any `json:"extra,omitempty"`
}
