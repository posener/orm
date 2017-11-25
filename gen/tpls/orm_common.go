package tpls

import (
	"database/sql"

	"github.com/posener/orm"
)

type ORM struct {
	dialect orm.Dialect
	db      *sql.DB
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
	orm.Insert
	orm *ORM
}

// Update is a struct to hold information for an INSERT statement
type Update struct {
	orm.Update
	orm *ORM
}

func (u *Update) Where(where orm.Where) *Update {
	u.Update.Where = where
	return u
}

// Delete is the struct that holds the SELECT data
type Delete struct {
	orm.Delete
	orm *ORM
}

// Where applies where conditions on the query
func (d *Delete) Where(w orm.Where) *Delete {
	d.Delete.Where = w
	return d
}
