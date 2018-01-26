package dialect

import (
	"context"
	"fmt"

	"github.com/posener/orm"
	"github.com/posener/orm/dialect/builder"
	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/postgres"
	"github.com/posener/orm/dialect/sqlite3"
	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/load"
)

// Constants that are names of the different supported dialects
const (
	Mysql    = "mysql"
	Postgres = "postgres"
	Sqlite3  = "sqlite3"
)

// API is an interface to represent an SQL dialect
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
	if p.Relations {
		return d.createRelationTables(conn, p)
	}
	return d.create(conn, p, p.Table, p.MarshaledTable)
}

func (d *dialect) create(conn orm.Conn, p *CreateParams, tableName, tableProperties string) ([]string, error) {
	table := new(Table)
	err := table.UnMarshal(tableProperties)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling table properties: %s", err)
	}

	if p.AutoMigrate {
		stmts, ok, err := d.autoMigrate(p.Ctx, conn, tableName, table)
		if err != nil {
			return nil, fmt.Errorf("automigration: %s", err)
		}
		if ok {
			return stmts, nil
		}
	}

	b := builder.New(d, "CREATE TABLE")
	if p.IfNotExists {
		b.Append("IF NOT EXISTS")
	}
	b.Quote(tableName)
	b.Open()
	buildTable(b, d, table)
	b.Close()
	return []string{b.Statement()}, nil
}

func (d *dialect) createRelationTables(conn orm.Conn, p *CreateParams) ([]string, error) {
	var stmts []string
	for name, props := range p.MarshaledRelationTables {
		stmt, err := d.create(conn, p, name, props)
		if err != nil {
			return nil, fmt.Errorf("creating table %s: %s", name, err)
		}
		stmts = append(stmts, stmt...)
	}
	return stmts, nil
}

func (d *dialect) autoMigrate(ctx context.Context, conn orm.Conn, tableName string, want *Table) ([]string, bool, error) {
	got, err := d.loadTable(ctx, conn, tableName)
	if err != nil {
		return nil, false, err
	}
	if got == nil {
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
		b := builder.New(d, "ALTER TABLE")
		b.Quote(tableName)
		b.Append("ADD COLUMN")
		buildColumn(b, d, col)
		stmts = append(stmts, b.Statement())
	}
	for _, fk := range diff.ForeignKeys {
		b := builder.New(d, "ALTER TABLE")
		b.Quote(tableName)
		b.Append("ADD")
		buildForeignKey(b, fk)
		stmts = append(stmts, b.Statement())
	}
	return stmts, true, nil
}

// Insert returns the SQL INSERT statement and arguments according to the given parameters
func (d *dialect) Insert(p *InsertParams) (string, []interface{}) {
	b := builder.New(d, "INSERT INTO")
	b.Quote(p.Table)
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
		b.Quote(col)
	}
	return b.Statement(), b.Args()
}

// Select returns the SQL SELECT statement and arguments according to the given parameters
func (d *dialect) Select(p *SelectParams) (string, []interface{}) {
	b := builder.New(d, "SELECT")
	buildSelectColumns(b, p)
	b.Append("FROM")
	b.Append(d.Quote(p.Table))
	buildJoin(b, p)
	buildWhere(b, p.Table, p.Where, p.Joins)
	buildGroupBy(b, p.Table, p)
	buildOrderBy(b, p.Table, p)
	buildPage(b, p.Page)
	return b.Statement(), b.Args()
}

// Delete returns the SQL DELETE statement and arguments according to the given parameters
func (d *dialect) Delete(p *DeleteParams) (string, []interface{}) {
	b := builder.New(d, "DELETE FROM")
	b.Append(d.Quote(p.Table))
	buildWhere(b, p.Table, p.Where, nil)
	return b.Statement(), b.Args()
}

// Update returns the SQL UPDATE statement and arguments according to the given parameters
func (d *dialect) Update(p *UpdateParams) (string, []interface{}) {
	b := builder.New(d, "UPDATE")
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
	buildWhere(b, p.Table, p.Where, nil)
	return b.Statement(), b.Args()
}

// Drop returns the SQL DROP statement and arguments according to the given parameters
func (d *dialect) Drop(p *DropParams) (string, []interface{}) {
	b := builder.New(d, "DROP TABLE")
	if p.IfExists {
		b.Append("IF EXISTS")
	}
	b.Append(d.Quote(p.Table))
	return b.Statement(), b.Args()
}

// buildTable returns all properties of SQL table, as should be given in the table CREATE statement
func buildTable(b *builder.Builder, d Dialect, t *Table) {
	for i, col := range t.Columns {
		buildColumn(b, d, col)
		if i != len(t.Columns)-1 {
			b.Comma()
		}
	}
	if len(t.PrimaryKeys) > 0 {
		b.Comma()
		b.Append("PRIMARY KEY")
		b.Open()
		buildQuoteSlice(b, t.PrimaryKeys...)
		b.Close()
	}
	if len(t.ForeignKeys) > 0 {
		b.Comma()
		for i, fk := range t.ForeignKeys {
			buildForeignKey(b, fk)
			if i != len(t.ForeignKeys)-1 {
				b.Comma()
			}
		}
	}
}

