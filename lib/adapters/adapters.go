package adapters

import "database/sql"

/*
Adapter - adapter intarface for DataServices
*/
type Adapter interface {
	Connect() error
	Exec(string) (sql.Result, error)
	QueryRow(string) (*sql.Row, error)
	Query(...interface{}) (*sql.Rows, error)
	Builder() Builder
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

// Source — why do I need that?
type Source string

// fields - WTF?
type fields string

/*
DBConfig - Database config
*/
type DBConfig struct {
	User,
	Password,
	Host string
	Port int
	Database,
	SSLmode string
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
