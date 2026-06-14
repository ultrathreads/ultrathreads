package controller

type Pagination struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"pageSize"`
	TotalItems int  `json:"totalItems"`
	HasPrev    bool `json:"hasPrev"`
	HasNext    bool `json:"hasNext"`
}

// DataEnvelope 仅承载旧格式兼容的列表数据，不含任何具体业务字段
type DataEnvelope[T any] struct {
	Results T `json:"results"`
}

// ListResponse 三个平级字段，完全通用
type ListResponse[T any] struct {
	Data    *DataEnvelope[T] `json:"data"`
	Meta    Pagination       `json:"meta"`
	Context interface{}      `json:"context,omitempty"`
}

func NewListResponse[T any](
	items T,
	page, pageSize, total int,
	ctx interface{},
) ListResponse[T] {
	return ListResponse[T]{
		Meta: Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			HasPrev:    page > 1,
			HasNext:    page*pageSize < total,
		},
		Context: ctx, // 直接透传，不解析、不假设内部结构
		Data: &DataEnvelope[T]{
			Results: items,
		},
	}
}