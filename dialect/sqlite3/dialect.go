package sqlite3

import (
	"fmt"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/format"
)

// Dialect represents the sqlite dialect
type Dialect struct{}

// Name returns the name of the dialect
func (*Dialect) Name() string {
	return "sqlite3"
}

// Create returns the SQL CREATE statement and arguments according to the given parameters
func (d *Dialect) Create(p *common.CreateParams) (string, []interface{}) {
	stmt := fmt.Sprintf(`CREATE TABLE %s '%s' ( %s )`,
		format.IfNotExists(p.IfNotExists),
		p.Table,
		p.ColumnsStatement,
	)

	return stmt, nil
}

// Insert returns the SQL INSERT statement and arguments according to the given parameters
func (d *Dialect) Insert(p *common.InsertParams) (string, []interface{}) {
	stmt := fmt.Sprintf(`INSERT INTO '%s' (%s) VALUES (%s)`,
		p.Table,
		format.AssignColumns(p.Assignments),
		common.QMarks(len(p.Assignments)),
	)

	var args []interface{}
	if p.Assignments != nil {
		args = append(args, p.Assignments.Args()...)
	}

	return stmt, args
}

// Select returns the SQL SELECT statement and arguments according to the given parameters
func (d *Dialect) Select(p *common.SelectParams) (string, []interface{}) {
	stmt := fmt.Sprintf("SELECT %s FROM '%s' %s %s %s %s %s",
		format.Columns(p.Table, p.Columns),
		p.Table,
		format.Where(p.Where),
		format.GroupBy(p.Table, p.Groups),
		format.OrderBy(p.Table, p.Orders),
		format.Page(p.Page),
	)

	var args []interface{}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}

// Delete returns the SQL DELETE statement and arguments according to the given parameters
func (d *Dialect) Delete(p *common.DeleteParams) (string, []interface{}) {
	stmt := fmt.Sprintf("DELETE FROM '%s' %s",
		p.Table,
		format.Where(p.Where),
	)

	var args []interface{}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}

// Update returns the SQL UPDATE statement and arguments according to the given parameters
func (d *Dialect) Update(p *common.UpdateParams) (string, []interface{}) {
	stmt := fmt.Sprintf(`UPDATE '%s' SET %s %s`,
		p.Table,
		format.AssignSets(p.Assignments),
		format.Where(p.Where),
	)

	var args []interface{}
	if p.Assignments != nil {
		args = append(args, p.Assignments.Args()...)
	}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}
