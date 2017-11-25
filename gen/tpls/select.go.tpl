package {{.Package}}

import (
	"database/sql/driver"
	"fmt"
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
	{{- end }}
	"github.com/posener/orm"
	"github.com/posener/orm/row"

    "{{.Type.ImportPath}}"
)

type {{.Type.Name}}Count struct {
    {{.Type.FullName}}
    Count int64
}

// Select is the struct that holds the SELECT data
type Select struct {
	orm.Select
	orm *ORM
	columns
}

// Where applies where conditions on the query
func (s *Select) Where(where orm.Where) *Select {
	s.Select.Where = where
	return s
}

// Limit applies rows limit on the query response
func (s *Select) Limit(limit int64) *Select {
	s.Select.Page.Limit = limit
	return s
}

// Page applies rows offset and limit on the query response
func (s *Select) Page(offset, limit int64) *Select {
	s.Select.Page.Offset = offset
	s.Select.Page.Limit = limit
	return s
}

// Query the database
func (s *Select) Query() ([]{{.Type.FullName}}, error) {
    rows, err := s.query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.FullName}}
	for rows.Next() {
		item, err := s.scan(row.Values(*rows))
        if err != nil {
			return nil, err
		}
		all = append(all, item.{{.Type.Name}})
	}
	return all, rows.Err()
}

// Count add a count column to the query
func (s *Select) Count() ([]{{.Type.Name}}Count, error) {
    s.columns.count = true
    rows, err := s.query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.Name}}Count
	for rows.Next() {
		item, err := s.scan(row.Values(*rows))
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
func (s *Select) OrderBy{{$f.Name}}(dir orm.OrderDir) *Select {
    s.Orders.Add("{{$f.SQL.Column}}", dir)
    return s
}

// GroupBy{{$f.Name}} make the query group by column {{$f.SQL.Column}}
func (s *Select) GroupBy{{$f.Name}}() *Select {
    s.Groups.Add("{{$f.SQL.Column}}")
    return s
}
{{ end -}}

// scanArgs are list of fields to be given to the sql Scan command
func (s *Select) scan(vals []driver.Value) (*{{.Type.Name}}Count, error) {
    var (
        row {{.Type.Name}}Count
        all = s.columns.selectAll()
        i = 0
    )
	{{ range $_, $f := .Type.Fields -}}
	if all || s.columns.Select{{$f.Name}} {
	    if vals[i] != nil {
            val, ok := vals[i].({{$f.ConvertType}})
            if !ok {
                return nil, fmt.Errorf("converting {{$f.Name}}: column %d with value %v to {{$f.Type}}", i, vals[i])
            }
            row.{{$f.Name}} = {{$f.Type}}(val)
        }
        i++
	}
	{{ end -}}
	if s.columns.count {
        var ok bool
        row.Count, ok = vals[i].(int64)
        if !ok {
            return nil, fmt.Errorf("converting COUNT(*): column %d with value %v to int64", i, vals[i])
        }
    }
	return &row, nil
}
