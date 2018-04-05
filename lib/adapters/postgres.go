package adapters

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

const driverName = "postgres"

/*
Postgres - low-level Postgres adapters for DataServices
*/
type Postgres struct {
	Config         DBConfig
	Conn           *sql.DB
	Source         string
	ConnectionInfo string
}

/*
NewPostgres - adapter constructor
*/
func NewPostgres(config DBConfig) *Postgres {
	return &Postgres{
		Config: config,
	}
}

/*
Connect - public method to connect.
Not very useful because all methods checking connections and connecting by default
*/
func (psql *Postgres) Connect() error {
	return psql.connect()
}

/*
Builder - returns Query builder (SQL) instance
*/
func (psql Postgres) Builder() Builder {
	return &SQLBuilder{}
}

/*
Exec - executing SQL-query and returning Result
*/
func (psql *Postgres) Exec(SQL string) (sql.Result, error) {
	if err := psql.checkConnection(); err != nil {
		return nil, err
	}
	return psql.Conn.Exec(SQL)
}

/*
Query - preparing query into Statement and executing SQL-query and returning *Rows
*/
func (psql *Postgres) Query(v ...interface{}) (*sql.Rows, error) {
	if err := psql.checkConnection(); err != nil {
		return nil, err
	}
	SQL := v[0].(string)
	values := v[1:]
	return psql.Conn.Query(SQL, values...)
}

/*
QueryRow - executing single row query. May be suitable for INSERT/UPDATE.
*/
func (psql *Postgres) QueryRow(SQL string) (*sql.Row, error) {
	if err := psql.checkConnection(); err != nil {
		return nil, err
	}
	return psql.Conn.QueryRow(SQL), nil
}

func (psql *Postgres) connect() error {
	config := psql.Config
	psql.ConnectionInfo = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.SSLmode,
	)
	log.Println("Connecting to Postgres: ", psql.ConnectionInfo)
	conn, err := sql.Open("postgres", psql.ConnectionInfo)
	if err != nil {
		fmt.Println("Postgres connection error", err)
		return err
	}
	if conn == nil {
		fmt.Println("Connection to postgres is nil")
	}
	psql.Conn = conn
	return nil
}

func (psql *Postgres) checkConnection() error {
	fmt.Printf("Connection: %+v\n", psql.Conn)
	if psql.Conn == nil {
		return psql.connect()
	}
	return nil
}
