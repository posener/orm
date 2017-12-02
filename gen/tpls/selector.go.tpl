package {{.Package}}

import (
	"database/sql/driver"
	"fmt"
    {{ range $_, $import := .Type.Imports -}}
    "{{$import}}"
    {{ end }}
    "github.com/posener/orm/common"
)

const errMsg = "converting %s: column %d with value %v (type %T) to %s"

// selector selects columns for SQL queries and for parsing SQL rows
type selector struct {
    {{ range $i, $f := .Type.Fields -}}
    {{ if not $f.IsReference -}}
    Select{{$f.Name}} bool
    {{ else -}}
    Join{{$f.Name}} {{$f.Name}}Scanner
    {{ end }}
    {{ end }}
    count bool // used for sql COUNT(*) column
}

// Columns are the names of selected columns
func (s *selector) Columns() []string {
	var cols []string
    {{ range $i, $f := .Type.Fields -}}
    {{ if not $f.IsReference -}}
    if s.Select{{$f.Name}} {
        cols = append(cols, "{{$f.Column}}")
    }
    {{ end }}
    {{ end }}
	return cols
}

// Joins are join options of the query
func (s *selector) Joins() []common.Join {
	var joins []common.Join
    {{ range $i, $f := .Type.References -}}
    if s.Join{{$f.Name}} != nil {
        joins = append(joins, common.Join{
            Column: "{{$f.Column}}",
            RefTable: "{{$f.Type.Table}}",
            RefColumn: "{{$f.Type.PrimaryKey.Column}}",
        })
    }
    {{ end }}
	return joins
}

// Count is true when a COUNT(*) column should be added to the query
func (s *selector) Count() bool {
    return s.count
}

// FirstCount scans an SQL row to a {{.Type.Name}}Count struct
func (s *selector) FirstCount(dialect string, vals []driver.Value) (*{{.Type.Name}}Count, error) {
    switch dialect {
    {{- range $_, $dialect := $.Dialects }}
    case "{{$dialect.Name}}":
        return s.scan{{$dialect.Name}}(vals)
    {{ end -}}
    default:
        return nil, fmt.Errorf("unsupported dialect %s", dialect)
    }
}
// First scans an SQL row to a {{.Type.Name}} struct
func (s *selector) First(dialect string, vals []driver.Value) (*{{.Type.ExtName}}, error) {
    item, err := s.FirstCount(dialect, vals)
    if err != nil {
        return nil, err
    }
    return &item.{{.Type.Name}}, nil
}

{{ range $_, $dialect := $.Dialects }}
// scan{{$dialect.Name}} scans {{$dialect.Name}} row to a {{$.Type.Name}} struct
func (s *selector) scan{{$dialect.Name}} (vals []driver.Value) (*{{$.Type.Name}}Count, error) {
    var (
        row {{$.Type.Name}}Count
        all = s.selectAll()
        i int
    )
    {{ range $_, $f := $.Type.Fields }}
    {{ if not $f.IsReference }}
    if all || s.Select{{$f.Name}} {
        if vals[i] != nil {
{{$dialect.ConvertValueCode $f}}
        }
        i++
    }
    {{ else }}
    if all || s.Join{{$f.Name}} != nil {
        i++
    }
    {{ end }}
    {{ end }}

    if s.count {
        switch val := vals[i].(type) {
        case int64:
            row.Count = val
        case []byte:
            row.Count = parseInt(val)
        default:
            return nil, fmt.Errorf(errMsg, "COUNT(*)", i, vals[i], vals[i], "int64, []byte")
        }
        i++
    }

    {{ range $_, $f := $.Type.References }}
    if j := s.Join{{$f.Name}}; j != nil {
        var err error
        row.{{$f.Name}}, err = j.First("{{$dialect.Name}}", vals[i:])
        if err != nil {
            return nil, err
        }
    }
    {{ end }}

    return &row, nil
}
{{ end }}

// selectAll returns true if no column was specifically selected
func (s *selector) selectAll() bool {
    return {{ range $i, $f := .Type.Fields -}}{{if not $f.IsReference }} !s.Select{{$f.Name}} && {{end}}{{end}} !s.count
}
