package querybuilder

// Paging 分页请求数据
type Paging struct {
	Page     int   `json:"page"`     // 页码
	PageSize int   `json:"pageSize"` // 每页条数
	Total    int64 `json:"total"`    // 改为 int64，匹配 GORM v2 Count 返回值，防止大表溢出
}

func (p *Paging) Offset() int {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Paging) TotalPage() int64 {
	if p.Total == 0 || p.PageSize == 0 {
		return 0
	}
	totalPage := p.Total / int64(p.PageSize)
	if p.Total%int64(p.PageSize) > 0 {
		totalPage++
	}
	return totalPage
}

func GetPaging(page, limit int) *Paging {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &Paging{Page: page, PageSize: limit}
}

// ParamPair 查询条件对
type ParamPair struct {
	Query string        // 查询语句
	Args  []interface{} // 参数列表
}

// OrderByCol 排序信息
type OrderByCol struct {
	Column string // 排序字段
	Asc    bool   // 是否正序
}

// PageResult 分页返回数据
type PageResult struct {
	Page    *Paging     `json:"page"`    // 分页信息
	Results interface{} `json:"results"` // 数据
}

// CursorResult 游标分页返回数据
type CursorResult struct {
	Results interface{} `json:"results"` // 数据
	Cursor  string      `json:"cursor"`  // 下一页游标
}