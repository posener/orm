package {{.Package}}

import (
    "context"

	"github.com/posener/orm/common"
    "{{.Type.ImportPath}}"
)

// {{.Type.Name}}Count is a struct for counting rows of type {{.Type.Name}}
type {{.Type.Name}}Count struct {
    {{.Type.ExtTypeName}}
    Count int64
}

// SelectBuilder builds an SQL SELECT statement parameters
type SelectBuilder struct {
	params common.SelectParams
	conn *conn
	selector selector
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
// Select{{$f.VarName}} adds {{$f.VarName}} to the selected column of a query
func (b *SelectBuilder) Select{{$f.VarName}}() *SelectBuilder {
    b.selector.Select{{$f.VarName}} = true
    return b
}

// OrderBy{{$f.VarName}} set order to the query results according to column {{$f.SQL.Column}}
func (b *SelectBuilder) OrderBy{{$f.VarName}}(dir common.OrderDir) *SelectBuilder {
    b.params.Orders.Add("{{$f.SQL.Column}}", dir)
    return b
}

// GroupBy{{$f.VarName}} make the query group by column {{$f.SQL.Column}}
func (b *SelectBuilder) GroupBy{{$f.VarName}}() *SelectBuilder {
    b.params.Groups.Add("{{$f.SQL.Column}}")
    return b
}
{{ end -}}

// Context sets the context for the SQL query
func (b *SelectBuilder) Context(ctx context.Context) *SelectBuilder {
	b.params.Ctx = ctx
	return b
}
