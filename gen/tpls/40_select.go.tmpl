type {{$.Type.PrefixPublic}}Scanner interface {
    Columns() []string
    First(dialect string, values []driver.Value{{if .Type.HasOneToManyRelation}}, exists map[{{.Type.PrimaryKey.Type.ExtName $.Type.Package}}]*{{.Type.ExtName $.Type.Package}}{{end}}) (*{{$.Type.ExtNaked $.Type.Package}}, error)
}

{{ range $_, $f := $.Type.References }}
type {{$.Type.PrefixPrivate}}{{$f.Name}}Scanner interface {
    Columns() []string
    First(dialect string, values []driver.Value{{if $f.Type.HasOneToManyRelation}}, exists map[{{$f.Type.PrimaryKey.Type.ExtName $.Type.Package}}]*{{$f.Type.ExtName $.Type.Package}}{{end}}) (*{{$f.Type.ExtNaked $.Type.Package}}, error)
}
{{ end }}

// {{.Type.Name}}Count is a struct for counting rows of type {{.Type.Name}}
type {{.Type.Name}}Count struct {
    {{.Type.ExtName $.Type.Package}}
    Count int64
}

// {{$.Type.PrefixPublic}}SelectBuilder builds an SQL SELECT statement parameters
type {{$.Type.PrefixPublic}}SelectBuilder struct {
	params common.SelectParams
	conn *{{.Type.PrefixPrivate}}Conn
	selector {{$.Type.PrefixPrivate}}Selector
}

func (b *{{$.Type.PrefixPublic}}SelectBuilder) Scanner() {{$.Type.PrefixPublic}}Scanner {
    return b.params.Columns.(*{{$.Type.PrefixPrivate}}Selector)
}

// Where applies where conditions on the query
func (b *{{$.Type.PrefixPublic}}SelectBuilder) Where(where common.Where) *{{$.Type.PrefixPublic}}SelectBuilder {
	b.params.Where = where
	return b
}

// Limit applies rows limit on the query response
func (b *{{$.Type.PrefixPublic}}SelectBuilder) Limit(limit int64) *{{$.Type.PrefixPublic}}SelectBuilder {
	b.params.Page.Limit = limit
	return b
}

// Page applies rows offset and limit on the query response
func (b *{{$.Type.PrefixPublic}}SelectBuilder) Page(offset, limit int64) *{{$.Type.PrefixPublic}}SelectBuilder {
	b.params.Page.Offset = offset
	b.params.Page.Limit = limit
	return b
}

{{ range $_, $f := .Type.Fields -}}
{{ if not $f.IsReference }}
// Select{{$f.Name}} adds {{$f.Name}} to the selected column of a query
func (b *{{$.Type.PrefixPublic}}SelectBuilder) Select{{$f.Name}}() *{{$.Type.PrefixPublic}}SelectBuilder {
    b.selector.Select{{$f.Name}} = true
    return b
}
{{ else }}
// Join{{$f.Name}} add a join query for {{$f.Name}}
func (b *{{$.Type.PrefixPublic}}SelectBuilder) Join{{$f.Name}}(scanner {{$.Type.PrefixPrivate}}{{$f.Name}}Scanner) *{{$.Type.PrefixPublic}}SelectBuilder {
    b.selector.Join{{$f.Name}} = scanner
    return b
}
{{ end }}

{{ if not $f.Type.Slice -}}
// OrderBy{{$f.Name}} set order to the query results according to column {{$f.Column}}
func (b *{{$.Type.PrefixPublic}}SelectBuilder) OrderBy{{$f.Name}}(dir common.OrderDir) *{{$.Type.PrefixPublic}}SelectBuilder {
    b.params.Orders.Add("{{$f.Column}}", dir)
    return b
}

// GroupBy{{$f.Name}} make the query group by column {{$f.Column}}
func (b *{{$.Type.PrefixPublic}}SelectBuilder) GroupBy{{$f.Name}}() *{{$.Type.PrefixPublic}}SelectBuilder {
    b.params.Groups.Add("{{$f.Column}}")
    return b
}
{{ end -}}
{{ end -}}

// Context sets the context for the SQL query
func (b *{{$.Type.PrefixPublic}}SelectBuilder) Context(ctx context.Context) *{{$.Type.PrefixPublic}}SelectBuilder {
	b.params.Ctx = ctx
	return b
}
