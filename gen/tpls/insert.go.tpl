package {{.Package}}

import (
    "database/sql"
    "fmt"
    "strings"
    "log"

    "{{.Type.ImportPath}}"
)

func (i Insert) {{.Type.Name}}(p *{{.Type.FullName}}) Insert {
	var j = i
	{{- range $_, $f := .Type.Fields}}
	j = j.add("{{$f.ColumnName}}", p.{{$f.Name}})
	{{- end}}
	return j
}

{{range $_, $f := .Type.Fields}}
func (i Insert) {{$f.Name}}(value {{$f.Type}}) Insert {
	return i.add("{{$f.ColumnName}}", value)
}
{{end}}

// Create creates a table for {{.Type.Name}}
func (i Insert) Exec(db *sql.DB) error {
	if len(i.cols) == 0 || len(i.values) == 0 {
		return fmt.Errorf("nothing to insert")
	}
	stmt := fmt.Sprintf(`INSERT INTO {{.Table}} (%s) VALUES (%s)`,
		strings.Join(i.cols, ", "),
		qMarks(len(i.values)),
	)

	log.Printf("Insert: '%v' (%v)", stmt, i.values)
	_, err := db.Exec(stmt, i.values...)
	return err
}

