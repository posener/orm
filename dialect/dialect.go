package dialect

import (
	"context"
	"fmt"
	"strings"

	"github.com/posener/orm"
	"github.com/posener/orm/dialect/migration"
	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/postgres"
	"github.com/posener/orm/dialect/sqlite3"
	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/load"
)

// dialect is an interface to represent an SQL dialect
// Objects that implement this interface, can convert query params, such as SelectParams or
// UpdateParams, and convert them to an SQL statement and a list of arguments, which can be used
// to invoke SQL Exec or Query functions.
type API interface {
	// Name returns the name of the dialect
	Name() string
	// Create returns the SQL CREATE statement and arguments according to the given parameters
	Create(orm.Conn, *CreateParams) ([]string, error)
	// Insert returns the SQL INSERT statement and arguments according to the given parameters
	Insert(*InsertParams) (string, []interface{})
	// Select returns the SQL SELECT statement and arguments according to the given parameters
	Select(*SelectParams) (string, []interface{})
	// Delete returns the SQL DELETE statement and arguments according to the given parameters
	Delete(*DeleteParams) (string, []interface{})
	// Update returns the SQL UPDATE statement and arguments according to the given parameters
	Update(*UpdateParams) (string, []interface{})
	// Drop returns the SQL DROP statement and arguments according to the given parameters
	Drop(*DropParams) (string, []interface{})
}

// Dialect is an interface for a dialect for generating ORM code
type Dialect interface {
	// Name returns the name of the dialect
	Name() string
	// GoTypeToColumnType gets a string that represents a go basic type
	// and returns an SQL type for a createColumn for a field of that type.
	GoTypeToColumnType(goType string, autoIncrement bool) *sqltypes.Type
	// Translate gets a MySQL statement and returns a corresponding statement
	// in a specific dialect
	Translate(string) string
	// ConvertValueCode returns code for converting a value for a field with
	// a given SQL type.
	ConvertValueCode(*load.Field) string
	// Quote returns the quoted form of an SQL variable
	Quote(string) string
	// ReplaceVars replaces question marks from sql query to the right variable of the dialect
	ReplaceVars(s string) string
}

var dialects = map[string]API{
	"mysql":    &dialect{name: "mysql", Dialect: new(mysql.Dialect)},
	"postgres": &dialect{name: "postgres", Dialect: new(postgres.Dialect)},
	"sqlite3":  &dialect{name: "sqlite3", Dialect: new(sqlite3.Dialect)},
}

// All returns all available dialects
func All() []API {
	var all []API
	for _, d := range dialects {
		all = append(all, d)
	}
	return all
}

// Get returns a dialect by name
func Get(name string) API {
	return dialects[name]
}

type dialect struct {
	Dialect
	name string
}

// Create returns the SQL CREATE statement and arguments according to the given parameters
func (d *dialect) Create(conn orm.Conn, p *CreateParams) ([]string, error) {
	table := new(migration.Table)
	err := table.UnMarshal(p.MarshaledTable)
	if err != nil {
		return nil, err
	}

	if p.AutoMigrate {
		stmts, ok, err := d.autoMigrate(p.Ctx, conn, p.Table, table)
		if err != nil {
			return nil, fmt.Errorf("automigration: %s", err)
		}
		if ok {
			return stmts, nil
		}
	}
	stmt := fmt.Sprintf("CREATE TABLE %s %s (%s)",
		d.ifNotExists(p.IfNotExists),
		d.Quote(p.Table),
		d.tableProperties(table),
	)
	return []string{stmt}, nil
}

func (d *dialect) autoMigrate(ctx context.Context, conn orm.Conn, tableName string, want *migration.Table) ([]string, bool, error) {
	got, err := migration.Load(ctx, conn, tableName)
	if err != nil {
		// XXX: Here we assume error is: table does not exists
		// if it is not, we should return the error and not nil
		return nil, false, nil
	}
	diff, err := got.Diff(want)
	if err != nil {
		return nil, false, fmt.Errorf("bad conditions: %s", err)
	}

	for len(diff.PrimaryKeys) > 0 {
		return nil, false, fmt.Errorf("not supported: add primary keys %s", diff.PrimaryKeys)
	}

	var stmts []string
	for _, col := range diff.Columns {
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s",
			d.Quote(tableName), d.createColumn(col)))
	}
	for _, fk := range diff.ForeignKeys {
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s",
			d.Quote(tableName), d.foreignKey(fk)))
	}
	return stmts, true, nil
}

