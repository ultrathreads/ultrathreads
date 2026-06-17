package dto

type SlugRequest struct {
    Slug string `uri:"slug" json:"slug" binding:"required"`
}
