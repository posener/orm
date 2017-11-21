package tpls

import (
	"database/sql"
)

// Execer is the interface for SQL update operations
type Execer interface {
	Exec() (sql.Result, error)
}

// Argser is an interface for returning a list of arguments
// This is used to pass to SQL statement
type Argser interface {
	Args() []interface{}
}
