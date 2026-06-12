package model

import (
	"html/template"
)

// UserInfo 用户信息响应
type UserInfo struct {
	Slug 		 string `json:"slug"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	Level        int    `json:"level"`
	LevelName    string `json:"levelName"`
	Website      string `json:"website"`
	Description  string `json:"description"`
	Score        int    `json:"score"`
	TopicCount   int64  `json:"topicCount"`
	CommentCount int64  `json:"commentCount"`
	PasswordSet  bool   `json:"passwordSet"`
	Status       int    `json:"status"`
	CreateTime   int64  `json:"createTime"`
}

type TagResponse struct {
	Slug 	string `json:"slug"`
	TagName string `json:"tagName"`
}

type ArticleSimpleResponse struct {
	Slug 	   string 		  `json:"slug"`
	User       *UserInfo      `json:"user"`
	Tags       []TagResponse  `json:"tags"`
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
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	TopicCount  int64  `json:"topicCount"`
}

// PostSimpleResponse 帖子列表返回实体
type PostSimpleResponse struct {
	Slug 			string		   `json:"slug"`
	ThreadSlug      string         `json:"threadSlug"`
	ParentSlug      string         `json:"parentSlug"`
	IsRoot          bool  		   `json:"isRoot"`
	Type            int            `json:"type"`
	Title           string         `json:"title"`
	IsPinned        bool           `json:"isPinned"`
	LastCommentTime int64          `json:"lastCommentTime"`
	ViewCount       int64          `json:"viewCount"`
	LikeCount       int64          `json:"likeCount"`
	CreateTime      int64          `json:"createTime"`
	ImageList       []string       `json:"imageList"`
	Node            *NodeResponse  `json:"node"`
	Tags            []TagResponse  `json:"tags"`
	User            *UserInfo      `json:"user"`
	LastCommentUser *UserInfo      `json:"lastCommentUser"`
}

// PostResponse 帖子详情返回实体
type PostResponse struct {
	PostSimpleResponse
	Content template.HTML `json:"content"`
	Toc     template.HTML `json:"toc"`
}

type FavoriteResponse struct {
	Slug 	   string    `json:"slug"`
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
	Slug 		 string     `json:"slug"`
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