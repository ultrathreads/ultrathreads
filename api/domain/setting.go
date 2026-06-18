package domain

// Setting 系统配置领域模型
type Setting struct {
	ID          int64
	Key         string
	Value       string
	Name        string
	Description string
	CreateTime  int64
	UpdateTime  int64
}

// SiteNav 站点导航
type SiteNav struct {
	Title string
	Url   string
}

// SiteTip 小贴士
type SiteTip struct {
	Title   string
	Content string
}

// ScoreConfig 积分配置
type ScoreConfig struct {
	PostPostScore    int
	PostCommentScore int
}

// ConfigData 配置返回结构体
type ConfigData struct {
	SiteTitle       string
	SiteDescription string
	SiteKeywords    []string
	SiteNavs        []SiteNav
	DefaultNodeId   int64
	RecommendTags   []string
}

// AppData 应用数据
type AppData struct {
	Name           string
	Version        string
	UserLevelAdmin int
}
