package adapters

import (
	"database/sql"

	"github.com/syndicatedb/vodka/builders"
)

/*
Adapter - adapter intarface for DataServices
*/
type Adapter interface {
	Connect() error
	Exec(string) (sql.Result, error)
	QueryRow(string) (*sql.Row, error)
	Query(...interface{}) (*sql.Rows, error)
	Builder() builders.Builder
}

/*
Config - Database config
*/
type Config struct {
	User,
	Password,
	Host string
	Port int
	Database,
	SSLmode string
}