// buildColumn is an SQL column definition, as given in the SQL CREATE statement
func buildColumn(b *builder.Builder, d Dialect, col Column) {
	b.Quote(col.Name)
	b.Append(d.GoTypeToColumnType(col.GoType, hasAutoIncrement(col.Options)).String())
	for _, opt := range col.Options {
		b.Append(d.Translate(opt))
	}
}

// buildForeignKey is the FOREIGN KEY statement
func buildForeignKey(b *builder.Builder, fk ForeignKey) {
	b.Append("FOREIGN KEY")
	b.Open()
	buildQuoteSlice(b, fk.Columns...)
	b.Close()
	b.Append("REFERENCES")
	b.Quote(fk.Table)
	b.Open()
	buildQuoteSlice(b, fk.ForeignColumns...)
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

// buildSelectColumns returns the columns selected for an SQL SELECT query
func buildSelectColumns(b *builder.Builder, p *SelectParams) {
	first := true
	columnRec(b, p.Table, p, &first)

	if p.Count {
		if !first {
			b.Comma()
		}
		b.Append("COUNT(*)")
	}
}

func columnRec(b *builder.Builder, table string, p *SelectParams, first *bool) {
	cols := p.SelectedColumns()
	for _, col := range cols {
		if !*first {
			b.Comma()
		}
		b.QuoteFullName(table, col)
		*first = false
	}
	for _, join := range p.Joins {
		columnRec(b, join.TableName(table), &join.SelectParams, first)
	}
}

// buildWhere returns a WHERE statement
func buildWhere(b *builder.Builder, table string, w Where, j []JoinParams) {
	first := true
	whereRec(b, table, w, j, &first)
}

// whereRec returns a WHERE statement for a recursive join statement
// it concat all the conditions with an AND operator
func whereRec(b *builder.Builder, table string, w Where, joins []JoinParams, first *bool) {
	if w != nil {
		if *first {
			b.Append("WHERE")
		} else {
			b.Append("AND")
		}
		*first = false
		w.Build(table, b)
	}
	for _, join := range joins {
		whereRec(b, join.TableName(table), join.SelectParams.Where, join.SelectParams.Joins, first)
	}
}

// buildGroupBy formats an SQL GROUP BY statement
func buildGroupBy(b *builder.Builder, table string, p *SelectParams) {
	first := true
	groupByRec(b, table, p, &first)
}

func groupByRec(b *builder.Builder, table string, p *SelectParams, first *bool) {
	for _, group := range p.Groups {
		if *first {
			b.Append("GROUP BY")
		} else {
			b.Comma()
		}
		*first = false
		b.QuoteFullName(table, group.Column)
	}
	for _, join := range p.Joins {
		groupByRec(b, join.TableName(table), &join.SelectParams, first)
	}
}

// buildOrderBy formats an SQL ORDER BY statement
func buildOrderBy(b *builder.Builder, table string, p *SelectParams) {
	first := true
	orderByRec(b, table, p, &first)
}

func orderByRec(b *builder.Builder, table string, p *SelectParams, first *bool) {
	for _, order := range p.Orders {
		if *first {
			b.Append("ORDER BY")
		} else {
			b.Comma()
		}
		*first = false
		b.QuoteFullName(table, order.Column)
		b.Append(string(order.Dir))
	}
	for _, join := range p.Joins {
		orderByRec(b, join.TableName(table), &join.SelectParams, first)
	}
}

// buildPage formats an SQL LIMIT...OFFSET statement
func buildPage(b *builder.Builder, p Page) {
	if p.Limit == 0 { // why would someone ask for a page of zero size?
		return
	}
	b.Append(fmt.Sprintf("LIMIT %d", p.Limit))
	if p.Offset != 0 {
		b.Append(fmt.Sprintf("OFFSET %d", p.Offset))
	}
}

// buildJoin extract SQL join list statement
func buildJoin(b *builder.Builder, p *SelectParams) {
	joinRec(b, p.Table, p)
}

func joinRec(b *builder.Builder, table string, p *SelectParams) {
	joins := p.Joins
	if len(joins) == 0 {
		return
	}
	for _, j := range joins {
		b.Append("LEFT OUTER JOIN")
		joinTable := j.TableName(table)
		b.Quote(j.Table)
		b.Append("AS")
		b.Quote(joinTable)
		b.Append("ON")
		b.Open()

		for i, pairing := range j.Pairings {
			if i > 0 {
				b.Append("AND")
			}
			b.QuoteFullName(table, pairing.Column)
			b.Append("=")
			b.QuoteFullName(joinTable, pairing.JoinedColumn)
		}
		b.Close()
		joinRec(b, j.TableName(table), &j.SelectParams)
	}
}

func buildQuoteSlice(b *builder.Builder, s ...string) {
	for i := range s {
		b.Quote(s[i])
		if i != len(s)-1 {
			b.Comma()
		}
	}
}
