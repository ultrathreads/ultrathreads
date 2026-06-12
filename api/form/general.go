package form

// form/general_get_dto.go
type GeneralGetDto struct {
    ID int64 `uri:"id" json:"id" binding:"required,min=1"`
}

type IdentifierDto struct {
    Slug string `uri:"slug" json:"slug" binding:"required"`
}
