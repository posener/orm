package tpls

import (
	"database/sql"
)

// Select is the struct that holds the SELECT data
type TDelete struct {
	Execer
	orm   *ORM
	where *Where
}

// Where applies where conditions on the query
func (d *TDelete) Where(w *Where) *TDelete {
	d.where = w
	return d
}

// Exec runs the delete statement on a given database.
func (d *TDelete) Exec() (sql.Result, error) {
	// create select statement
	stmt := d.String()
	args := d.where.Args()
	d.orm.log("Delete: '%v' %v", stmt, args)
	return d.orm.db.Exec(stmt, args...)
}
