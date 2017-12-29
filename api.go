package orm

import (
	"context"
	"database/sql"
	"errors"
)

// Errors exported by ORM package
var (
	ErrNotFound = errors.New("not Found")
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

// DB is an interface of functions of sql.DB which are used by orm struct.
type DB interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	Close() error
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

// Logger is a fmt.Printf - like function
type Logger func(string, ...interface{})

// GlobalLogger sets orm's global logger
// Running this function in parallel to query execution will result in
// race condition, please prepare the logger beforehand.
func GlobalLogger(l Logger) {
	log = l
}
