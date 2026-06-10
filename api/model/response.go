package model

import (
	"html/template"
)

// UserInfo 用户信息响应（✅ 安全修复：移除敏感字段）
type UserInfo struct {
	Id           int64  `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	Level        int    `json:"level"`
	LevelName    string `json:"levelName"`
	Website      string `json:"website"`
	Description  string `json:"description"`
	Score        int    `json:"score"`
	TopicCount   int64  `json:"topicCount"`   // ✅ int → int64，与 User Model 对齐
	CommentCount int64  `json:"commentCount"` // ✅ int → int64，与 User Model 对齐
	PasswordSet  bool   `json:"passwordSet"`
	Status       int    `json:"status"`
	CreateTime   int64  `json:"createTime"`
}

type TagResponse struct {
	TagId   int64  `json:"tagId"`
	TagName string `json:"tagName"`
}

type ArticleSimpleResponse struct {
	ArticleId  int64          `json:"articleId"`
	User       *UserInfo      `json:"user"`
	Tags       []TagResponse  `json:"tags"` // ✅ *[]T → []T
	Title      string         `json:"title"`
	Summary    string         `json:"summary"`
	Share      bool           `json:"share"`
	SourceUrl  string         `json:"sourceUrl"`
	ViewCount  int64          `json:"viewCount"`
	CreateTime int64          `json:"createTime"`
}

type ArticleResponse struct {
	ArticleSimpleResponse
	Content template.HTML `json:"content"`
	Toc     template.HTML `json:"toc"`
}

type NodeResponse struct {
	NodeId      int64  `json:"nodeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	TopicCount  int64  `json:"topicCount"`
}

// PostSimpleResponse 帖子列表返回实体
type PostSimpleResponse struct {
	Id              int64          `json:"id"`
	ThreadId        int64          `json:"threadId"`
	ParentId        int64          `json:"parentId"`
	Type            int            `json:"type"`
	Title           string         `json:"title"`
	IsPinned        bool           `json:"isPinned"`
	LastCommentTime int64          `json:"lastCommentTime"`
	ViewCount       int64          `json:"viewCount"`
	LikeCount       int64          `json:"likeCount"`
	CreateTime      int64          `json:"createTime"`
	ImageList       []string       `json:"imageList"`      // ✅ *[]string → []string
	Node            *NodeResponse  `json:"node"`
	Tags            []TagResponse  `json:"tags"`           // ✅ *[]TagResponse → []TagResponse
	User            *UserInfo      `json:"user"`
	LastCommentUser *UserInfo      `json:"lastCommentUser"`
}

// PostResponse 帖子详情返回实体
type PostResponse struct {
	PostSimpleResponse
	Content template.HTML `json:"content"`
	Toc     template.HTML `json:"toc"`
}

// CommentResponse 回帖详情返回实体
type CommentResponse struct {
	CommentId    int64            `json:"commentId"`
	User         *UserInfo        `json:"user"`
	EntityType   string           `json:"entityType"`
	EntityId     int64            `json:"entityId"`
	Content      template.HTML    `json:"content"`
	QuoteId      int64            `json:"quoteId"`
	Quote        *CommentResponse `json:"quote"`
	QuoteContent template.HTML    `json:"quoteContent"`
	Status       int              `json:"status"`
	CreateTime   int64            `json:"createTime"`
}

type FavoriteResponse struct {
	FavoriteId int64     `json:"favoriteId"`
	EntityType string    `json:"entityType"`
	EntityId   int64     `json:"entityId"`
	Deleted    bool      `json:"deleted"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	User       *UserInfo `json:"user"`
	Url        string    `json:"url"`
	CreateTime int64     `json:"createTime"`
}

// NotificationResponse 消息通知响应
type NotificationResponse struct {
	MessageId    int64     `json:"messageId"`
	From         *UserInfo `json:"from"`
	UserId       int64     `json:"userId"`
	Content      string    `json:"content"`
	QuoteContent string    `json:"quoteContent"`
	Type         int       `json:"type"`
	Icon         string    `json:"icon"`
	DetailUrl    string    `json:"detailUrl"`
	ExtraData    string    `json:"extraData"`
	Status       int       `json:"status"`
	CreateTime   int64     `json:"createTime"`
}