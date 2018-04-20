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
	Join(JoinParam) Builder
	Order(OrderParam) Builder
	Build() string
}

/*
JoinParam - joining Repository (table to query)
*/
type JoinParam struct {
	Source   string
	SourceID string
	Fields   []string
	On       []JoinParamOn
	Type     string
}

// JoinParamOn - join params and conditions
type JoinParamOn struct {
	Source    string
	SourceKey string
	JoinKey   string
	JoinValue interface{}
}

// OrderParam - ordering params
type OrderParam struct {
	OrderBy string
	Asc     bool
	Desc    bool
}
