package builders

// NewPostgres - Postgres SQL builder
func NewPostgres() Builder {
	return &postgres{}
}

/*
Builder - interface for query builder for adapter and data service
*/
type Builder interface {
	Select([]string) Builder
	Insert(string) Builder
	Update(string) Builder
	Delete() Builder
	ReturnID(string) Builder
	Values(interface{}) Builder
	Set(interface{}) Builder
	From(string) Builder
	Where(map[string]interface{}) Builder
	Limit(int, int) Builder
	Join(Join) Builder
	Order(OrderParam) Builder
	Build() string
}

/*
Join - joining Repository (table to query)
*/
type Join struct {
	Source    string
	Key       string
	TargetKey string
	Fields    []string
	Type      string
}

// OrderParam - ordering params
type OrderParam struct {
	OrderBy string
	Asc     bool
	Desc    bool
}
