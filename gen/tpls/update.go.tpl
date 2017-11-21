package {{.Package}}

import (
    "fmt"
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
	{{- end }}

    "{{.Type.ImportPath}}"
)

func (u *TUpdate) String() string {
	return fmt.Sprintf(`UPDATE {{.Type.Table}} SET %s %s`,
	    u.assignmentList(),
	    u.where.String(),
	)
}

// Insert{{.Type.Name}} creates an UPDATE statement according to the given object
func (o *ORM) Update{{.Type.Name}}(p *{{.Type.FullName}}) *TUpdate {
	u := o.Update()
	{{- range $_, $f := .Type.Fields}}
	u.add("{{$f.ColumnName}}", p.{{$f.Name}})
	{{- end}}
	return u
}

{{range $_, $f := .Type.Fields}}
// Set{{$f.Name}} sets value for column {{$f.ColumnName}} in the UPDATE statement
func (u *TUpdate) Set{{$f.Name}}(value {{$f.Type}}) *TUpdate {
	return u.add("{{$f.ColumnName}}", value)
}
{{end}}
