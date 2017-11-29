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
    Insert{{.Type.Name}}(*{{.Type.ExtTypeName}}) *InsertBuilder
    Update{{.Type.Name}}(*{{.Type.ExtTypeName}}) *UpdateBuilder

    Logger(Logger)
}

// Querier is the interface for a SELECT SQL statement
type Querier interface {
    Query() ([]{{.Type.ExtTypeName}}, error)
}

// Counter is the interface for a SELECT SQL statement for counting purposes
type Counter interface {
    Count() ([]{{.Type.Name}}Count, error)
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

// Insert{{.Type.Name}} returns an SQL INSERT statement builder filled with values of a given object
func (c *conn) Insert{{.Type.Name}}(p *{{.Type.ExtTypeName}}) *InsertBuilder {
	i := c.Insert()
	{{- range $_, $f := .Type.Fields }}
	{{- if not $f.SQL.Auto }}
	i.params.Assignments.Add("{{$f.SQL.Column}}", p.{{$f.Name}})
	{{- end -}}
	{{- end }}
	return i
}

// Update returns a builder of an SQL UPDATE statement
func (c *conn) Update() *UpdateBuilder {
	return &UpdateBuilder{
		params: common.UpdateParams{Table: table},
		conn: c,
    }
}

// Update{{.Type.Name}} returns an SQL UPDATE statement builder filled with values of a given object
func (c *conn) Update{{.Type.Name}}(p *{{.Type.ExtTypeName}}) *UpdateBuilder {
	u := c.Update()
	{{- range $_, $f := .Type.Fields }}
    {{- if not $f.SQL.Auto }}
	u.params.Assignments.Add("{{$f.SQL.Column}}", p.{{$f.Name}})
	{{- end -}}
	{{- end }}
	return u
}

// Delete returns a builder of an SQL DELETE statement
func (c *conn) Delete() *DeleteBuilder {
	return &DeleteBuilder{
		params: common.DeleteParams{Table: table},
		conn: c,
    }
}

{{- range $_, $f := .Type.Fields }}
{{ if not $f.SQL.Auto -}}
// Set{{$f.Name}} sets value for column {{$f.SQL.Column}} in the INSERT statement
func (i *InsertBuilder) Set{{$f.Name}}(value {{$f.ExtTypeName}}) *InsertBuilder {
	i.params.Assignments.Add("{{$f.SQL.Column}}", value)
	return i
}

// Set{{$f.Name}} sets value for column {{$f.SQL.Column}} in the UPDATE statement
func (u *UpdateBuilder) Set{{$f.Name}}(value {{$f.ExtTypeName}}) *UpdateBuilder {
	u.params.Assignments.Add("{{$f.SQL.Column}}", value)
	return u
}
{{ end -}}
{{ end -}}