// Insert returns the SQL INSERT statement and arguments according to the given parameters
func (d *dialect) Insert(p *InsertParams) (string, []interface{}) {
	stmt := d.ReplaceVars(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		d.Quote(p.Table),
		d.assignColumns(p.Assignments),
		QMarks(len(p.Assignments)),
	))

	var args []interface{}
	if p.Assignments != nil {
		args = append(args, p.Assignments.Args()...)
	}

	return stmt, args
}

// Select returns the SQL SELECT statement and arguments according to the given parameters
func (d *dialect) Select(p *SelectParams) (string, []interface{}) {
	stmt := d.ReplaceVars(fmt.Sprintf("SELECT %s FROM %s %s %s %s %s %s",
		d.selectColumns(p),
		d.Quote(p.Table),
		d.join(p),
		d.whereJoin(p.Table, p),
		d.groupBy(p.Table, p),
		d.orderBy(p.Table, p),
		d.page(p.Page),
	))

	return stmt, collectWhereArgs(p)
}

// collectWhereArgs collects arguments for WHERE statement from
// select params and all its nested join options
func collectWhereArgs(p *SelectParams) []interface{} {
	var args []interface{}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}
	for _, join := range p.Joins {
		args = append(args, collectWhereArgs(&join.SelectParams)...)
	}
	return args
}

// Delete returns the SQL DELETE statement and arguments according to the given parameters
func (d *dialect) Delete(p *DeleteParams) (string, []interface{}) {
	stmt := d.ReplaceVars(fmt.Sprintf("DELETE FROM %s %s",
		d.Quote(p.Table),
		d.where(p.Table, p.Where),
	))

	var args []interface{}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}

// Update returns the SQL UPDATE statement and arguments according to the given parameters
func (d *dialect) Update(p *UpdateParams) (string, []interface{}) {
	stmt := d.ReplaceVars(fmt.Sprintf("UPDATE %s SET %s %s",
		d.Quote(p.Table),
		d.assignSets(p.Assignments),
		d.where(p.Table, p.Where),
	))

	var args []interface{}
	if p.Assignments != nil {
		args = append(args, p.Assignments.Args()...)
	}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}

// Drop returns the SQL DROP statement and arguments according to the given parameters
func (d *dialect) Drop(p *DropParams) (string, []interface{}) {
	stmt := fmt.Sprintf("DROP TABLE %s %s",
		d.ifExists(p.IfExists),
		d.Quote(p.Table),
	)
	return stmt, nil
}

// join extract SQL join list statement
func (d *dialect) join(p *SelectParams) string {
	return strings.Join(d.joinParts(p.Table, p), " ")
}

func (d *dialect) joinParts(table string, p *SelectParams) []string {
	joins := p.Joins
	if len(joins) == 0 {
		return nil
	}
	var (
		tables    []string
		conds     []string
		recursive []string
	)
	for _, j := range joins {
		joinTable := j.TableName(table)
		tables = append(tables, fmt.Sprintf("%s AS %s", d.Quote(j.Table), d.Quote(joinTable)))
		for _, pairing := range j.Pairings {
			conds = append(conds, fmt.Sprintf("%s.%s = %s.%s",
				d.Quote(table), d.Quote(pairing.Column), d.Quote(joinTable), d.Quote(pairing.JoinedColumn)))
		}
		recursive = append(recursive, d.joinParts(j.TableName(table), &j.SelectParams)...)
	}

	var joinStmt string
	for i := range tables {
		joinStmt += fmt.Sprintf("LEFT OUTER JOIN %s ON (%s) ", tables[i], conds[i])
	}
	return append([]string{joinStmt}, recursive...)
}
