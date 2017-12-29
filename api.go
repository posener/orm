package orm

import (
	"context"
	"database/sql"
	"errors"
)

// Errors exported by ORM package
var (
	ErrNotFound = errors.New("not found")
)

// Op is an SQL comparison operation
type Op string

// Operators for SQL WHERE statements
const (
	OpEq   Op = "="
	OpNe   Op = "<>"
	OpGt   Op = ">"
	OpGE   Op = ">="
	OpLt   Op = "<"
	OpLE   Op = "<="
	OpLike Op = "LIKE"
)

// OrderDir is direction in which a column can be ordered
type OrderDir string

// Directions for SQL ORDER BY statements
const (
	Asc  OrderDir = "ASC"
	Desc OrderDir = "DESC"
)

// DB is SQL database interface
type DB interface {
	// *sql.DB APIs
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	Close() error

	// Driver returns the SQL driver name
	Driver() string
	// Logger sets a logger for SQL queries
	Logger(Logger)
	// Logf write to logger
	Logf(format string, args ...interface{})
}

// Open returns a new database for orm libraries
func Open(driverName, address string) (DB, error) {
	sqlDB, err := sql.Open(driverName, address)
	if err != nil {
		return nil, err
	}
	return &db{DB: sqlDB, name: driverName}, nil
}

type db struct {
	*sql.DB
	name string
	log  Logger
}

// Driver returns the driver name
func (d *db) Driver() string {
	return d.name
}

func (d *db) Logger(log Logger) {
	d.log = log
}

func (d *db) Logf(format string, args ...interface{}) {
	if d.log == nil {
		return
	}
	d.log(format, args...)
}

// Logger is a fmt.Printf - like function
type Logger func(string, ...interface{})
