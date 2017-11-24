package tpls

import (
	"database/sql"
	"fmt"
)

// TUpdate is a struct to hold information for an INSERT statement
type TUpdate struct {
	TInsert
	where *Where
}

func (u *TUpdate) Where(where *Where) *TUpdate {
	u.where = where
	return u
}

// Exec inserts the data to the given database
func (u *TUpdate) Exec() (sql.Result, error) {
	if len(u.assign) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}

	stmt := u.String()
	args := append(u.TInsert.Args(), u.where.Args()...)
	u.orm.log("Update: '%v' %v", stmt, args)
	return u.orm.db.Exec(stmt, args...)
}

// add adds a column and value to the UPDATE statement
func (u *TUpdate) add(name string, value interface{}) *TUpdate {
	u.TInsert.add(name, value)
	return u
}
