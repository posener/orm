package gen

import (
	"fmt"
	"strings"
	"text/template"
)

var tpl = template.Must(template.New("").
	Funcs(template.FuncMap{
		"plus1":    func(x int) int { return x + 1 },
		"backtick": func(s string) string { return fmt.Sprintf("`%s`", s) },
		"repeat":   func(s string, n int) string { return strings.Repeat(s, n) },
	}).Parse(`
{{ $name := $.Graph.Type.Name -}}
{{ $type := $.Graph.Type.Naked.Ext $.Package -}}
{{ $hasOneToManyRelation := $.Graph.Type.HasOneToManyRelation -}}
{{ $apiName := (print $name "ORM") -}}
{{ $conn := (print $.Private "Conn") -}}
{{ $countStruct := (print $name "Count") -}}

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
	"github.com/posener/orm/runtime"
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

// {{$.Private}}Column is for table column names
type {{$.Private}}Column string

const (
	{{ range $_, $f := $.Graph.Type.NonReferences -}}
	// {{$.Public}}Col{{$f.Name}} is used to select the {{$f.Name}} column in SELECT queries
	{{$.Public}}Col{{$f.Name}} {{$.Private}}Column = "{{$f.Column.Name}}"
	{{ end -}}
)

// {{$.Private}}OrderedColumns is an oredered list of all the columns in the table
var {{$.Private}}OrderedColumns = []string{
	{{ range $_, $f := $.Graph.Type.NonReferences -}}
	string({{$.Public}}Col{{$f.Name}}),
	{{ end -}}
}

func init() {
	var v interface{} = &{{$type}}{}

	// override tableName if the type implement the TableNamer interface
	if namer, ok := v.(runtime.TableNamer); ok {
		{{$.Private}}Table = namer.TableName()
	}
}

// {{$.Public}}SelectExecer executes select statements
type {{$.Public}}SelectExecer interface {
	Query() ([]{{$type}}, error)
	Count() ([]{{$countStruct}}, error)
	First() (*{{$type}}, error)
	Joiner() {{$.Public}}Joiner
}

// {{$.Public}}API is API for ORM operations
type {{$.Public}}API interface {
	// Select returns a builder for selecting rows from an SQL table
	Select(...{{$.Private}}OptSelect) {{$.Public}}SelectExecer
	// Insert returns a builder for inserting a row to an SQL table
	Insert() *{{$.Public}}InsertBuilder
	// Update returns a builder for updating a row in an SQL table
	Update() *{{$.Public}}UpdateBuilder
	// Delete returns a builder for deleting a row in an SQL table
	Delete() *{{$.Public}}DeleteBuilder
	// Where returns a builder to build a where statement to be used in a Where function
	Where() *{{$.Public}}WhereBuilder
	{{ if $.Graph.Type.PrimaryKeys -}}
	// Get returns an object by primary key
	// In case that the object was not found, it returns an error orm.ErrNotFound
	Get({{range $_, $pk := $.Graph.Type.PrimaryKeys}}{{$pk.PrivateName}} {{$pk.Type.Ext $.Package}},{{end}}) (*{{$type}}, error)
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

type {{$.Public}}Conn struct {
	{{$apiName}}
	// S returns a struct that holds different select options
	S {{$.Public}}SelectOpts
}

// New{{$apiName}} returns an conn object from a db instance
func New{{$apiName}}(conn orm.Conn) (*{{$.Public}}Conn, error) {
	d := dialect.Get(conn.Driver())
	if d == nil {
		return nil, fmt.Errorf("dialect %s does not exists", conn.Driver())
	}
	return &{{$.Public}}Conn{
		{{$apiName}}: &{{$conn}}{
			Conn:    conn,
			dialect: d,
		},
		S: {{$.Private}}SelectOpts,
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

func (c *{{$conn}}) Create() *{{$.Public}}CreateBuilder {
	return &{{$.Public}}CreateBuilder{
		params: runtime.CreateParams{
			Table: {{$.Private}}Table,
			MarshaledTable: {{$.Private}}TableProperties,
		},
		conn: c,
	}
}

type {{$.Private}}OptSelect func(*{{$.Private}}Select)

func (c *{{$conn}}) Select(opts ...{{$.Private}}OptSelect) {{$.Public}}SelectExecer {
	s := &{{$.Private}}Select{
		params: runtime.SelectParams{
			Table: {{$.Private}}Table,
			OrderedColumns: {{$.Private}}OrderedColumns,
		},
		conn: c,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type {{$.Public}}SelectOpts struct {
	// Columns enable selection of specific columns in the returned structs
	// only the given columns will be filed with values from the database.
	Columns func(cols ...{{$.Private}}Column) {{$.Private}}OptSelect
	// Where applies where conditions on the query
	Where func(where runtime.Where) {{$.Private}}OptSelect
	// Limit applies rows limit on the query response
	Limit func(limit int64) {{$.Private}}OptSelect
	// Page applies rows offset and limit on the query response
	Page func(offset, limit int64) {{$.Private}}OptSelect
	{{ range $_, $e := $.Graph.Out -}}
	{{ $f := $e.LocalField -}}
	// Join{{$f.Name}} add a join query for {{$f.Name}}
	// Based on a forward relation
	Join{{$f.Name}} func(joiner {{$.Private}}{{$f.Type.Name}}Joiner) {{$.Private}}OptSelect
	{{ end -}}
	{{ range $_, $e := $.Graph.In -}}
	{{ $f := $e.LocalField -}}
	// Join{{$f.Name}} add a join query for {{$f.Name}}
	// Based on a reversed relation
	Join{{$f.Name}} func(joiner {{$.Private}}{{$f.Type.Name}}Joiner) {{$.Private}}OptSelect
	{{ end -}}
	// OrderBy set order to the query results according to column
	OrderBy func(col {{$.Private}}Column, dir orm.OrderDir) {{$.Private}}OptSelect
	// GroupBy make the query group by column
	GroupBy func(col {{$.Private}}Column) {{$.Private}}OptSelect
	// Context sets the context for the SQL query
	Context func(context.Context) {{$.Private}}OptSelect
}

var {{$.Private}}SelectOpts = {{$.Public}}SelectOpts {
	Columns: func(cols ...{{$.Private}}Column) {{$.Private}}OptSelect {
		return func(s *{{$.Private}}Select) {
			s.params.Columns = make(map[string]bool, len(cols))
			for _, col := range cols {
				s.params.Columns[string(col)] = true
			}
		}
	},
	Where: func(where runtime.Where) {{$.Private}}OptSelect {
		return func(s *{{$.Private}}Select) {
			s.params.Where = where
		}
	},
	Limit: func(limit int64) {{$.Private}}OptSelect {
		return func(s *{{$.Private}}Select) {
			s.params.Page.Limit = limit
		}
	},
	Page: func(offset, limit int64) {{$.Private}}OptSelect {
		return func(s *{{$.Private}}Select) {
			s.params.Page.Offset = offset
			s.params.Page.Limit = limit
		}
	},
	{{ range $_, $e := $.Graph.Out -}}
	{{ $f := $e.LocalField -}}
	Join{{$f.Name}}: func(joiner {{$.Private}}{{$f.Type.Name}}Joiner) {{$.Private}}OptSelect {
		return func(s *{{$.Private}}Select) {
			s.scan{{$f.Name}} = joiner
			s.params.Joins = append(s.params.Joins, runtime.JoinParams{
				Pairings: []runtime.Pairing{
					{{ range $i, $pk := $e.RelationType.PrimaryKeys -}}
					{
						Column: "{{(index $e.SrcField.Columns $i).Name}}",
						JoinedColumn: "{{$pk.Column.Name}}",
					},
					{{ end -}}
				},
				SelectParams: joiner.Params(),
			})
		}
	},
	{{ end -}}
	{{ range $_, $e := $.Graph.In -}}
	{{ $f := $e.LocalField -}}
	Join{{$f.Name}}: func(joiner {{$.Private}}{{$f.Type.Name}}Joiner) {{$.Private}}OptSelect {
		return func(s *{{$.Private}}Select) {
			s.scan{{$f.Name}} = joiner
			s.params.Joins = append(s.params.Joins, runtime.JoinParams{
				Pairings: []runtime.Pairing{
					{{ range $i, $pk := $e.RelationType.PrimaryKeys -}}
					{
						Column: "{{$pk.Column.Name}}",
						JoinedColumn: "{{(index $e.SrcField.Columns $i).Name}}",
					},
					{{ end -}}
				},
				SelectParams: joiner.Params(),
			})
		}
	},
	{{ end -}}
	OrderBy: func(col {{$.Private}}Column, dir orm.OrderDir) {{$.Private}}OptSelect {
		return func(s *{{$.Private}}Select) {
			s.params.Orders.Add(string(col), dir)
		}
	},
	GroupBy: func(col {{$.Private}}Column) {{$.Private}}OptSelect {
		return func(s *{{$.Private}}Select) {
			s.params.Groups.Add(string(col))
		}
	},
	Context: func(ctx context.Context) {{$.Private}}OptSelect {
		return func(s *{{$.Private}}Select) {
			s.params.Ctx = ctx
		}
	},
}

// Joiner returns an object to be used in a join operation with {{$name}}
func (b *{{$.Private}}Select) Joiner() {{$.Public}}Joiner {
	return &{{$.Private}}Joiner{builder: b}
}


{{ range $_, $refType := $.Graph.Type.ReferencedTypes -}}
// {{$.Private}}{{$refType.Name}}Joiner is a scanner that defined by .Select().Joiner()
// of an ORM object for type {{$refType.Name}}
type {{$.Private}}{{$refType.Name}}Joiner interface {
	Params() runtime.SelectParams
	Scan(dialect string, values []driver.Value{{if $refType.HasOneToManyRelation}}, exists map[string]*{{$refType.Ext $.Package}}{{end}}) (*{{$refType.Ext $.Package}}, int, error)
}
{{ end -}}

// Insert returns a builder of an SQL INSERT statement
func (c *{{$conn}}) Insert() *{{$.Public}}InsertBuilder {
	return &{{$.Public}}InsertBuilder{
		params: runtime.InsertParams{Table: {{$.Private}}Table},
		conn: c,
	}
}

// Update returns a builder of an SQL UPDATE statement
func (c *{{$conn}}) Update() *{{$.Public}}UpdateBuilder {
	return &{{$.Public}}UpdateBuilder{
		params: runtime.UpdateParams{Table: {{$.Private}}Table},
		conn: c,
	}
}

// Delete returns a builder of an SQL DELETE statement
func (c *{{$conn}}) Delete() *{{$.Public}}DeleteBuilder {
	return &{{$.Public}}DeleteBuilder{
		params: runtime.DeleteParams{Table: {{$.Private}}Table},
		conn: c,
	}
}

// Where returns a builder of an SQL WHERE statement
func (c *{{$conn}}) Drop() *{{$.Public}}DropBuilder {
	return &{{$.Public}}DropBuilder{
		params: runtime.DropParams{
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
	params runtime.CreateParams
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

// Context sets the context for the SQL query
func (b *{{$.Public}}CreateBuilder) Context(ctx context.Context) *{{$.Public}}CreateBuilder {
	b.params.Ctx = ctx
	return b
}

// === InsertBuilder ===

// {{$.Public}}InsertBuilder builds an INSERT statement parameters
type {{$.Public}}InsertBuilder struct {
	params runtime.InsertParams
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
	params runtime.UpdateParams
	conn   *{{$conn}}
}

// Where sets the WHERE statement to the SQL query
func (b *{{$.Public}}UpdateBuilder) Where(where runtime.Where) *{{$.Public}}UpdateBuilder {
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
	params runtime.DeleteParams
	conn   *{{$conn}}
}

// Where applies where conditions on the SQL query
func (b *{{$.Public}}DeleteBuilder) Where(w runtime.Where) *{{$.Public}}DeleteBuilder {
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
	params runtime.DropParams
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

{{ if $.Graph.Type.PrimaryKeys -}}

// === Get ===

func (c *{{$conn}}) Get({{range $i, $pk := $.Graph.Type.PrimaryKeys}}key{{$i}} {{$pk.Type.Ext $.Package}},{{end}}) (*{{$type}}, error) {
	return c.Select({{$.Private}}SelectOpts.Where(
	{{- range $i, $pk := $.Graph.Type.PrimaryKeys -}}
		c.Where().{{$pk.Name}}(orm.OpEq, key{{$i}})
		{{- if ne (plus1 $i) (len $.Graph.Type.PrimaryKeys) -}}
		.And(
		{{- end -}}
	{{end}}
	{{- range $i, $_ := $.Graph.Type.PrimaryKeys -}}
	{{- if ne (plus1 $i) (len $.Graph.Type.PrimaryKeys) -}}
	)
	{{- end -}}
	{{- end -}}
	)).First()
}
{{ end -}}

// Exec creates a table for the given struct
func (b *{{$.Public}}CreateBuilder) Exec() error {
	b.params.Ctx = runtime.ContextOrBackground(b.params.Ctx)
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
		_, err := tx.ExecContext(b.params.Ctx, stmt)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// query is used by the Select.Query and Select.Limit functions
func (b *{{$.Private}}Select) query() (*sql.Rows, error) {
	stmt, args := b.conn.dialect.Select(&b.params)
	return b.conn.QueryContext(b.params.Ctx, stmt, args...)
}

// Exec inserts the data to the given database
func (b *{{$.Public}}InsertBuilder) Exec() (*{{$type}}, error) {
	if len(b.params.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to insert")
	}
	b.params.Ctx = runtime.ContextOrBackground(b.params.Ctx)
	stmt, args := b.conn.dialect.Insert(&b.params)
	res, err := b.conn.ExecContext(b.params.Ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	return {{$.Private}}ReturnObject(b.params.Assignments, res)
}

// Exec inserts the data to the given database
func (b *{{$.Public}}UpdateBuilder) Exec() (sql.Result, error) {
	if len(b.params.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	b.params.Ctx = runtime.ContextOrBackground(b.params.Ctx)
	stmt, args := b.conn.dialect.Update(&b.params)
	return b.conn.ExecContext(b.params.Ctx, stmt, args...)
}

// Exec runs the delete statement on a given database.
func (b *{{$.Public}}DeleteBuilder) Exec() (sql.Result, error) {
	stmt, args := b.conn.dialect.Delete(&b.params)
	return b.conn.ExecContext(runtime.ContextOrBackground(b.params.Ctx), stmt, args...)
}

// Query the database
func (b *{{$.Private}}Select) Query() ([]{{$type}}, error) {
	b.params.Ctx = runtime.ContextOrBackground(b.params.Ctx)
	rows, err := b.query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		items []{{$type}}
		{{ if $hasOneToManyRelation -}}
		// exists is a mapping from primary key to already parsed structs
		exists = make(map[string]*{{$type}})
		{{ end -}}
	)
	for rows.Next() {
		// check context cancellation
		if err := b.params.Ctx.Err(); err != nil  {
			return nil, err
		}
		item, _, err := b.scan(b.conn.dialect.Name(), runtime.Values(*rows){{if $hasOneToManyRelation}}, exists{{end}})
		if err != nil {
			return nil, err
		}

		{{ if $hasOneToManyRelation -}}
		hash := {{$.Private}}HashItem(item)
		if exist := exists[hash]; exist != nil {
			{{ range $_, $f := $.Graph.Type.References -}}
			{{ if $f.Type.Slice -}}
			exist.{{$f.Name}} = append(exist.{{$f.Name}}, item.{{$f.Name}}...)
			{{ end -}}
			{{ end -}}
		} else {
			items = append(items, *item)
			exists[hash] = &items[len(items)-1]
		}
		{{ else -}}
		items = append(items, *item)
		{{ end -}}
	}
	return items, rows.Err()
}

// Count add a count column to the query
func (b *{{$.Private}}Select) Count() ([]{{$countStruct}}, error) {
	b.params.Ctx = runtime.ContextOrBackground(b.params.Ctx)
	b.params.Count = true
	rows, err := b.query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		items []{{$countStruct}}
		{{ if $hasOneToManyRelation -}}
		// exists is a mapping from primary key to already parsed structs
		exists = make(map[string]*{{$type}})
		{{ end -}}
	)
	for rows.Next() {
		// check context cancellation
		if err := b.params.Ctx.Err(); err != nil  {
			return nil, err
		}
		item, _, err := b.scanCount(b.conn.dialect.Name(), runtime.Values(*rows){{if $hasOneToManyRelation}}, exists{{end}})
		if err != nil {
			return nil, err
		}

		{{ if $hasOneToManyRelation -}}
		hash := {{$.Private}}HashItem(item.{{$name}})
		if exist := exists[hash]; exist != nil {
			{{ range $_, $f := $.Graph.Type.References -}}
			{{ if $f.Type.Slice -}}
			exist.{{$f.Name}} = append(exist.{{$f.Name}}, item.{{$f.Name}}...)
			{{ end -}}
			{{ end -}}
		} else {
			items = append(items, *item)
			exists[hash] = items[len(items)-1].{{$name}}
		}
		{{ else -}}
		items = append(items, *item)
		{{ end -}}
	}
	return items, rows.Err()
}

// First returns the first row that matches the query.
// If no row matches the query, an ErrNotFound will be returned.
// This call cancels any paging that was set with the
// {{$.Private}}Select previously.
func (b *{{$.Private}}Select) First() (*{{$type}}, error) {
	b.params.Ctx = runtime.ContextOrBackground(b.params.Ctx)
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
	item, _, err := b.scan(b.conn.dialect.Name(), runtime.Values(*rows){{if $hasOneToManyRelation}}, nil{{end}})
	if err != nil {
		return nil, err
	}
	return item, rows.Err()
}

// {{$.Private}}ReturnObject builds {{$type}} from assignment values
// and from the sql result ID, for returning an object from INSERT transactions
func {{$.Private}}ReturnObject(assignments runtime.Assignments, res sql.Result) (*{{$type}}, error) {
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
	{{ if eq (len $.Graph.Type.PrimaryKeys) 1 -}}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	ret.{{(index $.Graph.Type.PrimaryKeys 0).AccessName}} = {{(index $.Graph.Type.PrimaryKeys 0).Type.Ext $.Package}}(id)
	{{ end -}}
	return ret, nil
}

{{ if $hasOneToManyRelation -}}
func {{$.Private}}HashItem(item *{{$name}}) string {
	return fmt.Sprintf("{{repeat "%v" (len $.Graph.Type.PrimaryKeys)}}",
	{{- range $f := $.Graph.Type.PrimaryKeys -}}
	{{if $f.Type.Pointer}}*{{end}}item.{{$f.AccessName}},
	{{- end -}}
	)
}
{{ end -}}

// Exec runs the drop statement on a given database.
func (b *{{$.Public}}DropBuilder) Exec() error {
	stmt, args := b.conn.dialect.Drop(&b.params)
	_, err := b.conn.ExecContext(runtime.ContextOrBackground(b.params.Ctx), stmt, args...)
	return err
}

// {{$.Graph.Type.Naked.Name}}Joiner is an interface for joining a {{$name}} in a SELECT statement
// in another type
type {{$name}}Joiner interface {
	Params() runtime.SelectParams
	Scan(dialect string, values []driver.Value{{if $hasOneToManyRelation}}, exists map[string]*{{$type}}{{end}}) (*{{$type}}, int, error)
}

// {{$countStruct}} is a struct for counting rows of type {{$name}}
type {{$countStruct}} struct {
	*{{$.Graph.Type.Ext $.Package}}
	Count int64
}

// {{$.Private}}Select builds an SQL SELECT statement parameters
type {{$.Private}}Select struct {
	params runtime.SelectParams
	conn *{{$conn}}
	{{ range $_, $f := $.Graph.Type.References -}}
	scan{{$f.Name}} {{$.Private}}{{$f.Type.Naked.Name}}Joiner
	{{ end -}}
}

// {{$.Private}}Joiner represents a builder that exposes only Params and Scan method
type {{$.Private}}Joiner struct {
	builder *{{$.Private}}Select
}

func (j *{{$.Private}}Joiner) Params() runtime.SelectParams {
	return j.builder.params
}

func (j *{{$.Private}}Joiner) Scan(dialect string, values []driver.Value{{if $hasOneToManyRelation}}, exists map[string]*{{$type}}{{end}}) (*{{$type}}, int, error) {
	return j.builder.scan(dialect, values{{if $hasOneToManyRelation}}, exists{{end}})
}

// scan an SQL row to a {{$name}} struct
// It returns the scanned {{$.Graph.Type.Ext $.Package}} and the number of scanned fields,
// and an error in case of failure.
func (s *{{$.Private}}Select) scan(dialect string, vals []driver.Value{{if $hasOneToManyRelation}}, exists map[string]*{{$.Graph.Type.Ext $.Package}}{{end}}) (*{{$.Graph.Type.Ext $.Package}}, int, error) {
	item, n, err := s.scanCount(dialect, vals{{if $hasOneToManyRelation}}, exists{{end}})
	if err != nil {
		return nil, n, err
	}
	return item.{{$name}}, n, nil
}

// ScanCount scans an SQL row to a {{$countStruct}} struct
func (s *{{$.Private}}Select) scanCount(dialect string, vals []driver.Value{{if $hasOneToManyRelation}}, exists map[string]*{{$.Graph.Type.Ext $.Package}}{{end}}) (*{{$countStruct}}, int, error) {
	switch dialect {
	{{ range $_, $dialect := $.Dialects -}}
	case "{{$dialect.Name}}":
		return s.scan{{$dialect.Name}}(vals{{if $hasOneToManyRelation}}, exists{{end}})
	{{ end -}}
	default:
		return nil, 0, fmt.Errorf("unsupported dialect %s", dialect)
	}
}

{{ range $_, $dialect := $.Dialects }}
// scan{{$dialect.Name}} scans {{$dialect.Name}} row to a {{$name}} struct
func (s *{{$.Private}}Select) scan{{$dialect.Name}} (vals []driver.Value{{if $hasOneToManyRelation}}, exists map[string]*{{$.Graph.Type.Ext $.Package}}{{end}}) (*{{$countStruct}}, int, error) {
	var (
		row = new({{$countStruct}})
		i int
		rowExists bool
		allNils = true
		all = s.params.SelectAll()
	)
	row.{{$name}} = new({{$type}})
	{{ range $i, $f := $.Graph.Type.NonReferences -}}
	// scan column {{$i}}
	if all || s.params.Columns["{{$f.Column.Name}}"] {
		if i >= len(vals) {
			return nil, 0, fmt.Errorf("not enough columns returned: %d", len(vals))
		}
		if vals[i] != nil && !rowExists {
			allNils = false
{{ $dialect.ConvertValueCode $f -}}
		}
		{{ if and $hasOneToManyRelation (or $f.Unique $f.PrimaryKey) -}}
		// check if we scanned this item in previous rows. If we did, set rowExists,
		// so other columns in this table won't be evaluated. We only need values
		// from other tables.
		if exists[{{$.Private}}HashItem(row.{{$name}})] != nil {
			rowExists = true
		}
		{{ end -}}
		i++
	}
	{{ end -}}

	if s.params.Count {
		switch val := vals[i].(type) {
		case int64:
			row.Count = val
		case []byte:
			row.Count = runtime.ParseInt(val)
		default:
			return nil, 0, runtime.ErrConvert("COUNT(*)", i, vals[i], "int64, []byte")
		}
		i++
	}

	{{ range $_, $f := $.Graph.Type.References -}}
	if s := s.scan{{$f.Name}}; s != nil {
		tmp, n, err := s.Scan("{{$dialect.Name}}", vals[i:]{{if $f.Type.HasOneToManyRelation}}, nil{{end}})
		if err != nil {
			return nil, 0, fmt.Errorf("sub scanning {{$f.AccessName}}, cols [%d:%d]: %s", i, len(vals), err)
		}
		// If the result is nil, we want to discard it.
		// This is possible since we are doing a left join, if there was no match in the
		// right table, all it's columns are set to nil, and the result of a Scan function
		// is nil also.
		if tmp != nil {
			{{ if $f.Type.Slice -}}
			row.{{$f.AccessName}} = append(row.{{$f.AccessName}}, {{if not $f.Type.Pointer}}*{{end}}tmp)
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
	if allNils {
		row.{{$name}} = nil
	}

	return row, i, nil
}
{{ end }}

type {{$.Public}}WhereBuilder struct {}

{{ range $_, $f := $.Graph.Type.NonReferences -}}
// {{$.Public}}Where{{$f.Name}} adds a condition on {{$f.Name}} to the WHERE statement
func (*{{$.Public}}WhereBuilder) {{$f.Name}}(op orm.Op, val {{$f.Type.Ext $.Package}}) runtime.Where {
	return runtime.NewWhere(op, "{{$f.Column.Name}}", val)
}

// {{$.Public}}Where{{$f.Name}}In adds an IN condition on {{$f.Name}} to the WHERE statement
func (*{{$.Public}}WhereBuilder) {{$f.Name}}In(vals ...{{$f.Type.Ext $.Package}}) runtime.Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return runtime.NewWhereIn("{{$f.Column.Name}}", args...)
}

// {{$.Public}}Where{{$f.Name}}Between adds a BETWEEN condition on {{$f.Name}} to the WHERE statement
func (*{{$.Public}}WhereBuilder) {{$f.Name}}Between(low, high {{$f.Type.Ext $.Package}}) runtime.Where {
	return runtime.NewWhereBetween("{{$f.Column.Name}}", low, high)
}
{{ end -}}
`))
