package adapters

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

/*
MySQL - low-level Postgres adapters for DataServices
*/
type MySQL struct {
	driverName     string
	config         DBConfig
	conn           *sql.DB
	connectionInfo string
}

/*
NewMySQL - adapter constructor
*/
func NewMySQL(config DBConfig) *MySQL {
	return &MySQL{
		config:     config,
		driverName: "mysql",
	}
}

/*
Connect - public method to connect.
Not very useful because all methods checking connections and connecting by default
*/
func (db *MySQL) Connect() error {
	return db.connect()
}

/*
Builder - returns Query builder (SQL) instance
*/
func (db MySQL) Builder() Builder {
	return &SQLBuilder{}
}

/*
Exec - executing SQL-query and returning *Rows
*/
func (db *MySQL) Exec(SQL string) (sql.Result, error) {
	if err := db.checkConnection(); err != nil {
		return nil, err
	}
	return db.conn.Exec(SQL)
}

/*
Query - preparing query into Statement and executing SQL-query and returning *Rows
*/
func (db *MySQL) Query(v ...interface{}) (*sql.Rows, error) {
	if err := db.checkConnection(); err != nil {
		return nil, err
	}
	SQL := v[0].(string)
	values := v[1:]
	return db.conn.Query(SQL, values...)
}

/*
QueryRow - executing single row query. May be suitable for INSERT/UPDATE.
*/
func (db *MySQL) QueryRow(SQL string) (*sql.Row, error) {
	if err := db.checkConnection(); err != nil {
		return nil, err
	}
	return db.conn.QueryRow(SQL), nil
}

func (db *MySQL) connect() error {
	config := db.config
	db.connectionInfo = fmt.Sprintf("%s:%s@tcp(%s:%v)/%v",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	log.Println("Connecting to MySQL: ", db.connectionInfo)
	conn, err := sql.Open(db.driverName, db.connectionInfo)
	if err != nil {
		fmt.Println("MySQL connection error", err)
		return err
	}
	if conn == nil {
		fmt.Println("Connection to MySQL is nil")
	}
	db.conn = conn
	return nil
}

func (db *MySQL) checkConnection() error {
	if db.conn == nil {
		return db.connect()
	}
	if db.conn.Stats().OpenConnections == 0 {
		return db.connect()
	}
	fmt.Printf("Connection: %+v\n", db.conn.Stats().OpenConnections)
	return nil
}