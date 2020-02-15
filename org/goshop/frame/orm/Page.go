package orm

type Page struct {
	PageNo          int
	PageSize        int
	TotalCount      int
	PageCount       int
	FirstPage       bool
	HasPrev         bool
	HasNext         bool
	LastPage        bool
	selectpagecount bool
}
