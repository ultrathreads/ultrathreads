package render

// PostRenderConfig 渲染配置（不对外暴露）
type postRenderConfig struct {
    includeContent         bool
    includeViewCount       bool
    includeLastCommentUser bool
    // 未来可继续扩展:
    // includeReplies      bool
    // includePreviewImages bool
}

// PostRenderOption 函数式选项类型
type PostRenderOption func(*postRenderConfig)

// WithContent 启用 Content 字段渲染
func WithContent() PostRenderOption {
    return func(c *postRenderConfig) {
        c.includeContent = true
    }
}

// WithViewCount 启用 Content 字段渲染
func WithViewCount() PostRenderOption {
    return func(c *postRenderConfig) {
        c.includeViewCount = true
    }
}

// WithLastCommentUser 启用 LastCommentUser 字段渲染（预留示例）
func WithLastCommentUser() PostRenderOption {
    return func(c *postRenderConfig) {
        c.includeLastCommentUser = true
    }
}
