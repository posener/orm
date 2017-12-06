package {{.Package}}

import (
    "context"
	"database/sql/driver"

	"github.com/posener/orm/common"
    "{{.Type.ImportPath}}"
)

type Scanner interface {
    Columns() []string
    First(dialect string, values []driver.Value{{if .Type.HasOneToManyRelation}}, exists map[{{.Type.PrimaryKey.Type.ExtName}}]*{{.Type.ExtName}}{{end}}) (*{{$.Type.ExtNaked}}, error)
}

{{ range $_, $f := $.Type.References }}
type {{$f.Name}}Scanner interface {
    Columns() []string
    First(dialect string, values []driver.Value{{if $f.Type.HasOneToManyRelation}}, exists map[{{$f.Type.PrimaryKey.Type.ExtName}}]*{{$f.Type.ExtName}}{{end}}) (*{{$f.Type.ExtNaked}}, error)
}
{{ end }}

// {{.Type.Name}}Count is a struct for counting rows of type {{.Type.Name}}
type {{.Type.Name}}Count struct {
    {{.Type.ExtName}}
    Count int64
}

// SelectBuilder builds an SQL SELECT statement parameters
type SelectBuilder struct {
	params common.SelectParams
	conn *conn
	selector selector
}

func (b *SelectBuilder) Scanner() Scanner {
    return b.params.Columns.(*selector)
}

// Where applies where conditions on the query
func (b *SelectBuilder) Where(where common.Where) *SelectBuilder {
	b.params.Where = where
	return b
}

// Limit applies rows limit on the query response
func (b *SelectBuilder) Limit(limit int64) *SelectBuilder {
	b.params.Page.Limit = limit
	return b
}

// Page applies rows offset and limit on the query response
func (b *SelectBuilder) Page(offset, limit int64) *SelectBuilder {
	b.params.Page.Offset = offset
	b.params.Page.Limit = limit
	return b
}

{{ range $_, $f := .Type.Fields -}}
{{ if not $f.IsReference }}
// Select{{$f.Name}} adds {{$f.Name}} to the selected column of a query
func (b *SelectBuilder) Select{{$f.Name}}() *SelectBuilder {
    b.selector.Select{{$f.Name}} = true
    return b
}
{{ else }}
// Join{{$f.Name}} add a join query for {{$f.Name}}
func (b *SelectBuilder) Join{{$f.Name}}(scanner {{$f.Name}}Scanner) *SelectBuilder {
    b.selector.Join{{$f.Name}} = scanner
    return b
}
{{ end }}

{{ if not $f.Type.Slice -}}
// OrderBy{{$f.Name}} set order to the query results according to column {{$f.Column}}
func (b *SelectBuilder) OrderBy{{$f.Name}}(dir common.OrderDir) *SelectBuilder {
    b.params.Orders.Add("{{$f.Column}}", dir)
    return b
}

// GroupBy{{$f.Name}} make the query group by column {{$f.Column}}
func (b *SelectBuilder) GroupBy{{$f.Name}}() *SelectBuilder {
    b.params.Groups.Add("{{$f.Column}}")
    return b
}
{{ end -}}
{{ end -}}

// Context sets the context for the SQL query
func (b *SelectBuilder) Context(ctx context.Context) *SelectBuilder {
	b.params.Ctx = ctx
	return b
}
