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

func Insert{{.Type.Name}}(p *{{.Type.FullName}}) *TInsert {
	var i TInsert
	{{- range $_, $f := .Type.Fields}}
	i.add("{{$f.ColumnName}}", p.{{$f.Name}})
	{{- end}}
	return &i
}

{{range $_, $f := .Type.Fields}}
func (i *TInsert) Set{{$f.Name}}(value {{$f.Type}}) *TInsert {
	return i.add("{{$f.ColumnName}}", value)
}
{{end}}
