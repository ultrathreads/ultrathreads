package model

var Models = []interface{}{
	// 用户与权限
	&User{}, &Role{}, &Permission{}, &UserRole{}, &RolePermission{},

	// 标签
	&Tag{},

	// 文章/帖子
	&Article{}, &ArticleTag{},
	&Post{}, &Node{}, &PostTag{}, &PostLike{},

	// 用户行为
	&Favorite{}, &Notification{},
	&UserReadState{}, &UserWatch{}, &UserScore{}, &UserScoreLog{},

	// 站点
	&Setting{}, &Link{}, &LoginSource{}, &Sitemap{},
}

type Model struct {
	ID int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

const (
	StatusOk      = 0 // 正常
	StatusDeleted = 1 // 删除
	StatusPending = 2 // 待审核

	UserLevelGeneral = 0  // 普通用户
	UserLevelAdmin   = 10 // 管理员

	ContentTypeHtml     = "html"
	ContentTypeMarkdown = "markdown"

	EntityTypeArticle = "article"
	EntityTypePost    = "post"
	EntityTypeUser    = "user"

	NotificationStatusUnread = 0 // 消息未读
	NotificationStatusReaded = 1 // 消息已读

	MsgTypeComment   = 0 // 回复消息
	MsgTypePostLike = 1 // 话题点赞
	MsgTypeUserWatch = 2 // 用户关注

	LoginSourceTypeGithub = "github"
	LoginSourceTypeGitee  = "gitee"
	LoginSourceTypeQQ     = "qq"

	ScoreTypeIncr = 0 // 积分+
	ScoreTypeDecr = 1 // 积分-

	PostTypeNormal  = 0 // 普通帖子
	PostTypeTwitter = 1 // 推文
)
