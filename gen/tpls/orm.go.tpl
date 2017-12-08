// {{$.Type.PrefixPrivate}}Table is SQL table name
const {{$.Type.PrefixPrivate}}Table = "{{.Type.Table}}"

// {{$.Type.PrefixPrivate}}CreateColumnsStatements are columns definitions in different dialects
var {{$.Type.PrefixPrivate}}CreateColumnsStatements = map[string]string{
    {{ range $_, $d := .Dialects -}}
    "{{$d.Name}}": "{{$d.ColumnsStatement $.Type}}",
    {{ end -}}
}

// {{$.Type.PrefixPublic}}API is the interface of the ORM object
type {{$.Type.PrefixPublic}}API interface {
    Close() error
    Create() *{{$.Type.PrefixPublic}}CreateBuilder
    Select() *{{$.Type.PrefixPublic}}SelectBuilder
    Insert() *{{$.Type.PrefixPublic}}InsertBuilder
    Update() *{{$.Type.PrefixPublic}}UpdateBuilder
    Delete() *{{$.Type.PrefixPublic}}DeleteBuilder

    Logger(orm.Logger)
}

// {{$.Type.PrefixPublic}}Open opens database connection
func {{$.Type.PrefixPublic}}Open(driverName, dataSourceName string) ({{$.Type.PrefixPublic}}API, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return {{$.Type.PrefixPublic}}New(driverName, db)
}

// {{$.Type.PrefixPublic}}New returns an conn object from a db instance
func {{$.Type.PrefixPublic}}New(driverName string, db orm.DB) ({{$.Type.PrefixPublic}}API, error) {
	d, err := dialect.New(driverName)
	if err != nil {
		return nil, err
	}
    return &{{.Type.PrefixPrivate}}Conn{dialect: d, db: db}, nil
}

// Create returns a builder of an SQL CREATE statement
func (c *{{.Type.PrefixPrivate}}Conn) Create() *{{$.Type.PrefixPublic}}CreateBuilder {
	return &{{$.Type.PrefixPublic}}CreateBuilder{
		params: common.CreateParams{
		    Table: {{$.Type.PrefixPrivate}}Table,
		    ColumnsStatement: {{$.Type.PrefixPrivate}}CreateColumnsStatements[c.dialect.Name()],
        },
	    conn: c,
    }
}

// Select returns a builder of an SQL SELECT statement
func (c *{{.Type.PrefixPrivate}}Conn) Select() *{{$.Type.PrefixPublic}}SelectBuilder {
	s := &{{$.Type.PrefixPublic}}SelectBuilder{
		params: common.SelectParams{Table: {{$.Type.PrefixPrivate}}Table},
		conn: c,
	}
    s.params.Columns = &s.selector
    return s
}

// Insert returns a builder of an SQL INSERT statement
func (c *{{.Type.PrefixPrivate}}Conn) Insert() *{{$.Type.PrefixPublic}}InsertBuilder {
	return &{{$.Type.PrefixPublic}}InsertBuilder{
		params: common.InsertParams{Table: {{$.Type.PrefixPrivate}}Table},
		conn: c,
	}
}

// Update returns a builder of an SQL UPDATE statement
func (c *{{.Type.PrefixPrivate}}Conn) Update() *{{$.Type.PrefixPublic}}UpdateBuilder {
	return &{{$.Type.PrefixPublic}}UpdateBuilder{
		params: common.UpdateParams{Table: {{$.Type.PrefixPrivate}}Table},
		conn: c,
    }
}

// Delete returns a builder of an SQL DELETE statement
func (c *{{.Type.PrefixPrivate}}Conn) Delete() *{{$.Type.PrefixPublic}}DeleteBuilder {
	return &{{$.Type.PrefixPublic}}DeleteBuilder{
		params: common.DeleteParams{Table: {{$.Type.PrefixPrivate}}Table},
		conn: c,
    }
}

// Insert{{.Type.Name}} returns an SQL INSERT statement builder filled with values of a given object
func (b *{{$.Type.PrefixPublic}}InsertBuilder) Insert{{.Type.Name}}(p *{{.Type.ExtName $.Type.Package}}) *{{$.Type.PrefixPublic}}InsertBuilder {
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.IsSettable -}}
	{{ if not $f.IsReference -}}
	b.params.Assignments.Add("{{$f.Column}}", p.{{$f.Name}})
	{{ else -}}
	{{ if $f.Type.Pointer -}}
	if p.{{$f.Name}} != nil {
	{{ end -}}
	b.params.Assignments.Add("{{$f.Column}}", p.{{$f.Name}}.{{$f.Type.PrimaryKey.Name}})
	{{ if $f.Type.Pointer -}}
	}
	{{ end -}}
	{{ end -}}
	{{ end -}}
	{{ end -}}
	return b
}


// Update{{.Type.Name}} update values for all struct fields
func (b *{{$.Type.PrefixPublic}}UpdateBuilder) Update{{.Type.Name}}(p *{{.Type.ExtName $.Type.Package}}) *{{$.Type.PrefixPublic}}UpdateBuilder {
	{{ range $_, $f := .Type.Fields -}}
    {{ if $f.IsSettable -}}
	{{ if not $f.IsReference -}}
	b.params.Assignments.Add("{{$f.Column}}", p.{{$f.Name}})
	{{ else -}}
	{{ if $f.Type.Pointer -}}
	if p.{{$f.Name}} != nil {
	{{ end -}}
	b.params.Assignments.Add("{{$f.Column}}", p.{{$f.Name}}.{{$f.Type.PrimaryKey.Name}})
	{{ if $f.Type.Pointer -}}
	}
	{{ end -}}
	{{ end -}}
	{{ end -}}
	{{ end -}}
	return b
}

{{- range $_, $f := .Type.Fields }}

{{ if $f.IsSettable -}}
// Set{{$f.Name}} sets value for column {{$f.Column}} in the INSERT statement
func (b *{{$.Type.PrefixPublic}}InsertBuilder) Set{{$f.Name}}(value {{$f.SetType.ExtName $.Type.Package}}) *{{$.Type.PrefixPublic}}InsertBuilder {
	b.params.Assignments.Add("{{$f.Column}}", value)
	return b
}

// Set{{$f.Name}} sets value for column {{$f.Column}} in the UPDATE statement
func (b *{{$.Type.PrefixPublic}}UpdateBuilder) Set{{$f.Name}}(value {{$f.SetType.ExtName $.Type.Package}}) *{{$.Type.PrefixPublic}}UpdateBuilder {
	b.params.Assignments.Add("{{$f.Column}}", value)
	return b
}

{{ end -}}
{{ end -}}
