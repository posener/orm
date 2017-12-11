package dialect

import (
	"fmt"
	"strings"

	"github.com/posener/orm/common"
)

const SQLite3 = "sqlite3"

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

// New returns a new dialect according to it's name
func New(name string) (Dialect, error) {
	return &dialect{name: name}, nil
}

// dialect represents the sqlite dialect
type dialect struct {
	name string
}

func (d *dialect) Name() string {
	return d.name
}

// Create returns the SQL CREATE statement and arguments according to the given parameters
func (d *dialect) Create(p *common.CreateParams) (string, []interface{}) {
	stmt := fmt.Sprintf("CREATE TABLE %s `%s` ( %s )",
		IfNotExists(p.IfNotExists),
		p.Table,
		p.ColumnsStatement,
	)

	return stmt, nil
}

// Insert returns the SQL INSERT statement and arguments according to the given parameters
func (d *dialect) Insert(p *common.InsertParams) (string, []interface{}) {
	stmt := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
		p.Table,
		assignColumns(p.Assignments),
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
	stmt := fmt.Sprintf("SELECT %s FROM `%s` %s %s %s %s %s",
		columns(p),
		p.Table,
		d.join(p),
		whereJoin(p.Table, p),
		groupBy(p.Table, p.Groups),
		orderBy(p.Table, p.Orders),
		page(p.Page),
	)

	return stmt, collectWhereArgs(p)
}

// collectWhereArgs collects arguments for WHERE statement from
// select params and all its nested join options
func collectWhereArgs(p *common.SelectParams) []interface{} {
	var args []interface{}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}
	for _, join := range p.Columns.Joins() {
		args = append(args, collectWhereArgs(&join.SelectParams)...)
	}
	return args
}

// Delete returns the SQL DELETE statement and arguments according to the given parameters
func (d *dialect) Delete(p *common.DeleteParams) (string, []interface{}) {
	stmt := fmt.Sprintf("DELETE FROM `%s` %s",
		p.Table,
		where(p.Table, p.Where),
	)

	var args []interface{}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}

// Update returns the SQL UPDATE statement and arguments according to the given parameters
func (d *dialect) Update(p *common.UpdateParams) (string, []interface{}) {
	stmt := fmt.Sprintf("UPDATE `%s` SET %s %s",
		p.Table,
		assignSets(p.Assignments),
		where(p.Table, p.Where),
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

// join extract SQL join list statement
func (d *dialect) join(p *common.SelectParams) string {
	return strings.Join(d.joinParts(p.Table, p), " ")
}

func (d *dialect) joinParts(table string, p *common.SelectParams) []string {
	joins := p.Columns.Joins()
	if len(joins) == 0 {
		return nil
	}
	var (
		tables    []string
		conds     []string
		recursive []string
	)
	for _, j := range joins {
		joinTable := j.TableName(p.Table)
		tables = append(tables, fmt.Sprintf("`%s` AS `%s`", j.Table, joinTable))
		for _, pairing := range j.Pairings {
			conds = append(conds, fmt.Sprintf("`%s`.`%s` = `%s`.`%s`", table, pairing.Column, joinTable, pairing.JoinedColumn))
		}
		recursive = append(recursive, d.joinParts(j.TableName(p.Table), &j.SelectParams)...)
	}

	tablesStmt := strings.Join(tables, ", ")
	condStmt := strings.Join(conds, " AND ")

	// sqlite3 requires the table statement not to be in braces
	if d.name != SQLite3 {
		tablesStmt = "(" + tablesStmt + ")"
	}

	joinStmt := fmt.Sprintf("JOIN %s ON (%s)", tablesStmt, condStmt)

	return append([]string{joinStmt}, recursive...)
}
