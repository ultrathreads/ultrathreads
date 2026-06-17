package dto

type IdRequest struct {
    ID int64 `uri:"id" json:"id" binding:"required,min=1"`
}

// 定义与前端 SiteSettings 对应的结构体
type SettingsRequest struct {
	SiteTitle       string   `json:"siteTitle" binding:"required"`
	SiteDescription string   `json:"siteDescription"`
	SiteKeywords    *string  `json:"siteKeywords"`
	SiteNavs        *string  `json:"siteNavs"`
	DefaultNodeId   int      `json:"defaultNodeId" binding:"required,min=1"`
	RecommendTags   []string `json:"recommendTags"`
}

// NodeCreateForm node create form
type NodeCreateForm struct {
	Name        string `form:"name" json:"name" binding:"required,max=32"`
	Description string `form:"description" json:"description"`
	Icon        string `form:"icon" json:"icon"`
	SortNo      int    `form:"sortNo" json:"sortNo"`
	Status      int    `form:"status" json:"status"`
}

// NodeUpdateForm node update form
type NodeUpdateForm struct {
	ID 	   	    int64  `uri:"id" json:"id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required,max=32"`
	Description string `form:"description" json:"description"`
	Icon        string `form:"icon" json:"icon"`
	SortNo      int    `form:"sortNo" json:"sortNo"`
	Status      int    `form:"status" json:"status"`
}

type SortRequest struct {
    Items []struct {
        ID     int `json:"id" binding:"required"`
        SortNo int `json:"sortNo" binding:"required"`
    } `json:"items" binding:"required,min=1"`
}


// TagCreateForm tag create form
type TagCreateForm struct {
	Name        string `form:"name" json:"name" binding:"required,max=32"`
	Description string `form:"description" json:"description"`
	Status      int    `form:"status" json:"status"`
}

// TagUpdateForm tag update form
type TagUpdateForm struct {
	ID 	   	    int64  `uri:"id" json:"id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required,max=32"`
	Description string `form:"description" json:"description"`
	Status      int    `form:"status" json:"status"`
}

// LinkUpdateForm link update form
type LinkUpdateForm struct {
	ID      int64  //非表单赋值
	Title   string `form:"title" json:"title" binding:"required"`
	URL     string `form:"url" json:"url" binding:"required"`
	Logo    string `form:"logo" json:"logo"`
	Summary string `form:"summary" json:"summary"`
	Status  int    `form:"status" json:"status"`
}

// LinkCreateForm node create form
type LinkCreateForm struct {
	Title   string `form:"title" json:"title" binding:"required"`
	URL     string `form:"url" json:"url" binding:"required"`
	Logo    string `form:"logo" json:"logo"`
	Summary string `form:"summary" json:"summary"`
	Status  int    `form:"status" json:"status"`
}

// CommentUpdateForm comment update form
type CommentUpdateForm struct {
	ID      int64  //非表单赋值
	Content string `form:"content" json:"content" binding:"required,min=3"`
	Status  int    `form:"status" json:"status"`
}
