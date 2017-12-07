package tpls

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/posener/orm"
	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect"
)

// BEGIN - the generator will drop everything above this line

// conn represents a DB connection for manipulating a given struct.
// All functions available to interact with an SQL table that is related
// to this struct, are done by an instance of this object.
// To get an instance of orm use Open or New functions.
type conn struct {
	dialect dialect.Dialect
	db      orm.DB
	logger  orm.Logger
}

func (c *conn) Close() error {
	return c.db.Close()
}

// Logger sets a logger to the conn package
func (c *conn) Logger(logger orm.Logger) {
	c.logger = logger
}

// CreateBuilder builds an SQL CREATE statement parameters
type CreateBuilder struct {
	params common.CreateParams
	conn   *conn
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
	conn   *conn
}

// Context sets the context for the SQL query
func (b *InsertBuilder) Context(ctx context.Context) *InsertBuilder {
	b.params.Ctx = ctx
	return b
}

// UpdateBuilder builds SQL INSERT statement parameters
type UpdateBuilder struct {
	params common.UpdateParams
	conn   *conn
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
	conn   *conn
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
func (c *conn) log(s string, args ...interface{}) {
	if c.logger == nil {
		return
	}
	c.logger(s, args...)
}

func parseInt(s []byte) int64 {
	i, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		log.Printf("Failed parsing %s to int", string(s))
	}
	return i
}

func parseUInt(s []byte) uint64 {
	i, err := strconv.ParseUint(string(s), 10, 64)
	if err != nil {
		log.Printf("Failed parsing %s to uint", string(s))
	}
	return i
}

func parseFloat(s []byte) float64 {
	i, err := strconv.ParseFloat(string(s), 64)
	if err != nil {
		log.Printf("Failed parsing %s to float", string(s))
	}
	return i
}

func parseTime(s []byte, precision int) time.Time {
	format := "2006-01-02 15:04:05"
	if precision > 0 {
		format += "." + strings.Repeat("0", precision)
	}
	t, err := time.Parse(format, string(s))
	if err != nil {
		log.Printf("Failed parsing '%s' to time.Time", string(s))
	}
	return t
}

func parseBool(s []byte) bool {
	return s[0] != 0
}
