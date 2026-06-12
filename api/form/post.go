package form

// RootPostCreateForm 创建根帖专用
type RootPostCreateForm struct {
	NodeSlug string   `json:"nodeSlug" binding:"required"`
	Title    string   `json:"title" binding:"required,min=5,max=100"`
	Content  string   `json:"content" binding:"required"`
	Tags     []string `json:"tags"`
	ImageList string `json:"imageList"`

	// 以下字段由 Controller 内部注入，不参与 JSON 绑定与校验
	UserSlug string   `json:"-"`
}

// ReplyCreateForm 回复专用
type ReplyCreateForm struct {
	Content    string   `json:"content" binding:"required"`
	ImageList  string   `json:"imageList"`

	// 以下字段由 Controller 内部注入，不参与 JSON 绑定与校验
	Title      string   `json:"-"`
	ParentSlug string   `json:"-"`
	UserSlug   string   `json:"-"`
}

// PostCreateForm post create form
type PostCreateForm struct {
	UserSlug   string   //非表单赋值
	Title      string   `form:"title" json:"title" binding:"required"`
	Content    string   `form:"content" json:"content" binding:"required"`
	NodeSlug   string   `form:"nodeSlug" json:"nodeSlug"`
	ParentSlug string   `form:"parentSlug" json:"parentSlug"`
	Tags       []string `form:"tags" json:"tags"`
	ImageList  string   `form:"imageList" json:"imageList"`
}

// PostUpdateForm post update form
type PostUpdateForm struct {
	Slug     string   //非表单赋值
	Title    string   `form:"title" json:"title" binding:"required"`
	Content  string   `form:"content" json:"content" binding:"required"`
	NodeSlug string   `form:"nodeSlug" json:"nodeSlug"`
	Tags     string `form:"tags" json:"tags"`
}
