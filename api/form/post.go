package form

// RootPostCreateForm 创建根帖专用
type RootPostCreateForm struct {
	NodeSlug  string   `json:"nodeSlug" binding:"required"`
	Title     string   `json:"title" binding:"required,min=5,max=100"`
	Content   string   `json:"content" binding:"required"`
	Tags      []string `json:"tags"`
	ImageList string `json:"imageList"`

	// 以下字段由 Controller 内部注入，不参与 JSON 绑定与校验
	UserSlug  string   `json:"-"`
}

// RootPostUpdateForm root post update form
type RootPostUpdateForm struct {
	Slug 	 string   `uri:"slug" json:"slug" binding:"required"`
	Title    string   `form:"title" json:"title"`
	Content  string   `form:"content" json:"content"`
	NodeSlug string   `form:"nodeSlug" json:"nodeSlug"`
	Tags     []string `form:"tags" json:"tags"`
}

// ReplyCreateForm 回复专用
type ReplyCreateForm struct {
	Slug 	   string   `uri:"slug" json:"slug" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	ImageList  string   `json:"imageList"`

	// 以下字段由 Controller 内部注入，不参与 JSON 绑定与校验
	Title      string   `json:"-"`
	ParentSlug string   `json:"-"`
	UserSlug   string   `json:"-"`
}

// ReplyUpdateForm reply update form
type ReplyUpdateForm struct {
	Slug 	 string   `uri:"slug" json:"slug" binding:"required"`
	Content  string   `form:"content" json:"content"`

	Title    string   `json:"-"`
}
