package {{.Package}}

import (
    "context"
	"github.com/posener/orm/common"

    "{{.Type.ImportPath}}"
)

type {{.Type.Name}}Count struct {
    {{.Type.FullName}}
    Count int64
}

// Select is the struct that holds the SELECT data
type Select struct {
	internal common.Select
	orm *orm
	columns columns
}

// Where applies where conditions on the query
func (s *Select) Where(where common.Where) *Select {
	s.internal.Where = where
	return s
}

// Limit applies rows limit on the query response
func (s *Select) Limit(limit int64) *Select {
	s.internal.Page.Limit = limit
	return s
}

// Page applies rows offset and limit on the query response
func (s *Select) Page(offset, limit int64) *Select {
	s.internal.Page.Offset = offset
	s.internal.Page.Limit = limit
	return s
}

// Query the database
func (s *Select) Query(ctx context.Context) ([]{{.Type.FullName}}, error) {
    rows, err := s.query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.FullName}}
	for rows.Next() {
	    // check context cancellation
	    if err := ctx.Err(); err != nil  {
	        return nil, err
	    }
		item, err := scan(s.orm.dialect.Name(), s.columns, rows)
        if err != nil {
			return nil, err
		}
		all = append(all, item.{{.Type.Name}})
	}
	return all, rows.Err()
}

// Count add a count column to the query
func (s *Select) Count(ctx context.Context) ([]{{.Type.Name}}Count, error) {
    s.columns.count = true
    rows, err := s.query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.Name}}Count
	for rows.Next() {
	    // check context cancellation
	    if err := ctx.Err(); err != nil  {
	        return nil, err
	    }
		item, err := scan(s.orm.dialect.Name(), s.columns, rows)
        if err != nil {
			return nil, err
		}
		all = append(all, *item)
	}
	return all, rows.Err()
}

{{ range $_, $f := .Type.Fields -}}
// Select{{$f.Name}} adds {{$f.Name}} to the selected column of a query
func (s *Select) Select{{$f.Name}}() *Select {
    s.columns.Select{{$f.Name}} = true
    return s
}

// OrderBy{{$f.Name}} set order to the query results according to column {{$f.SQL.Column}}
func (s *Select) OrderBy{{$f.Name}}(dir common.OrderDir) *Select {
    s.internal.Orders.Add("{{$f.SQL.Column}}", dir)
    return s
}

// GroupBy{{$f.Name}} make the query group by column {{$f.SQL.Column}}
func (s *Select) GroupBy{{$f.Name}}() *Select {
    s.internal.Groups.Add("{{$f.SQL.Column}}")
    return s
}
{{ end -}}

