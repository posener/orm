package {{.Package}}

import (
    "fmt"
    "strings"

    "{{.Type.ImportPath}}"
)

func (i *TInsert) String() string {
	return fmt.Sprintf(`INSERT INTO {{.Type.Table}} (%s) VALUES (%s)`,
		strings.Join(i.cols, ", "),
		qMarks(len(i.values)),
	)
}

// Insert{{.Type.Name}} creates an INSERT statement according to the given object
func (o *ORM) Insert{{.Type.Name}}(p *{{.Type.FullName}}) *TInsert {
	i := o.Insert()
	{{- range $_, $f := .Type.Fields}}
	i.add("{{$f.ColumnName}}", p.{{$f.Name}})
	{{- end}}
	return i
}

{{range $_, $f := .Type.Fields}}
// Set{{$f.Name}} sets value for column {{$f.ColumnName}} in the INSERT statement
func (i *TInsert) Set{{$f.Name}}(value {{$f.Type}}) *TInsert {
	return i.add("{{$f.ColumnName}}", value)
}
{{end}}
