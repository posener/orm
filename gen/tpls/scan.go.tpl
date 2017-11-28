package {{.Package}}
import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"unsafe"
	"fmt"
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }} "{{$f.ImportPath}}"{{ end }}
	{{- end }}
)

const errMsg = "converting %s: column %d with value %v (type %T) to %s"

// scanArgs are list of fields to be given to the sql Scan command
func (c *columns) scan(dialect string, rows *sql.Rows) (*{{.Type.Name}}Count, error) {
    switch dialect {
    {{- range $_, $dialect := $.Dialects }}
    case "{{$dialect.Name}}":
        return c.scan{{$dialect.Name}}(rows)
    {{ end -}}
    default:
        return nil, fmt.Errorf("unsupported dialect %s", dialect)
    }
}

{{ range $_, $dialect := $.Dialects }}
func (c *columns) scan{{$dialect.Name}} (rows *sql.Rows) (*{{$.Type.Name}}Count, error) {
    var (
        vals = values(*rows)
        row {{$.Type.Name}}Count
        all = c.selectAll()
        i = 0
    )
    {{ range $_, $f := $.Type.Fields }}
    if all || c.Select{{$f.Name}} {
        if vals[i] != nil {
{{$dialect.ConvertValueCode $f}}
        }
        i++
    }
    {{ end }}

    if c.count {
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

// Values is a hack to the sql.Rows struct
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
