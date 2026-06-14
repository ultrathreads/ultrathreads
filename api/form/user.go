package form

// UserUpdateForm user update form
type UserUpdateForm struct {
	Slug 	    string   `uri:"slug" json:"slug" binding:"required"`
	Nickname    string `form:"nickname" json:"nickname" binding:"required"`
	Avatar      string `form:"avatar" json:"avatar"`
	Website     string `form:"website" json:"website"`
	Description string `form:"description" json:"description"`
	Level       int    `form:"level" json:"level"`
}
