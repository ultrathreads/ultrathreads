package form

// PostCreateForm post create form
type PostCreateForm struct {
	UserID    int64  //非表单赋值
	Title     string `form:"title" json:"title" binding:"required"`
	Content   string `form:"content" json:"content" binding:"required"`
	NodeID 	  int64 `form:"nodeId" json:"nodeId" binding:"required"`
	ParentId  int64 `form:"parentId" json:"parentId"`
	Tags      []string `form:"tags" json:"tags"`
	ImageList string `form:"imageList" json:"imageList"`
}

// PostUpdateForm post update form
type PostUpdateForm struct {
	ID       int64  //非表单赋值
	PostSlug string `form:"post_slug" json:"post_slug" binding:"required"`
	Title    string `form:"title" json:"title" binding:"required"`
	Content  string `form:"content" json:"content" binding:"required"`
	NodeID   int64  `form:"nodeId" json:"nodeId" binding: "required"`
	Tags     string `form:"tags" json:"tags"`
}
