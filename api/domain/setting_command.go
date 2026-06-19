package domain

// UpdateSettingsCommand 更新系统设置命令
type UpdateSettingsCommand struct {
	SiteTitle       string
	SiteDescription string
	SiteKeywords    *string
	SiteNavs        *string
	DefaultNodeId   int
	RecommendTags   []string
}
