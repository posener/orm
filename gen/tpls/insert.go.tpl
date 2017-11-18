package {{.Package}}

import (
    "fmt"
    "strings"

    "{{.Type.ImportPath}}"
)

func (i TInsert) String() string {
	return fmt.Sprintf(`INSERT INTO {{.Type.Table}} (%s) VALUES (%s)`,
		strings.Join(i.cols, ", "),
		qMarks(len(i.values)),
	)
}

func (i TInsert) {{.Type.Name}}(p *{{.Type.FullName}}) TInsert {
	var j = i
	{{- range $_, $f := .Type.Fields}}
	j = j.add("{{$f.ColumnName}}", p.{{$f.Name}})
	{{- end}}
	return j
}

{{range $_, $f := .Type.Fields}}
func (i TInsert) {{$f.Name}}(value {{$f.Type}}) TInsert {
	return i.add("{{$f.ColumnName}}", value)
}
{{end}}
