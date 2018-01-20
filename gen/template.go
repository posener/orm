package gen

import (
	"fmt"
	"strings"
	"text/template"
)

var tpl = template.Must(template.New("").
	Funcs(template.FuncMap{
		"inc":          func(x int) int { return x + 1 },
		"dec":          func(x int) int { return x - 1 },
		"backtick":     func(s string) string { return fmt.Sprintf("`%s`", s) },
		"repeat":       strings.Repeat,
		"uncapitalize": func(s string) string { return strings.ToLower(string(s[0])) + s[1:] },
	}).Parse(`
{{ $name := $.Graph.Type.Name -}}
{{ $type := $.Graph.Type.Naked.Ext $.Package -}}
{{ $hasOneToManyRelation := $.Graph.Type.HasOneToManyRelation -}}
{{ $apiName := (print $name "ORM") -}}
{{ $conn := (print $.Private "Conn") -}}
{{ $countStruct := (print $name "Count") -}}
{{ $pks := $.Graph.Type.PrimaryKeys -}}

// Code generated by github.com/posener/orm; DO NOT EDIT
//
// ORM functions for type {{$name}}

package {{$.Package}}
import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/posener/orm"
	"github.com/posener/orm/dialect"
	{{ range $_, $import := $.Graph.Type.Imports -}}
	"{{$import}}"
	{{ end -}}
)

// {{$.Private}}Table is SQL table name
var {{$.Private}}Table = "{{$.Graph.Type.Table}}"

// {{$.Private}}TableProperties is a json representation of the table properties
// used for table creation and migration.
const {{$.Private}}TableProperties = {{backtick $.Table.Marshal}}

var {{$.Private}}RelationTablesProperties = map[string]string{
	{{ range $name, $table := $.RelationTables -}}
	"{{$name}}": {{backtick $table.Marshal}},
	{{ end -}}
}

// {{$.Private}}Column is for table column names
type {{$.Private}}Column string

const (
	{{ range $_, $f := $.Graph.Type.NonReferences -}}
	// {{$.Public}}Col{{$f.Name}} is used to select the {{$f.Name}} column in SELECT queries
	{{$.Public}}Col{{$f.Name}} {{$.Private}}Column = "{{$f.Column.Name}}"
	{{ end -}}
)

// {{$.Private}}OrderedColumns is an ordered list of all the columns in the table
var {{$.Private}}OrderedColumns = []string{
	{{ range $_, $f := $.Graph.Type.NonReferences -}}
	string({{$.Public}}Col{{$f.Name}}),
	{{ end -}}
}

func init() {
	var v interface{} = &{{$type}}{}

	// override tableName if the type implement the TableNamer interface
	if namer, ok := v.(dialect.TableNamer); ok {
		{{$.Private}}Table = namer.TableName()
	}
}

type {{$.Public}}API interface {
	// Select returns a builder for selecting rows from an SQL table
	Select(...{{$.Private}}Column) *{{$.Public}}SelectBuilder
	// Insert returns a builder for inserting a row to an SQL table
	Insert() *{{$.Public}}InsertBuilder
	// Update returns a builder for updating a row in an SQL table
	Update() *{{$.Public}}UpdateBuilder
	// Delete returns a builder for deleting a row in an SQL table
	Delete() *{{$.Public}}DeleteBuilder
	// Where returns a builder to build a where statement to be used in a Where function
	Where() *{{$.Public}}WhereBuilder
	{{ if $pks -}}
	// Get returns a builder for getting an object by primary key
	// It is a shortcut for writing a select with where statement and get the first.
	// In case that the object was not found, it returns an error orm.ErrNotFound
	Get({{range $_, $pk := $pks}}{{uncapitalize $pk.Name}} {{$pk.Type.Ext $.Package}},{{end}}) *{{$.Public}}GetBuilder
	{{ end -}}

	{{ range $_, $edge := $.Graph.RelTable -}}
	// {{$edge.Field.RelationName}} returns a relation handler for {{$edge.Field.Type.Name}}
	{{$edge.Field.RelationName}}() *{{$.Public}}{{$edge.Field.RelationName}}Builder
	{{ end -}}
}

// {{$apiName}} is the interface of the ORM object
type {{$apiName}} interface {
	{{$.Public}}API
	// Begin begins an SQL transaction and returns the transaction ORM object
	Begin(context.Context, *sql.TxOptions) ({{$apiName}}Tx, error)
	// Create returns a builder for creating an SQL table
	Create() *{{$.Public}}CreateBuilder
	// Drop returns a builder for dropping an SQL table
	Drop() *{{$.Public}}DropBuilder
}

// {{$apiName}} is the interface of the ORM object
type {{$apiName}}Tx interface {
	{{$.Public}}API
	Commit() error
	Rollback() error
}

// New{{$apiName}} returns an conn object from a db instance
func New{{$apiName}}(conn orm.Conn) ({{$apiName}}, error) {
	d := dialect.Get(conn.Driver())
	if d == nil {
		return nil, fmt.Errorf("dialect %s does not exists", conn.Driver())
	}
	return &{{$conn}}{
		Conn:    conn,
		dialect: d,
	}, nil
}

// {{$conn}} represents a DB connection for manipulating a given struct.
// All functions available to interact with an SQL table that is related
// to this struct, are done by an instance of this object.
// To get an instance of orm use Open or New functions.
type {{$conn}} struct {
	orm.Conn
	dialect dialect.API
}

func (c *{{$conn}}) Begin(ctx context.Context, opt *sql.TxOptions) ({{$apiName}}Tx, error) {
	tx, err := c.Conn.Begin(ctx, opt)
	if err != nil {
		return nil, err
	}
	retConn := *c
	retConn.Conn = tx
	return &retConn, nil
}

// Create returns a builder of an SQL CREATE statement
func (c *{{$conn}}) Create() *{{$.Public}}CreateBuilder {
	return &{{$.Public}}CreateBuilder{
		params: dialect.CreateParams{
			Table: {{$.Private}}Table,
			MarshaledTable: {{$.Private}}TableProperties,
			MarshaledRelationTables: {{$.Private}}RelationTablesProperties,
		},
		conn: c,
	}
}

// Select returns a builder of an SQL SELECT statement
func (c *{{$conn}}) Select(cols ...{{$.Private}}Column) *{{$.Public}}SelectBuilder {
	s := &{{$.Public}}SelectBuilder{
		params: dialect.SelectParams{
			Table: {{$.Private}}Table,
			OrderedColumns: {{$.Private}}OrderedColumns,
		},
		conn: c,
	}
	s.params.Columns = make(map[string]bool, len(cols))
	for _, col := range cols {
		s.params.Columns[string(col)] = true
	}
	return s
}

// Insert returns a builder of an SQL INSERT statement
func (c *{{$conn}}) Insert() *{{$.Public}}InsertBuilder {
	return &{{$.Public}}InsertBuilder{
		params: dialect.InsertParams{Table: {{$.Private}}Table},
		conn: c,
	}
}

// Update returns a builder of an SQL UPDATE statement
func (c *{{$conn}}) Update() *{{$.Public}}UpdateBuilder {
	return &{{$.Public}}UpdateBuilder{
		params: dialect.UpdateParams{Table: {{$.Private}}Table},
		conn: c,
	}
}

// Delete returns a builder of an SQL DELETE statement
func (c *{{$conn}}) Delete() *{{$.Public}}DeleteBuilder {
	return &{{$.Public}}DeleteBuilder{
		params: dialect.DeleteParams{Table: {{$.Private}}Table},
		conn: c,
	}
}

// Where returns a builder of an SQL WHERE statement
func (c *{{$conn}}) Drop() *{{$.Public}}DropBuilder {
	return &{{$.Public}}DropBuilder{
		params: dialect.DropParams{
			Table: {{$.Private}}Table,
		},
		conn: c,
	}
}

// Drop returns a builder of an SQL DROP statement
func (c *{{$conn}}) Where() *{{$.Public}}WhereBuilder {
	return &{{$.Public}}WhereBuilder{}
}

// === CreateBuilder ===

// {{$.Public}}CreateBuilder builds an SQL CREATE statement parameters
type {{$.Public}}CreateBuilder struct {
	params dialect.CreateParams
	conn   *{{$conn}}
}

// IfNotExists sets IF NOT EXISTS for the CREATE SQL statement
func (b *{{$.Public}}CreateBuilder) IfNotExists() *{{$.Public}}CreateBuilder {
	b.params.IfNotExists = true
	return b
}

// AutoMigrate sets auto-migration mode for table creation
func (b *{{$.Public}}CreateBuilder) AutoMigrate() *{{$.Public}}CreateBuilder {
	b.params.AutoMigrate = true
	return b
}

{{ if $.Graph.RelTable }}
// Relations makes Exec create relation tables instead of the type table
func (b *{{$.Public}}CreateBuilder) Relations() *{{$.Public}}CreateBuilder {
	b.params.Relations = true
	return b
}
{{ end -}}

// Context sets the context for the SQL query
func (b *{{$.Public}}CreateBuilder) Context(ctx context.Context) *{{$.Public}}CreateBuilder {
	b.params.Ctx = ctx
	return b
}

// === InsertBuilder ===

// {{$.Public}}InsertBuilder builds an INSERT statement parameters
type {{$.Public}}InsertBuilder struct {
	params dialect.InsertParams
	conn   *{{$conn}}
}

// Context sets the context for the SQL query
func (b *{{$.Public}}InsertBuilder) Context(ctx context.Context) *{{$.Public}}InsertBuilder {
	b.params.Ctx = ctx
	return b
}

// Insert{{$name}} returns an SQL INSERT statement builder filled with values of a given object
func (b *{{$.Public}}InsertBuilder) Insert{{$name}}(p *{{$type}}) *{{$.Public}}InsertBuilder {
	{{ range $_, $f := $.Graph.Type.Fields -}}
	{{ if $f.IsSettable -}}
	{{ if not $f.IsReference -}}
	b.params.Assignments.Add("{{$f.Column.Name}}", p.{{$f.AccessName}}, p.{{$f.AccessName}})
	{{ else -}}
	{{ if $f.Type.Pointer -}}
	if p.{{$f.Name}} != nil {
	{{ end -}}
	{{ range $i, $col := $f.Columns -}}
	b.params.Assignments.Add("{{$col.Name}}", p.{{$f.AccessName}}.{{(index $f.Type.PrimaryKeys $i).Name}}, p.{{$f.AccessName}})
	{{ end -}}
	{{ if $f.Type.Pointer -}}
	}
	{{ end -}}
	{{ end -}}
	{{ end -}}
	{{ end -}}
	return b
}

{{ range $_, $f := $.Graph.Type.Fields -}}
{{ if $f.IsSettable -}}
// Set{{$f.Name}} sets value for column in the INSERT statement
func (b *{{$.Public}}InsertBuilder) Set{{$f.Name}}(value {{$f.Type.Ext $.Package}}) *{{$.Public}}InsertBuilder {
	{{ if $f.IsReference -}}
	{{ range $i, $col := $f.Columns -}}
	b.params.Assignments.Add("{{$col.Name}}", value.{{(index $f.Type.PrimaryKeys $i).Name}}, value)
	{{ end -}}
	{{ else -}}
	b.params.Assignments.Add("{{$f.Column.Name}}", value, value)
	{{ end -}}
	return b
}
{{ end -}}
{{ end -}}

// === UpdateBuilder ===

// {{$.Public}}UpdateBuilder builds SQL INSERT statement parameters
type {{$.Public}}UpdateBuilder struct {
	params dialect.UpdateParams
	conn   *{{$conn}}
}

// Where sets the WHERE statement to the SQL query
func (b *{{$.Public}}UpdateBuilder) Where(where dialect.Where) *{{$.Public}}UpdateBuilder {
	b.params.Where = where
	return b
}

// Context sets the context for the SQL query
func (b *{{$.Public}}UpdateBuilder) Context(ctx context.Context) *{{$.Public}}UpdateBuilder {
	b.params.Ctx = ctx
	return b
}

// Update{{$name}} update values for all struct fields
func (b *{{$.Public}}UpdateBuilder) Update{{$name}}(p *{{$type}}) *{{$.Public}}UpdateBuilder {
	{{ range $_, $f := $.Graph.Type.Fields -}}
	{{ if $f.IsSettable -}}
	{{ if not $f.IsReference -}}
	b.params.Assignments.Add("{{$f.Column.Name}}", p.{{$f.AccessName}}, p.{{$f.AccessName}})
	{{ else -}}
	{{ if $f.Type.Pointer -}}
	if p.{{$f.Name}} != nil {
	{{ end -}}
	{{ range $i, $col := $f.Columns -}}
	b.params.Assignments.Add("{{$col.Name}}", p.{{$f.AccessName}}.{{(index $f.Type.PrimaryKeys $i).Name}}, p.{{$f.AccessName}})
	{{ end -}}
	{{ if $f.Type.Pointer -}}
	}
	{{ end -}}
	{{ end -}}
	{{ end -}}
	{{ end -}}
	return b
}

{{ range $_, $f := $.Graph.Type.Fields -}}
{{ if $f.IsSettable -}}
// Set{{$f.Name}} sets value for column in the UPDATE statement
func (b *{{$.Public}}UpdateBuilder) Set{{$f.Name}}(value {{$f.Type.Ext $.Package}}) *{{$.Public}}UpdateBuilder {
	{{ if $f.IsReference -}}
	{{ range $i, $col := $f.Columns -}}
	b.params.Assignments.Add("{{$col.Name}}", value.{{(index $f.Type.PrimaryKeys $i).Name}}, value)
	{{ end -}}
	{{ else -}}
	b.params.Assignments.Add("{{$f.Column.Name}}", value, value)
	{{ end -}}
	return b
}
{{ end -}}
{{ end -}}

// === DeleteBuilder ===

// {{$.Public}}DeleteBuilder builds SQL DELETE statement parameters
type {{$.Public}}DeleteBuilder struct {
	params dialect.DeleteParams
	conn   *{{$conn}}
}

// Where applies where conditions on the SQL query
func (b *{{$.Public}}DeleteBuilder) Where(w dialect.Where) *{{$.Public}}DeleteBuilder {
	b.params.Where = w
	return b
}

// Context sets the context for the SQL query
func (b *{{$.Public}}DeleteBuilder) Context(ctx context.Context) *{{$.Public}}DeleteBuilder {
	b.params.Ctx = ctx
	return b
}

// === DropBuilder ===

// {{$.Public}}DropBuilder builds an SQL DROP statement parameters
type {{$.Public}}DropBuilder struct {
	params dialect.DropParams
	conn   *{{$conn}}
}

// IfExists sets IF NOT EXISTS for the CREATE SQL statement
func (b *{{$.Public}}DropBuilder) IfExists() *{{$.Public}}DropBuilder {
	b.params.IfExists = true
	return b
}

// Context sets the context for the SQL query
func (b *{{$.Public}}DropBuilder) Context(ctx context.Context) *{{$.Public}}DropBuilder {
	b.params.Ctx = ctx
	return b
}

// Exec creates a table for the given struct
func (b *{{$.Public}}CreateBuilder) Exec() error {
	b.params.Ctx = dialect.ContextOrBackground(b.params.Ctx)
	stmts, err := b.conn.dialect.Create(b.conn.Conn, &b.params)
	if err != nil {
		return err
	}
	if len(stmts) == 0 {
		return nil
	}
	tx, err := b.conn.Conn.Begin(b.params.Ctx, nil)
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		_, err := tx.Exec(b.params.Ctx, stmt)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// query is used by the Select.Query and Select.Limit functions
func (b *{{$.Public}}SelectBuilder) query() (*sql.Rows, error) {
	stmt, args := b.conn.dialect.Select(&b.params)
	return b.conn.Query(b.params.Ctx, stmt, args...)
}

// Exec inserts the data to the given database
{{ $getID := (eq (len $pks) 1) -}}
func (b *{{$.Public}}InsertBuilder) Exec() (*{{$type}}, error) {
	if len(b.params.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to insert")
	}
	b.params.Ctx = dialect.ContextOrBackground(b.params.Ctx)

	{{ if $pks -}}
	if b.conn.dialect.Name() == dialect.Postgres {
		return b.execPostgres()
	}
	{{ end -}}

	stmt, args := b.conn.dialect.Insert(&b.params)
	ret := {{$.Private}}ReturnObject(b.params.Assignments)
	{{ if $getID -}}res{{ else }}_{{ end }}, err := b.conn.Exec(b.params.Ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	{{ if $getID -}}
	if id, err := res.LastInsertId(); err == nil {
		ret.{{(index $pks 0).AccessName}} = {{(index $pks 0).Type.Ext $.Package}}(id)
	}
	{{ end -}}
	return ret, nil
}

{{ if $pks -}}
// execProgress execute insert statement for postgres engine
// Postgres is different because it does not support LastInsertId
func (b *{{$.Public}}InsertBuilder) execPostgres() (*{{$type}}, error) {
	b.params.RetColumns = []string{
		{{- range $_, $pk := $pks -}}
		"{{- $pk.Column.Name -}}",
		{{- end -}}
	}
	stmt, args := b.conn.dialect.Insert(&b.params)
	ret := {{$.Private}}ReturnObject(b.params.Assignments)
	err := b.conn.QueryRow(b.params.Ctx, stmt, args...).Scan({{- range $_, $pk := $pks -}}&ret.{{- $pk.AccessName -}}, {{- end -}})
	return ret, err
}
{{ end -}}

// Exec inserts the data to the given database
func (b *{{$.Public}}UpdateBuilder) Exec() (sql.Result, error) {
	if len(b.params.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	b.params.Ctx = dialect.ContextOrBackground(b.params.Ctx)
	stmt, args := b.conn.dialect.Update(&b.params)
	return b.conn.Exec(b.params.Ctx, stmt, args...)
}

// Exec runs the delete statement on a given database.
func (b *{{$.Public}}DeleteBuilder) Exec() (sql.Result, error) {
	stmt, args := b.conn.dialect.Delete(&b.params)
	return b.conn.Exec(dialect.ContextOrBackground(b.params.Ctx), stmt, args...)
}

// Query the database
func (b *{{$.Public}}SelectBuilder) Query() ([]*{{$type}}, error) {
	b.params.Ctx = dialect.ContextOrBackground(b.params.Ctx)
	rows, err := b.query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		items []*{{$type}}
		// exists is a mapping from primary key to already parsed structs
		exists = make(map[string]interface{})
	)
	for rows.Next() {
		// check context cancellation
		if err := b.params.Ctx.Err(); err != nil  {
			return nil, err
		}
		item, _, err := b.scan(b.conn.dialect.Name(), dialect.Values(*rows), exists, "")
		if err != nil {
			return nil, err
		}
		if item != nil {
			items = append(items, item)
		}
	}
	return items, rows.Err()
}

// Count add a count column to the query
func (b *{{$.Public}}SelectBuilder) Count() ([]{{$countStruct}}, error) {
	b.params.Ctx = dialect.ContextOrBackground(b.params.Ctx)
	b.params.Count = true
	rows, err := b.query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		items []{{$countStruct}}
		// exists is a mapping from primary key to already parsed structs
		exists = make(map[string]interface{})
	)
	for rows.Next() {
		// check context cancellation
		if err := b.params.Ctx.Err(); err != nil  {
			return nil, err
		}
		item, _, err := b.scanCount(b.conn.dialect.Name(), dialect.Values(*rows), exists, "")
		if err != nil {
			return nil, err
		}
		if item != nil {
			items = append(items, *item)
		}
	}
	return items, rows.Err()
}

// First returns the first row that matches the query.
// If no row matches the query, an ErrNotFound will be returned.
// This call cancels any paging that was set with the
// {{$.Public}}SelectBuilder previously.
func (b *{{$.Public}}SelectBuilder) First() (*{{$type}}, error) {
	b.params.Ctx = dialect.ContextOrBackground(b.params.Ctx)
	b.params.Page.Limit = 1
	b.params.Page.Offset = 0
	rows, err := b.query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := rows.Next()
	if !found {
		return nil, orm.ErrNotFound
	}
	var exists = make(map[string]interface{})
	item, _, err := b.scan(b.conn.dialect.Name(), dialect.Values(*rows), exists, "")
	if err != nil {
		return nil, err
	}
	return item, rows.Err()
}

// {{$.Private}}ReturnObject builds {{$type}} from assignment values
// and from the sql result ID, for returning an object from INSERT transactions
func {{$.Private}}ReturnObject(assignments dialect.Assignments) *{{$type}} {
	ret := new({{$type}})
	for _, assign := range assignments {
		switch assign.Column {
		{{ range $_, $f := $.Graph.Type.Fields -}}
		{{ if not $f.Type.Slice -}}
		case "{{(index $f.Columns 0).Name}}":
			ret.{{$f.AccessName}} = assign.OriginalValue.({{$f.Type.Ext $.Package}})
		{{ end -}}
		{{ end -}}
		}
	}
	return ret
}

// Exec runs the drop statement on a given database.
func (b *{{$.Public}}DropBuilder) Exec() error {
	stmt, args := b.conn.dialect.Drop(&b.params)
	_, err := b.conn.Exec(dialect.ContextOrBackground(b.params.Ctx), stmt, args...)
	return err
}

// {{$.Graph.Type.Naked.Name}}Joiner is an interface for joining a {{$name}} in a SELECT statement
// in another type
type {{$name}}Joiner interface {
	Params() dialect.SelectParams
	Scan(dialect string, values []driver.Value, exists map[string]interface{}, prefix string) (*{{$type}}, int, error)
}

// {{$countStruct}} is a struct for counting rows of type {{$name}}
type {{$countStruct}} struct {
	*{{$.Graph.Type.Ext $.Package}}
	Count int64
}

// {{$.Public}}SelectBuilder builds an SQL SELECT statement parameters
type {{$.Public}}SelectBuilder struct {
	params dialect.SelectParams
	conn *{{$conn}}
	{{ range $_, $f := $.Graph.Type.References -}}
	scan{{$f.Name}} {{$.Private}}{{$f.Type.Naked.Name}}Joiner
	{{ end -}}
}

// {{$.Private}}Joiner represents a builder that exposes only Params and Scan method
type {{$.Private}}Joiner struct {
	builder *{{$.Public}}SelectBuilder
}

func (j *{{$.Private}}Joiner) Params() dialect.SelectParams {
	return j.builder.params
}

func (j *{{$.Private}}Joiner) Scan(dialect string, values []driver.Value, exists map[string]interface{}, prefix string) (*{{$type}}, int, error) {
	return j.builder.scan(dialect, values, exists, prefix)
}

// Joiner returns an object to be used in a join operation with {{$name}}
func (b *{{$.Public}}SelectBuilder) Joiner() {{$.Public}}Joiner {
	return &{{$.Private}}Joiner{builder: b}
}

// Where applies where conditions on the query
func (b *{{$.Public}}SelectBuilder) Where(where dialect.Where) *{{$.Public}}SelectBuilder {
	b.params.Where = where
	return b
}

// Limit applies rows limit on the query response
func (b *{{$.Public}}SelectBuilder) Limit(limit int64) *{{$.Public}}SelectBuilder {
	b.params.Page.Limit = limit
	return b
}

// Page applies rows offset and limit on the query response
func (b *{{$.Public}}SelectBuilder) Page(offset, limit int64) *{{$.Public}}SelectBuilder {
	b.params.Page.Offset = offset
	b.params.Page.Limit = limit
	return b
}

{{ range $_, $refType := $.Graph.Type.ReferencedTypes -}}
// {{$.Private}}{{$refType.Name}}Joiner is a scanner that defined by .Select().Joiner()
// of an ORM object for type {{$refType.Name}}
type {{$.Private}}{{$refType.Name}}Joiner interface {
	Params() dialect.SelectParams
	Scan(dialect string, values []driver.Value, exists map[string]interface{}, prefix string) (*{{$refType.Ext $.Package}}, int, error)
}
{{ end -}}

{{ range $_, $e := $.Graph.Out -}}
{{ $f := $e.Field -}}
// Join{{$f.Name}} add a join query for {{$f.Name}}
// Based on a forward relation
func (b *{{$.Public}}SelectBuilder) Join{{$f.Name}}(joiner {{$.Private}}{{$f.Type.Name}}Joiner) *{{$.Public}}SelectBuilder {
	b.scan{{$f.Name}} = joiner
	b.params.Joins = append(b.params.Joins, dialect.JoinParams{
		Pairings: []dialect.Pairing{
			{{ range $i, $pk := $e.RelationType.PrimaryKeys -}}
			{
				Column: "{{(index $e.Field.Columns $i).Name}}",
				JoinedColumn: "{{$pk.Column.Name}}",
			},
			{{ end -}}
		},
		SelectParams: joiner.Params(),
	})
	return b
}
{{ end -}}

{{ range $_, $e := $.Graph.In -}}
{{ $f := $e.Field -}}
// Join{{$f.Name}} add a join query for {{$f.Name}}
// Based on a reversed relation
func (b *{{$.Public}}SelectBuilder) Join{{$f.Name}}(joiner {{$.Private}}{{$f.Type.Name}}Joiner) *{{$.Public}}SelectBuilder {
	b.scan{{$f.Name}} = joiner
	b.params.Joins = append(b.params.Joins,
		dialect.JoinParams{
			Pairings: []dialect.Pairing{
				{{ range $i, $pk := $e.RelationType.PrimaryKeys -}}
				{
					Column: "{{$pk.Column.Name}}",
					JoinedColumn: "{{(index $e.SrcField.Columns $i).Name}}",
				},
				{{ end -}}
			},
			SelectParams: joiner.Params(),
		},
	)
	return b
}
{{ end -}}

{{ range $_, $e := $.Graph.RelTable -}}
{{ $f := $e.Field -}}
// Join{{$f.Name}} add a join query for {{$f.Name}}
// Based on a relation table
func (b *{{$.Public}}SelectBuilder) Join{{$f.Name}}(joiner {{$.Private}}{{$f.Type.Name}}Joiner) *{{$.Public}}SelectBuilder {
	b.scan{{$f.Name}} = joiner
	b.params.Joins = append(b.params.Joins, 
		dialect.JoinParams{
			Pairings: []dialect.Pairing{
				{{ range $i, $pk := $f.ParentType.PrimaryKeys -}}
				{
					Column: "{{$pk.Column.Name}}",
					JoinedColumn: "{{$f.ParentType.Table}}_{{$pk.Column.Name}}",
				},
				{{ end -}}
			},
			SelectParams: dialect.SelectParams{
				Table: "{{$e.Field.RelationTable}}",
				Joins: []dialect.JoinParams{
					{
						Pairings: []dialect.Pairing{
							{{ range $i, $pk := $f.Type.PrimaryKeys -}}
							{
								Column: "{{$f.Type.Table}}_{{$pk.Column.Name}}",
								JoinedColumn: "{{$pk.Column.Name}}",
							},
							{{ end -}}
						},
						SelectParams: joiner.Params(),
					},
				},
			},
		},
	)
	return b
}
{{ end -}}

// OrderBy set order to the query results according to column
func (b *{{$.Public}}SelectBuilder) OrderBy(col {{$.Private}}Column, dir orm.OrderDir) *{{$.Public}}SelectBuilder {
	b.params.Orders.Add(string(col), dir)
	return b
}

// GroupBy make the query group by column
func (b *{{$.Public}}SelectBuilder) GroupBy(col {{$.Private}}Column) *{{$.Public}}SelectBuilder {
	b.params.Groups.Add(string(col))
	return b
}

// Context sets the context for the SQL query
func (b *{{$.Public}}SelectBuilder) Context(ctx context.Context) *{{$.Public}}SelectBuilder {
	b.params.Ctx = ctx
	return b
}

// scan an SQL row to a {{$name}} struct
// It returns the scanned {{$.Graph.Type.Ext $.Package}} and the number of scanned fields,
// and an error in case of failure.
func (b *{{$.Public}}SelectBuilder) scan(dialect string, vals []driver.Value, exists map[string]interface{}, prefix string) (*{{$.Graph.Type.Ext $.Package}}, int, error) {
	item, n, err := b.scanCount(dialect, vals, exists, prefix)
	if err != nil {
		return nil, n, err
	}
	return item.{{$name}}, n, nil
}

// ScanCount scans an SQL row to a {{$countStruct}} struct
func (b *{{$.Public}}SelectBuilder) scanCount(dialect string, vals []driver.Value, exists map[string]interface{}, prefix string) (*{{$countStruct}}, int, error) {
	switch dialect {
	{{ range $_, $dialect := $.Dialects -}}
	case "{{$dialect.Name}}":
		return b.scan{{$dialect.Name}}(vals, exists, prefix)
	{{ end -}}
	default:
		return nil, 0, fmt.Errorf("unsupported dialect %s", dialect)
	}
}


{{ if len $pks -}}
func {{$.Private}}HashItem(item *{{$name}}) string {
	return fmt.Sprintf("{{$.Private}}/{{repeat "%v" (len $pks)}}",
	{{- range $f := $pks -}}
	{{if $f.Type.Pointer}}*{{end}}item.{{$f.AccessName}},
	{{- end -}}
	)
}
{{ end -}}

type {{$.Private}}Exists struct{
	self *{{$type}}
	related map[string]interface{}
}

{{ range $_, $dialect := $.Dialects }}
// scan{{$dialect.Name}} scans {{$dialect.Name}} row to a {{$name}} struct
func (b *{{$.Public}}SelectBuilder) scan{{$dialect.Name}} (vals []driver.Value, exists map[string]interface{}, prefix string) (*{{$countStruct}}, int, error) {
	var (
		row = new({{$countStruct}})
		i int
		skipRow = false
		allNils = true
		all = b.params.SelectAll()
		{{ if $.Graph.Type.References -}}
		existingRow = &{{$.Private}}Exists{related: make(map[string]interface{})}
		{{ end -}}
	)
	row.{{$name}} = new({{$type}})
	{{ range $i, $f := $.Graph.Type.NonReferences -}}
	// scan column {{$i}}
	if all || b.params.Columns["{{$f.Column.Name}}"] {
		if i >= len(vals) {
			return nil, 0, fmt.Errorf("not enough columns returned: %d", len(vals))
		}
		if vals[i] != nil && !skipRow {
			allNils = false
{{ $dialect.ConvertValueCode $f -}}
		}
		{{ if or $f.Unique $f.PrimaryKey -}}
		{{/* Test if the field is the last primary key */}}
		{{ if eq $f.Name (index $pks (dec (len $pks))).Name -}}
		{{ if $hasOneToManyRelation -}}
		// check if we scanned this item in previous rows. If we did, set existingRow
		// so other columns in this table won't be evaluated. And joined items could
		// be joined to the already scanned item.
		hash := prefix + {{$.Private}}HashItem(row.{{$name}}) 
		if exists[hash] == nil {
			exists[hash] = &{{$.Private}}Exists{self: row.{{$name}}, related: make(map[string]interface{})}
		} else {
			skipRow = true	
		}
		existingRow = exists[hash].(*{{$.Private}}Exists)
		{{ end -}}
		{{ end -}}
		{{ end -}}
		i++
	}
	{{ end -}}

	if b.params.Count {
		switch val := vals[i].(type) {
		case int64:
			row.Count = val
		case []byte:
			row.Count = dialect.ParseInt(val)
		default:
			return nil, 0, dialect.ErrConvert("COUNT(*)", i, vals[i], "int64, []byte")
		}
		i++
	}

	{{ range $_, $f := $.Graph.Type.References -}}
	if b := b.scan{{$f.Name}}; b != nil {
		tmp, n, err := b.Scan("{{$dialect.Name}}", vals[i:], existingRow.related, "{{$f.Name}}")
		if err != nil {
			return nil, 0, fmt.Errorf("sub scanning {{$f.AccessName}}, cols [%d:%d]: %s", i, len(vals), err)
		}
		// If the result is nil, we want to discard it.
		// This is possible since we are doing a left join, if there was no match in the
		// right table, all it's columns are set to nil, and the result of a Scan function
		// is nil also.
		if tmp != nil {
			{{ if $f.Type.Slice -}}
			existingRow.self.{{$f.AccessName}} = append(existingRow.self.{{$f.AccessName}}, {{if not $f.Type.Pointer}}*{{end}}tmp)
			{{ else -}}
			row.{{$f.AccessName}} = {{ if not $f.Type.Pointer}}*{{end}}tmp
			{{ end -}}
		}
		i += n
	}
	{{ end }}

	// If all values were nil, there was not any actual row returned,
	// this could happen in case that the scanned row is the right table in case of an
	// outer left join statement. We set the result to nil, so it ill be discarded.
	if allNils || skipRow {
		row.{{$name}} = nil
	}

	return row, i, nil
}
{{ end }}

type {{$.Public}}WhereBuilder struct {}

{{ range $_, $f := $.Graph.Type.NonReferences -}}
// {{$.Public}}Where{{$f.Name}} adds a condition on {{$f.Name}} to the WHERE statement
func (*{{$.Public}}WhereBuilder) {{$f.Name}}(op orm.Op, val {{$f.Type.Ext $.Package}}) dialect.Where {
	return dialect.NewWhere(op, "{{$f.Column.Name}}", val)
}

// {{$.Public}}Where{{$f.Name}}In adds an IN condition on {{$f.Name}} to the WHERE statement
func (*{{$.Public}}WhereBuilder) {{$f.Name}}In(vals ...{{$f.Type.Ext $.Package}}) dialect.Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return dialect.NewWhereIn("{{$f.Column.Name}}", args...)
}

// {{$.Public}}Where{{$f.Name}}Between adds a BETWEEN condition on {{$f.Name}} to the WHERE statement
func (*{{$.Public}}WhereBuilder) {{$f.Name}}Between(low, high {{$f.Type.Ext $.Package}}) dialect.Where {
	return dialect.NewWhereBetween("{{$f.Column.Name}}", low, high)
}
{{ end -}}


{{ if $pks -}}
// === GetBuilder ===

func (c *{{$conn}}) Get({{range $i, $pk := $pks}}key{{$i}} {{$pk.Type.Ext $.Package}},{{end}}) *{{.Public}}GetBuilder {
	return &{{$.Public}}GetBuilder{
		conn: c,
		{{ range $i, $pk := $pks -}}
		key{{$i}}: key{{$i}},
		{{ end -}}
	}
}

type {{$.Public}}GetBuilder struct {
	conn *{{$conn}}
	ctx context.Context

	{{ range $i, $pk := $pks -}}
	key{{$i}} {{$pk.Type.Ext $.Package}}
	{{ end -}}
}

func (g *{{$.Public}}GetBuilder) Context(ctx context.Context) *{{$.Public}}GetBuilder {
	g.ctx = ctx
	return g
}

func (g *{{$.Public}}GetBuilder) Exec() (*{{$type}}, error) {
	return g.conn.Select().
		Context(dialect.ContextOrBackground(g.ctx)).
		Where(
		{{- range $i, $pk := $pks -}}
			g.conn.Where().{{$pk.Name}}(orm.OpEq, g.key{{$i}})
			{{- if ne (inc $i) (len $pks) -}}
			.And(
			{{- end -}}
		{{end}}
		{{- range $i, $_ := $pks -}}
		{{- if ne (inc $i) (len $pks) -}}
		)
		{{- end -}}
		{{- end -}}
		).
		First()
}
{{ end -}}


{{ range $_, $edge := $.Graph.RelTable -}}
{{ $relName := (print $.Public $edge.Field.RelationName "Builder") -}}
{{ $localPks := $edge.Field.ParentType.PrimaryKeys -}}
{{ $remotePks := $edge.Field.Type.PrimaryKeys -}}
func (c *{{$conn}}) {{$edge.Field.RelationName}}() *{{$relName}} {
	return &{{$relName}}{conn: c}
}

type {{$relName}} struct {
	conn *{{$conn}}
	ctx context.Context
}

func (r *{{$relName}}) Context(ctx context.Context) *{{$relName}} {
	r.ctx = ctx
	return r
}

func (r *{{$relName}}) Add(
	{{- range $_, $pk := $localPks -}}
	{{uncapitalize $edge.Field.ParentType.Name}}{{$pk.Name}} {{$pk.Type.Naked.Ext $.Package}},
	{{- end -}}
	{{- range $_, $pk := $remotePks -}}
	{{uncapitalize $edge.Field.Type.Name}}{{$pk.Name}} {{$pk.Type.Naked.Ext $.Package}},
	{{- end -}}
	) error {
	stmt, args := r.conn.dialect.Insert(&dialect.InsertParams{
		Table: "{{$edge.Field.RelationTable}}",
		Assignments: dialect.Assignments{
			{{ range $_, $pk := $localPks -}}
			{{ $tp := $edge.Field.ParentType -}}
			{Column: "{{$tp.Table}}_{{$pk.Column.Name}}", ColumnValue: {{uncapitalize $tp.Name}}{{$pk.Name}}},
			{{ end -}}
			{{ range $_, $pk := $remotePks -}}
			{{ $tp := $edge.Field.Type -}}
			{Column: "{{$tp.Table}}_{{$pk.Column.Name}}", ColumnValue: {{uncapitalize $tp.Name}}{{$pk.Name}}},
			{{ end -}}
		},
	})
	_, err := r.conn.Exec(dialect.ContextOrBackground(r.ctx), stmt, args...)
	return err
}

func (r *{{$relName}}) Remove(
	{{- range $_, $pk := $edge.Field.ParentType.PrimaryKeys -}}
	{{uncapitalize $edge.Field.ParentType.Name}}{{$pk.Name}} {{$pk.Type.Naked.Ext $.Package}},
	{{- end -}}
	{{- range $_, $pk := $edge.Field.Type.PrimaryKeys -}}
	{{uncapitalize $edge.Field.Type.Name}}{{$pk.Name}} {{$pk.Type.Naked.Ext $.Package}},
	{{- end -}}
	) error {
	stmt, args := r.conn.dialect.Delete(&dialect.DeleteParams{
		Table: "{{$edge.Field.RelationTable}}",
		Where: dialect.And(
			{{ range $_, $pk := $localPks -}}
			{{ $tp := $edge.Field.ParentType -}}
			dialect.NewWhere(orm.OpEq, "{{$tp.Table}}_{{$pk.Column.Name}}", {{uncapitalize $tp.Name}}{{$pk.Name}}),
			{{ end -}}
			{{ range $_, $pk := $remotePks -}}
			{{ $tp := $edge.Field.Type -}}
			dialect.NewWhere(orm.OpEq, "{{$tp.Table}}_{{$pk.Column.Name}}", {{uncapitalize $tp.Name}}{{$pk.Name}}),
			{{ end -}}
		),
	})
	_, err := r.conn.Exec(dialect.ContextOrBackground(r.ctx), stmt, args...)
	return err
}
{{ end -}}
`))
