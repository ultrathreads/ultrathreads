package form

// form/general_get_dto.go
type GeneralGetDto struct {
    ID int64 `uri:"id" binding:"required,min=1"` // ✅ 必须是 uri tag，不是 json/form
}