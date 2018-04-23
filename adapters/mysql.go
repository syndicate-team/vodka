package adapters

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/syndicatedb/vodka/builders"
)

/*
MySQL - low-level Postgres adapters for DataServices
*/
type MySQL struct {
	driverName     string
	config         Config
	conn           *sql.DB
	connectionInfo string
}

/*
NewMySQL - adapter constructor
*/
func NewMySQL(config Config) *MySQL {
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
func (db MySQL) Builder() builders.Builder {
	return builders.NewPostgres()
}

/*
Exec - executing SQL-query and returning *Rows
*/
func (db *MySQL) Exec(SQL string) (res sql.Result, err error) {
	if err = db.checkConnection(); err != nil {
		return
	}
	res, err = db.conn.Exec(SQL)
	if err != nil {
		if isInvalidConnection(err) {
			db.closeConnection()
			return db.Exec(SQL)
		}
	}
	return
}

/*
Query - preparing query into Statement and executing SQL-query and returning *Rows
*/
func (db *MySQL) Query(v ...interface{}) (rows *sql.Rows, err error) {
	if err = db.checkConnection(); err != nil {
		return
	}
	SQL := v[0].(string)
	values := v[1:]
	rows, err = db.conn.Query(SQL, values...)
	if err != nil {
		if isInvalidConnection(err) {
			db.closeConnection()
			return db.Query(v...)
		}
	}
	return
}

/*
QueryRow - executing single row query. May be suitable for INSERT/UPDATE.
*/
func (db *MySQL) QueryRow(SQL string) (row *sql.Row, err error) {
	if err := db.checkConnection(); err != nil {
		return nil, err
	}
	row, err = db.conn.QueryRow(SQL), nil
	if err != nil {
		if isInvalidConnection(err) {
			db.closeConnection()
			return db.QueryRow(SQL)
		}
	}
	return
}

func isInvalidConnection(err error) bool {
	return strings.Index(err.Error(), "invalid connection") != -1
}

func (db *MySQL) connect() error {
	config := db.config
	db.connectionInfo = fmt.Sprintf("%s:%s@tcp(%s:%v)/%v?charset=utf8mb4,utf8",
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

func (db *MySQL) closeConnection() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}
