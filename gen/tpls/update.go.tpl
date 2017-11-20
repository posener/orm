package {{.Package}}

import (
    "fmt"

    "{{.Type.ImportPath}}"
)

func (u *TUpdate) String() string {
	return fmt.Sprintf(`UPDATE {{.Type.Table}} SET %s %s`,
	    u.assignmentList(),
	    u.where.String(),
	)
}

func (o *ORM) Update{{.Type.Name}}(p *{{.Type.FullName}}) *TUpdate {
	u := o.Update()
	{{- range $_, $f := .Type.Fields}}
	u.add("{{$f.ColumnName}}", p.{{$f.Name}})
	{{- end}}
	return u
}

{{range $_, $f := .Type.Fields}}
func (u *TUpdate) Set{{$f.Name}}(value {{$f.Type}}) *TUpdate {
	return u.add("{{$f.ColumnName}}", value)
}
{{end}}
