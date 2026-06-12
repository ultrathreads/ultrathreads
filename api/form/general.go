package form

// form/general_get_dto.go
type GeneralGetDto struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

type IdentifierDto struct {
    Slug string `uri:"slug" binding:"required"`
}
