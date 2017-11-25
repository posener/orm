package tpls

import (
	"context"
	"database/sql"

	"github.com/posener/orm/common"
)

// DB is an interface of functions of sql.DB which are used by ORM struct.
type DB interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	Close() error
}

// ORM represents an ORM of a given struct.
// All functions available to interact with an SQL table that is related
// to this struct, are done by an instance of this object.
// To get an instance of ORM use Open or New functions.
type ORM struct {
	dialect common.Dialect
	db      DB
	logger  Logger
}

func (o *ORM) Close() error {
	return o.db.Close()
}

// Create returns a struct for a CREATE statement
func (o *ORM) Create() *Create {
	return &Create{orm: o}
}

// Logger sets a logger to the ORM package
func (o *ORM) Logger(logger Logger) {
	o.logger = logger
}

// Create is a struct that holds data for the CREATE statement
type Create struct {
	orm *ORM
}

// Insert is a struct to hold information for an INSERT statement
type Insert struct {
	internal common.Insert
	orm      *ORM
}

// Update is a struct to hold information for an INSERT statement
type Update struct {
	internal common.Update
	orm      *ORM
}

func (u *Update) Where(where common.Where) *Update {
	u.internal.Where = where
	return u
}

// Delete is the struct that holds the SELECT data
type Delete struct {
	internal common.Delete
	orm      *ORM
}

// Where applies where conditions on the query
func (d *Delete) Where(w common.Where) *Delete {
	d.internal.Where = w
	return d
}
