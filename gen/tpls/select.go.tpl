package {{.Package}}

import (
	"log"
	"strings"

    "{{.Type.ImportPath}}"
)

// String returns the SQL query string
func (s *TSelect) String() string {
    return strings.Join([]string{
        "SELECT", s.selectString(), "FROM {{.Type.Table}}",
        s.where.String(),
        s.page.String(),
    }, " ")

}

// Exec runs the Query on a given database.
func (s *TSelect) Query() ([]{{.Type.FullName}}, error) {
	// create select statement
	stmt := s.String()
	args := s.where.Args()
	log.Printf("Query: '%v' %v", stmt, args)
	rows, err := s.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.FullName}}
	for rows.Next() {
		var i {{.Type.FullName}}
		if err := rows.Scan(s.scanArgs(&i)...); err != nil {
			return nil, err
		}
		all = append(all, i)
	}
	return all, rows.Err()
}

{{ range $_, $f := .Type.Fields -}}
// Select{{$f.Name}} Add {{$f.Name}} to the selected column of a query
func (s *TSelect) Select{{$f.Name}}() *TSelect {
    s.columns = append(s.columns, "{{$f.ColumnName}}")
    return s
}
{{ end -}}

// scanArgs are list of fields to be given to the sql Scan command
func (s *TSelect) scanArgs(p *{{.Type.FullName}}) []interface{} {
	if len(s.columns) == 0 {
        // add to args all the fields of p
        return []interface{}{
            {{ range $_, $f := .Type.Fields -}}
            &p.{{$f.Name}},
            {{ end }}
        }
	}
	m := s.columnsMap()
	args := make([]interface{}, len(s.columns))
	{{ range $_, $f := .Type.Fields -}}
	if i := m["{{$f.ColumnName}}"]; i != 0 {
		args[i-1] = &p.{{$f.Name}}
	}
	{{ end -}}
	return args
}
