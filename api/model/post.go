package model

// 话题
type Post struct {
	Model
	ThreadId          int64  `gorm:"not null;default:0;index:idx_thread_id" json:"threadId" form:"threadId"`
	ParentId          int64  `gorm:"not null;default:0;index:idx_parent_id" json:"parentId" form:"parentId"`
	Type              int    `gorm:"not null;index:idx_post_type" json:"type" form:"type"`          // 类型
	NodeId            int64  `gorm:"not null;index:idx_node_id;" json:"nodeId" form:"nodeId"`        // 节点编号
	UserId            int64  `gorm:"not null;index:idx_post_user_id;" json:"userId" form:"userId"`  // 用户
	Title             string `gorm:"size:128" json:"title" form:"title"`                             // 标题
	Content           string `gorm:"type:longtext" json:"content" form:"content"`                    // 内容
	ImageList         string `gorm:"type:longtext" json:"imageList" form:"imageList"`                // 图片
	IsPinned 		  bool   `gorm:"not null;default:false;index:idx_post_is_pinned" json:"isPinned" form:"isPinned"` // 是否置顶
	Recommend         bool   `gorm:"not null;index:idx_recommend" json:"recommend" form:"recommend"` // 是否推荐
	ViewCount         int64  `gorm:"not null" json:"viewCount" form:"viewCount"`                     // 查看数量
	LikeCount         int64  `gorm:"not null" json:"likeCount" form:"likeCount"`                     // 点赞数量
	Status            int    `gorm:"index:idx_post_status;" json:"status" form:"status"`
	LastCommentUserId int64  `gorm:"index:idx_post_last_comment_user_id" json:"lastCommentUserId" form:"lastCommentUserId"` // 最后回复时间                            // 状态：0：正常、1：删除
	LastCommentTime   int64  `gorm:"index:idx_post_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"`        // 最后回复时间
	CreateTime        int64  `gorm:"index:idx_post_create_time" json:"createTime" form:"createTime"`                        // 创建时间
	ExtraData         string `gorm:"type:text" json:"extraData" form:"extraData"`                                            // 扩展数据
}

func (p *Post) IsRoot() bool {
	return p.ParentId == 0
}

// 主题标签
type PostTag struct {
	Model
	PostId         int64 `gorm:"not null;index:idx_post_tag_post_id;" json:"postId" form:"postId"`                // 主题编号
	TagId           int64 `gorm:"not null;index:idx_post_tag_tag_id;" json:"tagId" form:"tagId"`                      // 标签编号
	Status          int64 `gorm:"not null;index:idx_post_tag_status" json:"status" form:"status"`                     // 状态：正常、删除
	LastCommentTime int64 `gorm:"index:idx_post_tag_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"` // 最后回复时间
	CreateTime      int64 `json:"createTime" form:"createTime"`                                                        // 创建时间
}

// 话题点赞
type PostLike struct {
	Model
	UserId     int64 `gorm:"not null;index:idx_post_like_user_id;" json:"userId" form:"userId"`    // 用户
	PostId    int64 `gorm:"not null;index:idx_post_like_post_id;" json:"postId" form:"postId"` // 主题编号
	CreateTime int64 `json:"createTime" form:"createTime"`                                          // 创建时间
}
