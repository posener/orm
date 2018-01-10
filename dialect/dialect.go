package dialect

import (
	"context"
	"fmt"

	"github.com/posener/orm"
	"github.com/posener/orm/dialect/migration"
	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/postgres"
	"github.com/posener/orm/dialect/sqlite3"
	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/load"
)

const (
	Mysql    = "mysql"
	Postgres = "postgres"
	Sqlite3  = "sqlite3"
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
	Var(int) string
}

var dialects = map[string]API{
	Mysql:    &dialect{name: Mysql, Dialect: new(mysql.Dialect)},
	Postgres: &dialect{name: Postgres, Dialect: new(postgres.Dialect)},
	Sqlite3:  &dialect{name: Sqlite3, Dialect: new(sqlite3.Dialect)},
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

	b := newBuilder(d, "CREATE TABLE")
	if p.IfNotExists {
		b.Append("IF NOT EXISTS")
	}
	b.Append(d.Quote(p.Table))
	b.Open()
	tableProperties(b, table)
	b.Close()
	return []string{b.Statement()}, nil
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
		b := newBuilder(d, "ALTER TABLE")
		b.Append(b.Quote(tableName))
		b.Append("ADD COLUMN")
		createColumn(b, col)
		stmts = append(stmts, b.Statement())
	}
	for _, fk := range diff.ForeignKeys {
		b := newBuilder(d, "ALTER TABLE")
		b.Append(b.Quote(tableName))
		b.Append("ADD CONSTRAINT")
		foreignKey(b, fk)
		stmts = append(stmts, b.Statement())
	}
	return stmts, true, nil
}

// Insert returns the SQL INSERT statement and arguments according to the given parameters
func (d *dialect) Insert(p *InsertParams) (string, []interface{}) {
	b := newBuilder(d, "INSERT INTO")
	b.Append(d.Quote(p.Table))
	b.Open()
	for i, assignment := range p.Assignments {
		b.Append(d.Quote(assignment.Column))
		if i != len(p.Assignments)-1 {
			b.Comma()
		}
	}
	b.Close()
	b.Append("VALUES")
	b.Open()
	for i, arg := range p.Assignments.Args() {
		b.Var(arg)
		if i != len(p.Assignments.Args())-1 {
			b.Comma()
		}
	}
	b.Close()

	for i, col := range p.RetColumns {
		if i == 0 {
			b.Append("RETURNING")
		} else {
			b.Comma()
		}
		b.Append(b.Quote(col))
	}

	return b.Statement(), b.Args()
}

// Select returns the SQL SELECT statement and arguments according to the given parameters
func (d *dialect) Select(p *SelectParams) (string, []interface{}) {
	b := newBuilder(d, "SELECT")
	selectColumns(b, p)
	b.Append("FROM")
	b.Append(d.Quote(p.Table))
	join(b, p)
	where(b, p.Table, p.Where, p.Joins)
	groupBy(b, p.Table, p)
	orderBy(b, p.Table, p)
	page(b, p.Page)
	return b.Statement(), b.Args()
}

// Delete returns the SQL DELETE statement and arguments according to the given parameters
func (d *dialect) Delete(p *DeleteParams) (string, []interface{}) {
	b := newBuilder(d, "DELETE FROM")
	b.Append(d.Quote(p.Table))
	where(b, p.Table, p.Where, nil)
	return b.Statement(), b.Args()
}

// Update returns the SQL UPDATE statement and arguments according to the given parameters
func (d *dialect) Update(p *UpdateParams) (string, []interface{}) {
	b := newBuilder(d, "UPDATE")
	b.Append(d.Quote(p.Table))
	b.Append("SET")
	for i, assign := range p.Assignments {
		b.Append(d.Quote(assign.Column))
		b.Append("=")
		b.Var(assign.ColumnValue)
		if i != len(p.Assignments)-1 {
			b.Comma()
		}
	}
	where(b, p.Table, p.Where, nil)
	return b.Statement(), b.Args()
}

// Drop returns the SQL DROP statement and arguments according to the given parameters
func (d *dialect) Drop(p *DropParams) (string, []interface{}) {
	b := newBuilder(d, "DROP TABLE")
	if p.IfExists {
		b.Append("IF EXISTS")
	}
	b.Append(d.Quote(p.Table))
	return b.Statement(), b.Args()
}

// tableProperties returns all properties of SQL table, as should be given in the table CREATE statement
func tableProperties(b *builder, t *migration.Table) {
	for i, col := range t.Columns {
		createColumn(b, col)
		if i != len(t.Columns)-1 {
			b.Comma()
		}
	}
	if len(t.PrimaryKeys) > 0 {
		b.Comma()
		b.Append("PRIMARY KEY")
		b.Open()
		quoteSlice(b, t.PrimaryKeys)
		b.Close()
	}
	if len(t.ForeignKeys) > 0 {
		b.Comma()
		for i, fk := range t.ForeignKeys {
			foreignKey(b, fk)
			if i != len(t.ForeignKeys)-1 {
				b.Comma()
			}
		}
	}
}

