package {{.Package}}

import (
    "database/sql"
    "fmt"

    "github.com/posener/orm/dialect/{{.Dialect.Name}}"
)

const createString = `{{.Dialect.Create}}`

// Exec creates a table for the given struct
func (c *Create) Exec() (sql.Result, error) {
	c.orm.log("Create: '%v'", createString)
	return c.orm.db.Exec(createString)
}

func (s *Select) query() (*sql.Rows, error) {
    stmt, args := {{.Dialect.Name}}.Select(&s.internal)
	s.orm.log("Query: '%v' %v", stmt, args)
	return s.orm.db.Query(stmt, args...)
}

// Exec inserts the data to the given database
func (i *Insert) Exec() (sql.Result, error) {
	if len(i.internal.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to insert")
	}
	stmt, args := {{.Dialect.Name}}.Insert(&i.internal)
	i.orm.log("Insert: '%v' %v", stmt, args)
	return i.orm.db.Exec(stmt, args...)
}

// Exec inserts the data to the given database
func (u *Update) Exec() (sql.Result, error) {
	if len(u.internal.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
    stmt, args := {{.Dialect.Name}}.Update(&u.internal)
	u.orm.log("Update: '%v' %v", stmt, args)
	return u.orm.db.Exec(stmt, args...)
}

// Exec runs the delete statement on a given database.
func (d *Delete) Exec() (sql.Result, error) {
    stmt, args := {{.Dialect.Name}}.Delete(&d.internal)
	d.orm.log("Delete: '%v' %v", stmt, args)
	return d.orm.db.Exec(stmt, args...)
}
