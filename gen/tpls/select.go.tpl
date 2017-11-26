package {{.Package}}

import (
	"github.com/posener/orm/common"

    "{{.Type.ImportPath}}"
)

type {{.Type.Name}}Count struct {
    {{.Type.FullName}}
    Count int64
}

// SelectBuilder builds an SQL SELECT statement parameters
type SelectBuilder struct {
	params common.SelectParams
	orm *orm
	columns columns
}

// Where applies where conditions on the query
func (s *SelectBuilder) Where(where common.Where) *SelectBuilder {
	s.params.Where = where
	return s
}

// Limit applies rows limit on the query response
func (s *SelectBuilder) Limit(limit int64) *SelectBuilder {
	s.params.Page.Limit = limit
	return s
}

// Page applies rows offset and limit on the query response
func (s *SelectBuilder) Page(offset, limit int64) *SelectBuilder {
	s.params.Page.Offset = offset
	s.params.Page.Limit = limit
	return s
}

{{ range $_, $f := .Type.Fields -}}
// Select{{$f.Name}} adds {{$f.Name}} to the selected column of a query
func (s *SelectBuilder) Select{{$f.Name}}() *SelectBuilder {
    s.columns.Select{{$f.Name}} = true
    return s
}

// OrderBy{{$f.Name}} set order to the query results according to column {{$f.SQL.Column}}
func (s *SelectBuilder) OrderBy{{$f.Name}}(dir common.OrderDir) *SelectBuilder {
    s.params.Orders.Add("{{$f.SQL.Column}}", dir)
    return s
}

// GroupBy{{$f.Name}} make the query group by column {{$f.SQL.Column}}
func (s *SelectBuilder) GroupBy{{$f.Name}}() *SelectBuilder {
    s.params.Groups.Add("{{$f.SQL.Column}}")
    return s
}
{{ end -}}