// createColumn is an SQL column definition, as given in the SQL CREATE statement
func createColumn(b *builder, col migration.Column) {
	b.Append(b.Quote(col.Name))
	b.Append(b.GoTypeToColumnType(col.GoType, hasAutoIncrement(col.Options)).String())
	for _, opt := range col.Options {
		b.Append(b.Translate(opt))
	}
}

// foreignKey is teh FOREIGN KEY statement
func foreignKey(b *builder, fk migration.ForeignKey) {
	b.Append("FOREIGN KEY")
	b.Open()
	quoteSlice(b, fk.Columns)
	b.Close()
	b.Append("REFERENCES")
	b.Append(b.Quote(fk.Table))
	b.Open()
	quoteSlice(b, fk.ForeignColumns)
	b.Close()
}

func hasAutoIncrement(options []string) bool {
	for _, opt := range options {
		if opt == "AUTO_INCREMENT" {
			return true
		}
	}
	return false
}

// selectColumns returns the columns selected for an SQL SELECT query
func selectColumns(b *builder, p *SelectParams) {
	noColumn := columnsColumnRec(b, p.Table, p, true)

	if p.Count {
		if !noColumn {
			b.Comma()
		}
		b.Append("COUNT(*)")
	}
}

func columnsColumnRec(b *builder, table string, p *SelectParams, first bool) bool {
	cols := p.SelectedColumns()
	for _, col := range cols {
		if !first {
			b.Comma()
		}
		b.Append(fmt.Sprintf("%s.%s", b.Quote(table), b.Quote(col)))
		first = false
	}
	for _, join := range p.Joins {
		if !columnsColumnRec(b, join.TableName(table), &join.SelectParams, first) {
			first = false
		}
	}
	return first
}

// where returns a WHERE statement
func where(b *builder, table string, w Where, j []JoinParams) {
	whereRec(b, table, w, j, true)
}

// whereRec returns a WHERE statement for a recursive join statement
// it concat all the conditions with an AND operator
func whereRec(b *builder, table string, w Where, joins []JoinParams, first bool) {
	if w != nil {
		if first {
			b.Append("WHERE")
		} else {
			b.Append("AND")
		}
		first = false
		w.Build(table, b)
	}
	for _, join := range joins {
		whereRec(b, join.TableName(table), join.SelectParams.Where, join.SelectParams.Joins, first)
	}
}

// groupBy formats an SQL GROUP BY statement
func groupBy(b *builder, table string, p *SelectParams) {
	groupByRec(b, table, p, true)
}

func groupByRec(b *builder, table string, p *SelectParams, first bool) {
	for _, group := range p.Groups {
		if first {
			b.Append("GROUP BY")
		} else {
			b.Comma()
		}
		first = false
		b.Append(fmt.Sprintf("%s.%s", b.Quote(table), b.Quote(group.Column)))
	}
	for _, join := range p.Joins {
		groupByRec(b, join.TableName(table), &join.SelectParams, first)
	}
}

// orderBy formats an SQL ORDER BY statement
func orderBy(b *builder, table string, p *SelectParams) {
	orderByRec(b, table, p, true)
}

func orderByRec(b *builder, table string, p *SelectParams, first bool) {
	for _, order := range p.Orders {
		if first {
			b.Append("ORDER BY")
		} else {
			b.Comma()
		}
		first = false
		b.Append(fmt.Sprintf("%s.%s", b.Quote(table), b.Quote(order.Column)))
		b.Append(string(order.Dir))
	}
	for _, join := range p.Joins {
		orderByRec(b, join.TableName(table), &join.SelectParams, first)
	}
}

// page formats an SQL LIMIT...OFFSET statement
func page(b *builder, p Page) {
	if p.Limit == 0 { // why would someone ask for a page of zero size?
		return
	}
	b.Append(fmt.Sprintf("LIMIT %d", p.Limit))
	if p.Offset != 0 {
		b.Append(fmt.Sprintf("OFFSET %d", p.Offset))
	}
}

// join extract SQL join list statement
func join(b *builder, p *SelectParams) {
	joinRec(b, p.Table, p)
}

func joinRec(b *builder, table string, p *SelectParams) {
	joins := p.Joins
	if len(joins) == 0 {
		return
	}
	for _, j := range joins {
		b.Append("LEFT OUTER JOIN")
		joinTable := j.TableName(table)
		b.Append(b.Quote(j.Table))
		b.Append("AS")
		b.Append(b.Quote(joinTable))
		b.Append("ON")
		b.Open()

		for i, pairing := range j.Pairings {
			if i > 0 {
				b.Append("AND")
			}
			b.Append(fmt.Sprintf("%s.%s", b.Quote(table), b.Quote(pairing.Column)))
			b.Append("=")
			b.Append(fmt.Sprintf("%s.%s", b.Quote(joinTable), b.Quote(pairing.JoinedColumn)))
		}
		b.Close()
		joinRec(b, j.TableName(table), &j.SelectParams)
	}
}

func quoteSlice(b *builder, s []string) {
	for i := range s {
		b.Append(b.Quote(s[i]))
		if i != len(s)-1 {
			b.Comma()
		}
	}
}
