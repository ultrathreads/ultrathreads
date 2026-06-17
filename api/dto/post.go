package dto

import (
	"html/template"
)

// ==================== 请求 DTO (Request) ====================

// CreateRootPostRequest 创建根帖请求
type CreateRootPostRequest struct {
	NodeSlug  string   `json:"nodeSlug" binding:"required"`
	Title     string   `json:"title" binding:"required,min=5,max=100"`
	Content   string   `json:"content" binding:"required"`
	Tags      []string `json:"tags"`
	ImageList []string `json:"imageList"`
}

// UpdateRootPostRequest 更新根帖请求
// URI 参数与 Body 分离，避免同一个 struct 混用 uri/json tag
type UpdateRootPostRequest struct {
	Slug      string   `uri:"slug" binding:"required"`
	Title     *string  `json:"title" binding:"omitempty"`
	Content   *string  `json:"content" binding:"omitempty"`
	NodeSlug  *string  `json:"nodeSlug" binding:"omitempty"`
	Tags      []string `json:"tags"`
	ImageList []string `json:"imageList"`
}

// CreateReplyRequest 创建回复请求
type CreateReplyRequest struct {
	ThreadSlug string   `uri:"slug" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	ImageList  []string `json:"imageList"`
	ParentSlug string   `json:"parentSlug,omitempty"`
}

// UpdateReplyRequest 更新回复请求
type UpdateReplyRequest struct {
	Slug    string  `uri:"slug" binding:"required"`
	Content *string `json:"content" binding:"omitempty"`
}

// ==================== 响应 DTO (Response) ====================

// PostLiteResponse 帖子信息（轻量）
type PostLiteResponse struct {
	Slug            string        `json:"slug"`
	ThreadSlug      string        `json:"threadSlug"`
	ParentSlug      string        `json:"parentSlug"`
	IsRoot          bool          `json:"isRoot"`
	Type            int           `json:"type"`
	Title           string        `json:"title"`
	IsPinned        bool          `json:"isPinned"`
	LastCommentTime int64         `json:"lastCommentTime"`
	ViewCount       int64         `json:"viewCount"`
	LikeCount       int64         `json:"likeCount"`
	CreatedAt       int64         `json:"createdAt"`
	NodeSlug        string        `json:"nodeSlug"`
}

// PostResponse 帖子详情（完整）
type PostResponse struct {
	PostLiteResponse                   // 复用基础字段
	RawContent string             `json:"rawContent"`
	Content    template.HTML      `json:"content"`
	Toc        template.HTML      `json:"toc"`
	ImageList  []string           `json:"imageList"`
}