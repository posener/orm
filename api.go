package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
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

	// SQL DB actions
	Exec(context.Context, string, ...interface{}) (sql.Result, error)
	Query(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRow(context.Context, string, ...interface{}) *sql.Row

	// Non transaction functions

	// Close database connection
	// will panic if in current connection is a transaction
	Close() error
	// Begin a database transaction
	// It returns a new Conn object
	// will panic if in current connection is a transaction
	Begin(context.Context, *sql.TxOptions) (Conn, error)
	// ConnDB returns the underlying sql connection
	// Will panic if the current connection is a transaction
	ConnDB() *sql.DB

	// Transaction only functions

	// Commit a transaction
	// Will panic if the current connection is not a transaction
	Commit() error
	// Rollback a transaction
	// Will panic if the current connection is not a transaction
	Rollback() error
	// ConnTx returns the underlying sql transaction connection
	// Will panic if the current connection is not a transaction
	ConnTx() *sql.Tx
}

// Logger is a fmt.Printf - like function
type Logger func(string, ...interface{})

// Open returns a new database for orm libraries
func Open(driverName, address string, options ...Option) (Conn, error) {
	sqlDB, err := sql.Open(driverName, address)
	if err != nil {
		return nil, err
	}
	c := &conn{DB: sqlDB, name: driverName}
	for _, option := range options {
		option(c)
	}
	return c, nil
}

// Option is a Conn modifier function that is used when opening a new connection
type Option func(c *conn)

// OptLogger sets a logger for the connection
func OptLogger(log Logger) Option {
	return func(c *conn) {
		c.log = log
	}
}

type conn struct {
	*sql.DB
	*sql.Tx
	name string
	log  Logger
	txID int
}

func (c *conn) Begin(ctx context.Context, opts *sql.TxOptions) (Conn, error) {
	tx, err := c.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	txConn := *c // copy connection object
	txConn.DB = nil
	txConn.Tx = tx
	txConn.txID = rand.Intn(9999)
	c.logf(fmt.Sprintf("Tx[%d] Begin", txConn.txID))
	return &txConn, nil
}

func (c *conn) Commit() error {
	c.logf(fmt.Sprintf("Tx[%d] Commit", c.txID))
	return c.Tx.Commit()
}

func (c *conn) Rollback() error {
	c.logf(fmt.Sprintf("Tx[%d] Rollback", c.txID))
	return c.Tx.Rollback()
}

func (c *conn) Exec(ctx context.Context, stmt string, args ...interface{}) (sql.Result, error) {
	if c.Tx != nil {
		c.logf("Tx[%d] Exec: %v %v", c.txID, stmt, args)
		return c.Tx.ExecContext(ctx, stmt, args...)
	}
	c.logf("Exec: %v %v", stmt, args)
	return c.DB.ExecContext(ctx, stmt, args...)
}

func (c *conn) Query(ctx context.Context, stmt string, args ...interface{}) (*sql.Rows, error) {
	if c.Tx != nil {
		c.logf("Tx[%d] Query: %v %v", c.txID, stmt, args)
		return c.Tx.QueryContext(ctx, stmt, args...)
	}
	c.logf("Query: %v %v", stmt, args)
	return c.DB.QueryContext(ctx, stmt, args...)
}

func (c *conn) QueryRow(ctx context.Context, stmt string, args ...interface{}) *sql.Row {
	if c.Tx != nil {
		c.logf("Tx[%d] QueryRow: %v %v", c.txID, stmt, args)
		return c.Tx.QueryRowContext(ctx, stmt, args...)
	}
	c.logf("QueryRow: %v %v", stmt, args)
	return c.DB.QueryRowContext(ctx, stmt, args...)
}

func (c *conn) ConnDB() *sql.DB {
	return c.DB
}

func (c *conn) ConnTx() *sql.Tx {
	return c.Tx
}

func (c *conn) Driver() string {
	return c.name
}

func (c *conn) logf(format string, args ...interface{}) {
	if c.log == nil {
		return
	}
	c.log(format, args...)
}
