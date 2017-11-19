package tpls

import (
	"fmt"
	"log"
)

// TInsert is a struct to hold information for an INSERT statement
type TInsert struct {
	fmt.Stringer
	cols   []string
	values []interface{}
}

// Insert returns a new INSERT statement
func Insert() *TInsert {
	return &TInsert{}
}

// Exec inserts the data to the given database
func (i *TInsert) Exec(db SQLExecer) error {
	if len(i.cols) == 0 || len(i.values) == 0 {
		return fmt.Errorf("nothing to insert")
	}

	stmt := i.String()
	log.Printf("Insert: '%v' (%v)", stmt, i.values)
	_, err := db.Exec(stmt, i.values...)
	return err
}

func (i *TInsert) add(name string, value interface{}) *TInsert {
	i.cols = append(i.cols, name)
	i.values = append(i.values, value)
	return i
}
