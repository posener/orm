package {{.PackageName}}

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/posener/orm/tools"

	"{{.Type.ImportPath}}"
)

func NewInsert() Insert {
	return Insert{}
}

type Insert struct {
	cols   []string
	values []interface{}
}

{{range $_, $f := .Type.Fields}}
func (i Insert) {{$f.Name}}(value {{$f.Type}}) Insert {
	return i.add("{{$f.ColumnName}}", value)
}
{{end}}

func (i Insert) {{.Type.Name}}(p *{{.Type.FullName}}) Insert {
	var j = i
	{{- range $_, $f := .Type.Fields}}
	j = j.add("{{$f.ColumnName}}", p.{{$f.Name}})
	{{- end}}
	return j
}

// Create creates a table for example.Person
func (i Insert) Exec(db *sql.DB) error {
	if len(i.cols) == 0 || len(i.values) == 0 {
		return fmt.Errorf("nothing to insert")
	}
	stmt := fmt.Sprintf(`INSERT INTO person (%s) VALUES (%s)`,
		strings.Join(i.cols, ", "),
		tools.QMarks(len(i.values)),
	)

	log.Printf("Insert: '%v' (%v)", stmt, i.values)
	_, err := db.Exec(stmt, i.values...)
	return err
}

func (i Insert) add(name string, value interface{}) Insert {
	i.cols = append(i.cols, name)
	i.values = append(i.values, value)
	return i
}
