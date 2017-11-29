package orm

import (
	"context"
	"database/sql"
	"errors"

	"github.com/posener/orm/common"
)

// Errors exported by ORM package
var (
	ErrNotFound = errors.New("Not Found")
)

// Operators for SQL WHERE statements
const (
	OpEq   common.Op = "="
	OpNe   common.Op = "<>"
	OpGt   common.Op = ">"
	OpGE   common.Op = ">="
	OpLt   common.Op = "<"
	OpLE   common.Op = "<="
	OpLike common.Op = "LIKE"
)

// Directions for SQL ORDER BY statements
const (
	Asc  common.OrderDir = "ASC"
	Desc common.OrderDir = "DESC"
)

// DB is an interface of functions of sql.DB which are used by orm struct.
type DB interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	Close() error
}
