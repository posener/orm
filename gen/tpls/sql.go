package tpls

import "database/sql"

type SQLExecer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type SQLQuerier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
