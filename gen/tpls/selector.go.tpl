package {{.Package}}

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"unsafe"
	"fmt"
    {{ range $_, $import := .Type.FieldsImports -}}
    "{{$import}}"
    {{ end }}
)

const errMsg = "converting %s: column %d with value %v (type %T) to %s"

// selector selects columns for SQL queries and for parsing SQL rows
type selector struct {
    {{ range $i, $f := .Type.Fields -}}
    Select{{$f.VarName}} bool
    {{ end }}
    count bool // used for sql COUNT(*) column
}

// Columns are the names of selected columns
func (s *selector) Columns() []string {
	var cols []string
    {{ range $i, $f := .Type.Fields -}}
    if s.Select{{$f.VarName}} {
        cols = append(cols, "{{$f.SQL.Column}}")
    }
    {{ end }}
	return cols
}

// Count is true when a COUNT(*) column should be added to the query
func (s *selector) Count() bool {
    return s.count
}

// scan an SQL row to a {{.Type.Name}}Count struct
func (s *selector) scan(dialect string, rows *sql.Rows) (*{{.Type.Name}}Count, error) {
    switch dialect {
    {{- range $_, $dialect := $.Dialects }}
    case "{{$dialect.Name}}":
        return s.scan{{$dialect.Name}}(rows)
    {{ end -}}
    default:
        return nil, fmt.Errorf("unsupported dialect %s", dialect)
    }
}

{{ range $_, $dialect := $.Dialects }}
// scan{{$dialect.Name}} scans {{$dialect.Name}} row to a {{$.Type.Name}} struct
func (s *selector) scan{{$dialect.Name}} (rows *sql.Rows) (*{{$.Type.Name}}Count, error) {
    var (
        vals = values(*rows)
        row {{$.Type.Name}}Count
        all = s.selectAll()
        i = 0
    )
    {{ range $_, $f := $.Type.Fields }}
    if all || s.Select{{$f.VarName}} {
        if vals[i] != nil {
{{$dialect.ConvertValueCode $f}}
        }
        i++
    }
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
    }

    return &row, nil
}
{{ end }}

// selectAll returns true if no column was specifically selected
func (s *selector) selectAll() bool {
    return {{ range $i, $f := .Type.Fields -}} !s.Select{{$f.VarName}} && {{end}} !s.count
}

// values is a hack to the sql.Rows struct
// since the rows struct does not expose it's lastcols values, or a way to give
// a custom scanner to the Scan method.
// See issue https://github.com/golang/go/issues/22544
func values(r sql.Rows) []driver.Value {
	// some ugly hack to access lastcols field
	rs := reflect.ValueOf(&r).Elem()
	rf := rs.FieldByName("lastcols")

	// overcome panic reflect.Value.Interface: cannot return value obtained from unexported field or method
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	return rf.Interface().([]driver.Value)
}
