package querybuilder

import (
	"gorm.io/gorm"

	"ultrathreads/util/log"
)

type QueryBuilder struct {
	SelectCols []string     // 要查询的字段，如果为空，表示查询所有字段
	Params     []ParamPair  // 参数
	Orders     []OrderByCol // 排序
	Paging     *Paging      // 分页
}

// NewQueryBuilder 创建新的查询构建器
func NewQueryBuilder(selectCols ...string) *QueryBuilder {
	s := &QueryBuilder{}
	if len(selectCols) > 0 {
		s.SelectCols = append(s.SelectCols, selectCols...)
	}
	return s
}

func (s *QueryBuilder) Eq(column string, args ...interface{}) *QueryBuilder {
	s.Where(column+" = ?", args...)
	return s
}

func (s *QueryBuilder) NotEq(column string, args ...interface{}) *QueryBuilder {
	s.Where(column+" <> ?", args...)
	return s
}

func (s *QueryBuilder) Gt(column string, args ...interface{}) *QueryBuilder {
	s.Where(column+" > ?", args...)
	return s
}

func (s *QueryBuilder) Gte(column string, args ...interface{}) *QueryBuilder {
	s.Where(column+" >= ?", args...)
	return s
}

func (s *QueryBuilder) Lt(column string, args ...interface{}) *QueryBuilder {
	s.Where(column+" < ?", args...)
	return s
}

func (s *QueryBuilder) Lte(column string, args ...interface{}) *QueryBuilder {
	s.Where(column+" <= ?", args...)
	return s
}

func (s *QueryBuilder) Like(column string, str string) *QueryBuilder {
	s.Where(column+" LIKE ?", "%"+str+"%")
	return s
}

func (s *QueryBuilder) Starting(column string, str string) *QueryBuilder {
	s.Where(column+" LIKE ?", str+"%")
	return s
}

func (s *QueryBuilder) Ending(column string, str string) *QueryBuilder {
	s.Where(column+" LIKE ?", "%"+str)
	return s
}

func (s *QueryBuilder) In(column string, params interface{}) *QueryBuilder {
	s.Where(column+" IN (?)", params)
	return s
}

func (s *QueryBuilder) Where(query string, args ...interface{}) *QueryBuilder {
	s.Params = append(s.Params, ParamPair{Query: query, Args: args})
	return s
}

func (s *QueryBuilder) Asc(column string) *QueryBuilder {
	s.Orders = append(s.Orders, OrderByCol{Column: column, Asc: true})
	return s
}

func (s *QueryBuilder) Desc(column string) *QueryBuilder {
	s.Orders = append(s.Orders, OrderByCol{Column: column, Asc: false})
	return s
}

func (s *QueryBuilder) Limit(limit int) *QueryBuilder {
	s.Page(1, limit)
	return s
}

func (s *QueryBuilder) Page(page, pageSize int) *QueryBuilder {
	if s.Paging == nil {
		s.Paging = &Paging{Page: page, PageSize: pageSize}
	} else {
		s.Paging.Page = page
		s.Paging.PageSize = pageSize
	}
	return s
}

// Build 将查询条件应用到 GORM v2 DB 实例
func (s *QueryBuilder) Build(db *gorm.DB) *gorm.DB {
	ret := db.Session(&gorm.Session{}) // ✅ v2 最佳实践：创建新会话避免污染全局 db

	if len(s.SelectCols) > 0 {
		ret = ret.Select(s.SelectCols)
	}

	// where
	for _, param := range s.Params {
		ret = ret.Where(param.Query, param.Args...)
	}

	// order
	for _, order := range s.Orders {
		direction := "DESC"
		if order.Asc {
			direction = "ASC"
		}
		ret = ret.Order(order.Column + " " + direction)
	}

	// paging
	if s.Paging != nil && s.Paging.PageSize > 0 {
		ret = ret.Limit(s.Paging.PageSize)
	}
	if s.Paging != nil && s.Paging.Offset() > 0 {
		ret = ret.Offset(s.Paging.Offset())
	}

	return ret
}

// Find 执行列表查询
func (s *QueryBuilder) Find(db *gorm.DB, out interface{}) {
	if err := s.Build(db).Find(out).Error; err != nil {
		log.Error("QueryBuilder.Find error: %v", err)
	}
}

// FindOne 查询单条记录
func (s *QueryBuilder) FindOne(db *gorm.DB, out interface{}) error {
	// ✅ v2 中查询单条应使用 First 而非 Find+Limit
	return s.Build(db).First(out).Error
}

// Count 统计总数
func (s *QueryBuilder) Count(db *gorm.DB, model interface{}) int64 {
	ret := db.Session(&gorm.Session{}).Model(model)

	for _, param := range s.Params {
		ret = ret.Where(param.Query, param.Args...)
	}

	var count int64 // ✅ v2 推荐 int64，避免大表溢出
	if err := ret.Count(&count).Error; err != nil {
		log.Error("QueryBuilder.Count error: %v", err)
	}
	return count
}