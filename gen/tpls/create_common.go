package tpls

import (
	"database/sql"
)

// TCreate is a struct that holds data for the CREATE statement
type TCreate struct {
	Execer
	orm *ORM
}

// Exec creates a table for the given struct
func (c *TCreate) Exec() (sql.Result, error) {
	stmt := c.String()
	c.orm.log("Create: '%v'", stmt)
	return c.orm.db.Exec(stmt)
}
