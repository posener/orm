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

// Conn is a database connection interface
// It has common functions for a database connection and a database transaction
type Conn interface {
	// Driver returns the SQL driver name
	Driver() string
	// Logger sets a logger for SQL queries
	Logger(Logger)
	// Logf write to logger
	Logf(format string, args ...interface{})

	// ExecContext executes an SQL statement
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	// QueryContext executes an SQL query statement
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)

	// Close database connection
	// will panic if in current connection is a transaction
	Close() error
	// Begin a database transaction
	// It returns a new Conn object
	// will panic if in current connection is a transaction
	Begin(context.Context, *sql.TxOptions) (Conn, error)

	// Commit a transaction
	// Will panic if the current connection is not a transaction
	Commit() error
	// Rollback a transaction
	// Will panic if the current connection is not a transaction
	Rollback() error
}

// Open returns a new database for orm libraries
func Open(driverName, address string) (Conn, error) {
	sqlDB, err := sql.Open(driverName, address)
	if err != nil {
		return nil, err
	}
	return &conn{DB: sqlDB, name: driverName}, nil
}

type conn struct {
	*sql.DB
	*sql.Tx
	name string
	log  Logger
}

func (c *conn) Begin(ctx context.Context, opts *sql.TxOptions) (Conn, error) {
	t, err := c.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	tx := *c // copy connection object
	tx.DB = nil
	tx.Tx = t
	return &tx, nil
}

func (c *conn) ExecContext(ctx context.Context, stmt string, args ...interface{}) (sql.Result, error) {
	if c.Tx != nil {
		return c.Tx.ExecContext(ctx, stmt, args...)
	}
	return c.DB.ExecContext(ctx, stmt, args...)
}

func (c *conn) QueryContext(ctx context.Context, stmt string, args ...interface{}) (*sql.Rows, error) {
	if c.Tx != nil {
		return c.Tx.QueryContext(ctx, stmt, args...)
	}
	return c.DB.QueryContext(ctx, stmt, args...)
}

// Driver returns the driver name
func (c *conn) Driver() string {
	return c.name
}

func (c *conn) Logger(log Logger) {
	c.log = log
}

func (c *conn) Logf(format string, args ...interface{}) {
	if c.log == nil {
		return
	}
	c.log(format, args...)
}

// Logger is a fmt.Printf - like function
type Logger func(string, ...interface{})
