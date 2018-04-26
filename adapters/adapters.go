package adapters

import (
	"database/sql"
	"time"

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
KVAdapter - Key/value adapter intarface for DataServices
*/
type KVAdapter interface {
	Connect() error
	Get(key string) ([]byte, error)
	Set(key string, value interface{}, expiry time.Duration) error
	SetJSON(key string, value interface{}, expiry time.Duration) error
	Del(key string) error
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
