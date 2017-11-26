package {{.Package}}

import (
    "database/sql"
    {{ range $_, $f := .Type.Fields -}}
    {{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
    {{- end }}

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect"

    "{{.Type.ImportPath}}"
)

const table = "{{.Type.Table}}"

var createColumnsStatements = map[string]string{
    {{ range $_, $d := .Dialects -}}
    "{{$d.Name}}": "{{$d.ColumnsStatement}}",
    {{ end -}}
}

// Open opens database connection
func Open(driverName, dataSourceName string) (API, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	dialect, err := dialect.New(driverName)
	if err != nil {
		return nil, err
	}
	return &orm{dialect: dialect, db: db}, nil
}

// New returns an orm object from a db instance
func New(driverName string, db DB) (API, error) {
	dialect, err := dialect.New(driverName)
	if err != nil {
		return nil, err
	}
    return &orm{dialect: dialect, db: db}, nil
}

// Create returns a builder of an SQL CREATE statement
func (o *orm) Create() *CreateBuilder {
	return &CreateBuilder{
		params: common.CreateParams{
		    Table: table,
		    ColumnsStatement: createColumnsStatements[o.dialect.Name()],
        },
	    orm: o,
    }
}

// Select returns a builder of an SQL SELECT statement
func (o *orm) Select() *SelectBuilder {
	s := &SelectBuilder{
		params: common.SelectParams{Table: table},
		orm: o,
	}
    s.params.Columns = &s.columns
    return s
}

// Insert returns a builder of an SQL INSERT statement
func (o *orm) Insert() *InsertBuilder {
	return &InsertBuilder{
		params: common.InsertParams{Table: table},
		orm: o,
	}
}

// Insert{{.Type.Name}} returns an SQL INSERT statement builder filled with values of a given object
func (o *orm) Insert{{.Type.Name}}(p *{{.Type.FullName}}) *InsertBuilder {
	i := o.Insert()
	{{- range $_, $f := .Type.Fields }}
	{{- if not $f.SQL.Auto }}
	i.params.Assignments.Add("{{$f.SQL.Column}}", p.{{$f.Name}})
	{{- end -}}
	{{- end }}
	return i
}

// Update returns a builder of an SQL UPDATE statement
func (o *orm) Update() *UpdateBuilder {
	return &UpdateBuilder{
		params: common.UpdateParams{Table: table},
		orm: o,
    }
}

// Update{{.Type.Name}} returns an SQL UPDATE statement builder filled with values of a given object
func (o *orm) Update{{.Type.Name}}(p *{{.Type.FullName}}) *UpdateBuilder {
	u := o.Update()
	{{- range $_, $f := .Type.Fields }}
    {{- if not $f.SQL.Auto }}
	u.params.Assignments.Add("{{$f.SQL.Column}}", p.{{$f.Name}})
	{{- end -}}
	{{- end }}
	return u
}

// Delete returns a builder of an SQL DELETE statement
func (o *orm) Delete() *DeleteBuilder {
	return &DeleteBuilder{
		params: common.DeleteParams{Table: table},
		orm: o,
    }
}

{{- range $_, $f := .Type.Fields }}
{{ if not $f.SQL.Auto -}}
// Set{{$f.Name}} sets value for column {{$f.SQL.Column}} in the INSERT statement
func (i *InsertBuilder) Set{{$f.Name}}(value {{$f.Type}}) *InsertBuilder {
	i.params.Assignments.Add("{{$f.SQL.Column}}", value)
	return i
}

// Set{{$f.Name}} sets value for column {{$f.SQL.Column}} in the UPDATE statement
func (u *UpdateBuilder) Set{{$f.Name}}(value {{$f.Type}}) *UpdateBuilder {
	u.params.Assignments.Add("{{$f.SQL.Column}}", value)
	return u
}
{{ end -}}
{{ end -}}
