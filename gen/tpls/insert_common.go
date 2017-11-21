package tpls

import (
	"database/sql"
	"fmt"
)

// TInsert is a struct to hold information for an INSERT statement
type TInsert struct {
	Execer
	Argser
	fmt.Stringer
	orm    *ORM
	cols   []string
	values []interface{}
}

// Exec inserts the data to the given database
func (i *TInsert) Exec() (sql.Result, error) {
	if len(i.cols) == 0 || len(i.values) == 0 {
		return nil, fmt.Errorf("nothing to insert")
	}

	stmt := i.String()
	args := i.Args()
	i.orm.log("Insert: '%v' %v", stmt, args)
	return i.orm.db.Exec(stmt, args...)
}

// Args returns a list of arguments for the INSERT statement
func (i *TInsert) Args() []interface{} {
	return i.values
}

func (i *TInsert) add(name string, value interface{}) *TInsert {
	i.cols = append(i.cols, name)
	i.values = append(i.values, value)
	return i
}
