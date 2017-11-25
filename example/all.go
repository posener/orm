package example

import "time"

//go:generate orm -name All -dialect sqlite3

// All is to test generation of variant fields and types
type All struct {
	Int        int    `sql:"primary key;autoincrement"`
	String     string `sql:"type:VARCHAR(100);not null"`
	Bool       bool
	unexported int

	Time time.Time

	Select int // test a case where field is a reserved name
}
