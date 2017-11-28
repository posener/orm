package tpls

import (
	"context"
	"database/sql"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect"
)

// DB is an interface of functions of sql.DB which are used by orm struct.
type DB interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	Close() error
}

// Logger is a fmt.Printf - like function
type Logger func(string, ...interface{})

// orm represents an orm of a given struct.
// All functions available to interact with an SQL table that is related
// to this struct, are done by an instance of this object.
// To get an instance of orm use Open or New functions.
type orm struct {
	dialect dialect.Dialect
	db      DB
	logger  Logger
}

func (o *orm) Close() error {
	return o.db.Close()
}

// Logger sets a logger to the orm package
func (o *orm) Logger(logger Logger) {
	o.logger = logger
}

// CreateBuilder builds an SQL CREATE statement parameters
type CreateBuilder struct {
	params common.CreateParams
	orm    *orm
}

// IfNotExists sets IF NOT EXISTS for the CREATE SQL statement
func (b *CreateBuilder) IfNotExists() *CreateBuilder {
	b.params.IfNotExists = true
	return b
}

// Context sets the context for the SQL query
func (b *CreateBuilder) Context(ctx context.Context) *CreateBuilder {
	b.params.Ctx = ctx
	return b
}

// InsertBuilder builds an INSERT statement parameters
type InsertBuilder struct {
	params common.InsertParams
	orm    *orm
}

// Context sets the context for the SQL query
func (b *InsertBuilder) Context(ctx context.Context) *InsertBuilder {
	b.params.Ctx = ctx
	return b
}

// UpdateBuilder builds SQL INSERT statement parameters
type UpdateBuilder struct {
	params common.UpdateParams
	orm    *orm
}

// Where sets the WHERE statement to the SQL query
func (b *UpdateBuilder) Where(where common.Where) *UpdateBuilder {
	b.params.Where = where
	return b
}

// Context sets the context for the SQL query
func (b *UpdateBuilder) Context(ctx context.Context) *UpdateBuilder {
	b.params.Ctx = ctx
	return b
}

// DeleteBuilder builds SQL DELETE statement parameters
type DeleteBuilder struct {
	params common.DeleteParams
	orm    *orm
}

// Where applies where conditions on the SQL query
func (b *DeleteBuilder) Where(w common.Where) *DeleteBuilder {
	b.params.Where = w
	return b
}

// Context sets the context for the SQL query
func (b *DeleteBuilder) Context(ctx context.Context) *DeleteBuilder {
	b.params.Ctx = ctx
	return b
}

// log if a logger was set
func (o *orm) log(s string, args ...interface{}) {
	if o.logger == nil {
		return
	}
	o.logger(s, args...)
}
