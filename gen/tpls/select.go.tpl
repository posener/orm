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
    s.columns.add(colCount)
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
    s.columns.add("`{{$f.ColumnName}}`")
    return s
}

// OrderBy{{$f.Name}} set order to the query results according to column {{$f.ColumnName}}
func (s *TSelect) OrderBy{{$f.Name}}(dir OrderDir) *TSelect {
    s.orderBy.add("`{{$f.ColumnName}}`", dir)
    return s
}

// GroupBy{{$f.Name}} make the query group by column {{$f.ColumnName}}
func (s *TSelect) GroupBy{{$f.Name}}() *TSelect {
    s.groupBy.add("`{{$f.ColumnName}}`")
    return s
}
{{ end -}}

// scanArgs are list of fields to be given to the sql Scan command
func (s *TSelect) scan(vals []driver.Value) (*{{.Type.Name}}Count, error) {
    var row {{.Type.Name}}Count
	if len(s.columns) == 0 {
        // add to args all the fields of row
        {{ range $i, $f := .Type.Fields -}}
        if vals[{{$i}}] != nil {
            val, ok := vals[{{$i}}].({{$f.SQL.ConvertType}})
            if !ok {
                return nil, fmt.Errorf("converting {{$f.Name}} column {{$i}} with value %v to {{$f.Type}}", vals[{{$i}}])
            }
            row.{{$f.Name}} = {{$f.Type}}(val)
        }
        {{ end }}
	}
	m := s.columns.indexMap()
	{{ range $_, $f := .Type.Fields -}}
	if i := m["`{{$f.ColumnName}}`"]-1; i != -1 {
	    val, ok := vals[i].({{$f.SQL.ConvertType}})
        if !ok {
            return nil, fmt.Errorf("converting {{$f.Name}}: column %d with value %v to {{$f.Type}}", i, vals[i])
        }
        row.{{$f.Name}} = {{$f.Type}}(val)
	}
	{{ end -}}
	if i := m[colCount]-1; i != -1 {
        var ok bool
        row.Count, ok = vals[i].(int64)
        if !ok {
            return nil, fmt.Errorf("converting COUNT(*): column %d with value %v to int64", i, vals[i])
        }
    }
	return &row, nil
}
