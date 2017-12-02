package {{.Package}}

import (
    "database/sql"
    {{ range $_, $import := .Type.FieldsImports -}}
    "{{$import}}"
    {{ end }}

	"github.com/posener/orm"
	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect"

    "{{.Type.ImportPath}}"
)

// table is SQL table name
const table = "{{.Type.Table}}"

// createColumnsStatements are columns definitions in different dialects
var createColumnsStatements = map[string]string{
    {{ range $_, $d := .Dialects -}}
    "{{$d.Name}}": "{{$d.ColumnsStatement}}",
    {{ end -}}
}

// API is the interface of the ORM object
type API interface {
    Close() error
    Create() *CreateBuilder
    Select() *SelectBuilder
    Insert() *InsertBuilder
    Update() *UpdateBuilder
    Delete() *DeleteBuilder

    Logger(orm.Logger)
}

// Querier is the interface for a SELECT SQL statement
type Querier interface {
    Query() ([]{{.Type.ExtTypeName}}, error)
}

// Counter is the interface for a SELECT SQL statement for counting purposes
type Counter interface {
    Count() ([]{{.Type.Name}}Count, error)
}

// Firster is the interface for a SELECT SQL statement for getting only the
// first item. if no item matches the query, an `orm.ErrNotFound` will be returned.
type Firster interface {
	First() (*{{.Type.ExtTypeName}}, error)
}

// Open opens database connection
func Open(driverName, dataSourceName string) (API, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	d, err := dialect.New(driverName)
	if err != nil {
		return nil, err
	}
	return &conn{dialect: d, db: db}, nil
}

// New returns an conn object from a db instance
func New(driverName string, db orm.DB) (API, error) {
	d, err := dialect.New(driverName)
	if err != nil {
		return nil, err
	}
    return &conn{dialect: d, db: db}, nil
}

// Create returns a builder of an SQL CREATE statement
func (c *conn) Create() *CreateBuilder {
	return &CreateBuilder{
		params: common.CreateParams{
		    Table: table,
		    ColumnsStatement: createColumnsStatements[c.dialect.Name()],
        },
	    conn: c,
    }
}

// Select returns a builder of an SQL SELECT statement
func (c *conn) Select() *SelectBuilder {
	s := &SelectBuilder{
		params: common.SelectParams{Table: table},
		conn: c,
	}
    s.params.Columns = &s.selector
    return s
}

// Insert returns a builder of an SQL INSERT statement
func (c *conn) Insert() *InsertBuilder {
	return &InsertBuilder{
		params: common.InsertParams{Table: table},
		conn: c,
	}
}

// Update returns a builder of an SQL UPDATE statement
func (c *conn) Update() *UpdateBuilder {
	return &UpdateBuilder{
		params: common.UpdateParams{Table: table},
		conn: c,
    }
}

// Delete returns a builder of an SQL DELETE statement
func (c *conn) Delete() *DeleteBuilder {
	return &DeleteBuilder{
		params: common.DeleteParams{Table: table},
		conn: c,
    }
}

// Insert{{.Type.Name}} returns an SQL INSERT statement builder filled with values of a given object
func (b *InsertBuilder) Insert{{.Type.Name}}(p *{{.Type.ExtName}}) *InsertBuilder {
	{{- range $_, $f := .Type.Fields }}
	{{- if $f.IsSettable }}
	b.params.Assignments.Add("{{$f.Column}}", p.{{$f.Name}}{{if $f.IsReference }}.{{$f.Type.PrimaryKey.Name}}{{end}})
	{{- end -}}
	{{- end }}
	return b
}


// Update{{.Type.Name}} update values for all struct fields
func (b *UpdateBuilder) Update{{.Type.Name}}(p *{{.Type.ExtName}}) *UpdateBuilder {
	{{- range $_, $f := .Type.Fields }}
    {{- if $f.IsSettable }}
	b.params.Assignments.Add("{{$f.Column}}", p.{{$f.Name}}{{if $f.IsReference }}.{{$f.Type.PrimaryKey.Name}}{{end}})
	{{- end -}}
	{{- end }}
	return b
}

{{- range $_, $f := .Type.Fields }}

{{ if $f.IsSettable -}}
// Set{{$f.Name}} sets value for column {{$f.Column}} in the INSERT statement
func (b *InsertBuilder) Set{{$f.Name}}(value {{$f.SetType}}) *InsertBuilder {
	b.params.Assignments.Add("{{$f.Column}}", value)
	return b
}

// Set{{$f.Name}} sets value for column {{$f.Column}} in the UPDATE statement
func (b *UpdateBuilder) Set{{$f.Name}}(value {{$f.SetType}}) *UpdateBuilder {
	b.params.Assignments.Add("{{$f.Column}}", value)
	return b
}

{{ end -}}
{{ end -}}
