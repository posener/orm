package dialect

import (
	"fmt"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/format"
	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/sqlite3"
)

// dialect is an interface to represent an SQL dialect
// Objects that implement this interface, can convert query params, such as SelectParams or
// UpdateParams, and convert them to an SQL statement and a list of arguments, which can be used
// to invoke SQL Exec or Query functions.
type Dialect interface {
	// Name returns the name of the dialect
	Name() string
	// Create returns the SQL CREATE statement and arguments according to the given parameters
	Create(*common.CreateParams) (string, []interface{})
	// Insert returns the SQL INSERT statement and arguments according to the given parameters
	Insert(*common.InsertParams) (string, []interface{})
	// Select returns the SQL SELECT statement and arguments according to the given parameters
	Select(*common.SelectParams) (string, []interface{})
	// Delete returns the SQL DELETE statement and arguments according to the given parameters
	Delete(*common.DeleteParams) (string, []interface{})
	// Update returns the SQL UPDATE statement and arguments according to the given parameters
	Update(*common.UpdateParams) (string, []interface{})
}

// DialectImplementer are methods that a dialect implementer should implement
type DialectImplementer interface {
	// Name returns the name of the dialect
	Name() string
	// TableQuote takes table name and return the name quoted, according to the syntax
	TableQuote(string) string
}

// New returns a new dialect according to it's name
func New(name string) (Dialect, error) {
	switch name {
	case "mysql":
		return &dialect{DialectImplementer: new(mysql.Dialect)}, nil
	case "sqlite3":
		return &dialect{DialectImplementer: new(sqlite3.Dialect)}, nil
	default:
		return nil, fmt.Errorf("unsupported dialect %s", name)
	}
}

// dialect represents the sqlite dialect
type dialect struct {
	DialectImplementer
}

// Create returns the SQL CREATE statement and arguments according to the given parameters
func (d *dialect) Create(p *common.CreateParams) (string, []interface{}) {
	stmt := fmt.Sprintf(`CREATE TABLE %s %s ( %s )`,
		format.IfNotExists(p.IfNotExists),
		d.TableQuote(p.Table),
		p.ColumnsStatement,
	)

	return stmt, nil
}

// Insert returns the SQL INSERT statement and arguments according to the given parameters
func (d *dialect) Insert(p *common.InsertParams) (string, []interface{}) {
	stmt := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
		d.TableQuote(p.Table),
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
func (d *dialect) Select(p *common.SelectParams) (string, []interface{}) {
	stmt := fmt.Sprintf("SELECT %s FROM %s %s %s %s %s %s",
		format.Columns(p.Table, p.Columns),
		d.TableQuote(p.Table),
		format.Join(p.Table, p.Columns),
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
func (d *dialect) Delete(p *common.DeleteParams) (string, []interface{}) {
	stmt := fmt.Sprintf("DELETE FROM %s %s",
		d.TableQuote(p.Table),
		format.Where(p.Where),
	)

	var args []interface{}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}

// Update returns the SQL UPDATE statement and arguments according to the given parameters
func (d *dialect) Update(p *common.UpdateParams) (string, []interface{}) {
	stmt := fmt.Sprintf(`UPDATE %s SET %s %s`,
		d.TableQuote(p.Table),
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
