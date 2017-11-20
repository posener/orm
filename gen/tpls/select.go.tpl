package {{.Package}}

import (
	"strings"

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
        "SELECT", s.columns.String(), "FROM {{.Type.Table}}",
        s.where.String(),
        s.groupBy.String(),
        s.page.String(),
    }, " ")

}

// Query the database
func (s *TSelect) Query() ([]{{.Type.FullName}}, error) {
	// create select statement
	stmt := s.String()
	args := s.where.Args()
	s.orm.log("Query: '%v' %v", stmt, args)
	rows, err := s.orm.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.FullName}}
	for rows.Next() {
		var item {{.Type.Name}}Count
		if err := rows.Scan(s.scanArgs(&item)...); err != nil {
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
	rows, err := s.orm.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.Name}}Count
	for rows.Next() {
		var item {{.Type.Name}}Count
		if err := rows.Scan(s.scanArgs(&item)...); err != nil {
			return nil, err
		}
		all = append(all, item)
	}
	return all, rows.Err()
}

{{ range $_, $f := .Type.Fields -}}
// Select{{$f.Name}} Add {{$f.Name}} to the selected column of a query
func (s *TSelect) Select{{$f.Name}}() *TSelect {
    s.columns.add("{{$f.ColumnName}}")
    return s
}

// GroupBy{{$f.Name}} make the query group by column {{$f.ColumnName}}
func (s *TSelect) GroupBy{{$f.Name}}() *TSelect {
    s.groupBy.add("{{$f.ColumnName}}")
    return s
}
{{ end -}}

// scanArgs are list of fields to be given to the sql Scan command
func (s *TSelect) scanArgs(p *{{.Type.Name}}Count) []interface{} {
	if len(s.columns) == 0 {
        // add to args all the fields of p
        return []interface{}{
            {{ range $_, $f := .Type.Fields -}}
            &p.{{$f.Name}},
            {{ end }}
        }
	}
	m := s.columns.indexMap()
	args := make([]interface{}, len(s.columns))
	{{ range $_, $f := .Type.Fields -}}
	if i := m["{{$f.ColumnName}}"]; i != 0 {
		args[i-1] = &p.{{$f.Name}}
	}
	{{ end -}}
	if i := m[colCount]; i != 0 {
        args[i-1] = &p.Count
    }
	return args
}
