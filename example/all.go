package example

import "time"

//go:generate orm -name All

// All is to test generation of variant fields and types
type All struct {
	Int        int    `sql:"primary key"`
	String     string `sql:"type:VARCHAR(100);not null"`
	Bool       bool
	unexported int

	Auto int `sql:"auto increment"`

	Time time.Time

	// test a case where field is a reserved name
	Select int
}
