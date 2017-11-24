package {{.Package}}

import (
	"strings"
	"database/sql/driver"
	"fmt"
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
	{{- end }}

    "{{.Type.ImportPath}}"
)

const colCount = "COUNT(*)"

// TSelect is the struct that holds the SELECT data
type TSelect struct {
	Querier
	Argser
	orm *ORM
	columns columns
	where *Where
	groupBy
	orderBy
	page Page
}

func (s *TSelect) Args() []interface{} {
	return s.where.Args()
}

// Where applies where conditions on the query
func (s *TSelect) Where(where *Where) *TSelect {
	s.where = where
	return s
}

// Limit applies rows limit on the query response
func (s *TSelect) Limit(limit int64) *TSelect {
	s.page.limit = limit
	return s
}

// Page applies rows offset and limit on the query response
func (s *TSelect) Page(offset, limit int64) *TSelect {
	s.page.offset = offset
	s.page.limit = limit
	return s
}

type {{.Type.Name}}Count struct {
    {{.Type.FullName}}
    Count int64
}

// String returns the SQL query string
func (s *TSelect) String() string {
    return strings.Join([]string{
        "SELECT", s.columns.String(), "FROM '{{.Type.Table}}'",
        s.where.String(),
        s.groupBy.String(),
        s.orderBy.String(),
        s.page.String(),
    }, " ")

}

// Query the database
func (s *TSelect) Query() ([]{{.Type.FullName}}, error) {
	// create select statement
	stmt := s.String()
	args := s.Args()
	s.orm.log("Query: '%v' %v", stmt, args)
	dbRows, err := s.orm.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer dbRows.Close()

	rows := Rows{Rows: dbRows} // this is a hack to access lastcols field

	// extract rows to structures
	var all []{{.Type.FullName}}
	for rows.Next() {
		item, err := s.scan(rows.Values())
        if err != nil {
			return nil, err
		}
		all = append(all, item.{{.Type.Name}})
	}
	return all, rows.Err()
}

// Count add a count column to the query
func (s *TSelect) Count() ([]{{.Type.Name}}Count, error) {
    s.columns.count = true
	// create select statement
	stmt := s.String()
	args := s.where.Args()
	s.orm.log("Count: '%v' %v", stmt, args)
	dbRows, err := s.orm.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer dbRows.Close()

	rows := Rows{Rows: dbRows} // this is a hack to access lastcols field

	// extract rows to structures
	var all []{{.Type.Name}}Count
	for rows.Next() {
		item, err := s.scan(rows.Values())
        if err != nil {
			return nil, err
		}
		all = append(all, *item)
	}
	return all, rows.Err()
}

{{ range $_, $f := .Type.Fields -}}
// Select{{$f.Name}} Add {{$f.Name}} to the selected column of a query
func (s *TSelect) Select{{$f.Name}}() *TSelect {
    s.columns.Select{{$f.Name}} = true
    return s
}

// OrderBy{{$f.Name}} set order to the query results according to column {{$f.SQL.Column}}
func (s *TSelect) OrderBy{{$f.Name}}(dir OrderDir) *TSelect {
    s.orderBy.add("`{{$f.SQL.Column}}`", dir)
    return s
}

// GroupBy{{$f.Name}} make the query group by column {{$f.SQL.Column}}
func (s *TSelect) GroupBy{{$f.Name}}() *TSelect {
    s.groupBy.add("`{{$f.SQL.Column}}`")
    return s
}
{{ end -}}

// scanArgs are list of fields to be given to the sql Scan command
func (s *TSelect) scan(vals []driver.Value) (*{{.Type.Name}}Count, error) {
    var (
        row {{.Type.Name}}Count
        all = s.columns.selectAll()
        i = 0
    )
	{{ range $_, $f := .Type.Fields -}}
	if all || s.columns.Select{{$f.Name}} {
	    if vals[i] != nil {
            val, ok := vals[i].({{$f.SQL.ConvertType}})
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
