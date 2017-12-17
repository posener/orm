package orm

import (
	"context"
	"database/sql"
	"errors"

	"github.com/posener/orm/runtime"
)

// Errors exported by ORM package
var (
	ErrNotFound = errors.New("Not Found")
)

// Operators for SQL WHERE statements
const (
	OpEq   runtime.Op = "="
	OpNe   runtime.Op = "<>"
	OpGt   runtime.Op = ">"
	OpGE   runtime.Op = ">="
	OpLt   runtime.Op = "<"
	OpLE   runtime.Op = "<="
	OpLike runtime.Op = "LIKE"
)

// Directions for SQL ORDER BY statements
const (
	Asc  runtime.OrderDir = "ASC"
	Desc runtime.OrderDir = "DESC"
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
