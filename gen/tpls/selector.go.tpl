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
    {{ end -}}
    {{ end -}}
    count bool // used for sql COUNT(*) column
}

// Columns are the names of selected columns
func (s *selector) Columns() []string {
	var cols []string
    {{ range $_, $f := .Type.Fields -}}
    {{ if not $f.IsReference -}}
    if s.Select{{$f.Name}} {
        cols = append(cols, "{{$f.Column}}")
    }
    {{ end -}}
    {{ end -}}
	return cols
}

// Joins are join options of the query
func (s *selector) Joins() []common.JoinParams {
	var joins []common.JoinParams
    {{ range $_, $f := .Type.References -}}
    {{ if $f.Type.Slice -}}
    if selector := s.Join{{$f.Name}}; selector != nil {
        // adding join of one to many relation, column in other type points to this type
        joins = append(joins, common.JoinParams{
            ForeignKey: common.ForeignKey{
                // column in this type that the other type is pointing on
                Column: "{{$f.ForeignKey.Column}}",
                // other type table
                RefTable: "{{$f.ForeignKey.RefTable}}",
                // other type column that points to this type
                RefColumn: "{{$f.ForeignKey.RefColumn}}",
            },
            SelectColumns: selector.Columns(),
        })
    }
    {{ else -}}
    if selector := s.Join{{$f.Name}}; selector != nil {
        // join that this type points to another type's primary key
        // this types [Column] points to [RefTable].[RefColumn]
        joins = append(joins, common.JoinParams{
            ForeignKey: common.ForeignKey{
                Column: "{{$f.Column}}",
                RefTable: "{{$f.Type.Table}}",
                RefColumn: "{{$f.Type.PrimaryKey.Column}}",
            },
            SelectColumns: selector.Columns(),
        })
    }
    {{ end }}
    {{ end }}
	return joins
}

// Count is true when a COUNT(*) column should be added to the query
func (s *selector) Count() bool {
    return s.count
}

// FirstCount scans an SQL row to a {{.Type.Name}}Count struct
func (s *selector) FirstCount(dialect string, vals []driver.Value{{if .Type.HasOneToManyRelation}}, exists map[{{.Type.PrimaryKey.Type.ExtName}}]*{{.Type.ExtName}}{{end}}) (*{{.Type.Name}}Count, error) {
    switch dialect {
    {{- range $_, $dialect := $.Dialects }}
    case "{{$dialect.Name}}":
        return s.scan{{$dialect.Name}}(vals{{if $.Type.HasOneToManyRelation}}, exists{{end}})
    {{ end -}}
    default:
        return nil, fmt.Errorf("unsupported dialect %s", dialect)
    }
}
// First scans an SQL row to a {{.Type.Name}} struct
func (s *selector) First(dialect string, vals []driver.Value{{if .Type.HasOneToManyRelation}}, exists map[{{.Type.PrimaryKey.Type.ExtName}}]*{{.Type.ExtName}}{{end}}) (*{{.Type.ExtName}}, error) {
    item, err := s.FirstCount(dialect, vals{{if .Type.HasOneToManyRelation}}, exists{{end}})
    if err != nil {
        return nil, err
    }
    return &item.{{.Type.Name}}, nil
}

{{ range $_, $dialect := $.Dialects }}
// scan{{$dialect.Name}} scans {{$dialect.Name}} row to a {{$.Type.Name}} struct
func (s *selector) scan{{$dialect.Name}} (vals []driver.Value{{if $.Type.HasOneToManyRelation}}, exists map[{{$.Type.PrimaryKey.Type.ExtName}}]*{{$.Type.ExtName}}{{end}}) (*{{$.Type.Name}}Count, error) {
    var (
        row {{$.Type.Name}}Count
        all = s.selectAll()
        i int
        rowExists bool
    )
    {{ range $_, $f := $.Type.Fields }}
    {{ if not $f.IsReference }}
    if all || s.Select{{$f.Name}} {
        if vals[i] != nil && !rowExists {
{{$dialect.ConvertValueCode $f}}
        }
        {{ if and $.Type.HasOneToManyRelation -}}
        {{ if eq $f.Name $.Type.PrimaryKey.Name -}}
        // check if we scanned this item in previous rows. If we did, set `rowExists`,
        // so other columns in this table won't be evaluated. We only need values
        // from other tables.
        if exists[row.{{$.Type.PrimaryKey.Name}}] != nil {
            rowExists = true
        }
        {{ end -}}
        {{ end -}}
        i++
    }
    {{ else if not $f.Type.Slice }}
    if all { // skip foreign key column
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
        tmp, err := j.First("{{$dialect.Name}}", vals[i:])
        if err != nil {
            return nil, err
        }
        {{ if $f.Type.Slice -}}
        row.{{$f.Name}} = append(row.{{$f.Name}}, {{if not $f.Type.Pointer}}*{{end}}tmp)
        {{ else -}}
        row.{{$f.Name}} = {{ if not $f.Type.Pointer}}*{{end}}tmp
        {{ end -}}
    }
    {{ end }}

    return &row, nil
}
{{ end }}

// selectAll returns true if no column was specifically selected
func (s *selector) selectAll() bool {
    return {{ range $i, $f := .Type.Fields -}}{{if not $f.IsReference }} !s.Select{{$f.Name}} && {{end}}{{end}} !s.count
}
