package model

import "ultrathreads/util/querybuilder"

// 替换原 PostSimpleResponse
type PostItem struct {
	Slug     string `json:"slug"`
	ThreadSlug string `json:"threadSlug"`
	ParentSlug string `json:"parentSlug"`
	NodeSlug string `json:"nodeSlug"`
	UserSlug string `json:"userSlug"`
	Title    string `json:"title"`

	CreateTime      int64          `json:"createTime"`
	LastCommentTime int64          `json:"lastCommentTime"`

	//先考虑兼容
	Node            *NodeResponse  `json:"node"`
	Tags            []TagResponse  `json:"tags"`
	User            *UserInfo      `json:"user"`
	LastCommentUser *UserInfo      `json:"lastCommentUser"`
}

// 侧载-用户（主key：slug）
type UserIncluded struct {
	Slug   string `json:"slug"` // 对外唯一key
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// 侧载-板块（主key：slug）
type NodeIncluded struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// 最终返回体
type PostListWithIncluded struct {
	Data  []PostItem              `json:"data"`
	Meta     querybuilder.Paging  `json:"meta"`
	LastRead map[string]int64     `json:"lastReadAtMap"`
	Included struct {
		Users []UserIncluded `json:"users"`
		Nodes []NodeIncluded `json:"nodes"`
	} `json:"included"`
}