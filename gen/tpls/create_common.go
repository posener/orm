package tpls

import (
	"database/sql"
	"fmt"
)

// Create is a struct that holds data for the CREATE statement
type Create struct {
	Execer
	fmt.Stringer
	orm *ORM
}

// Exec creates a table for the given struct
func (c *Create) Exec() (sql.Result, error) {
	stmt := c.String()
	c.orm.log("Create: '%v'", stmt)
	return c.orm.db.Exec(stmt)
}
