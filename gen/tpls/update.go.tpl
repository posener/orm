package {{.Package}}

import (
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
	{{- end }}

    "{{.Type.ImportPath}}"
)

// Insert{{.Type.Name}} creates an UPDATE statement according to the given object
func (o *orm) Update{{.Type.Name}}(p *{{.Type.FullName}}) *Update {
	u := o.Update()
	{{- range $_, $f := .Type.Fields}}
	u.internal.Assignments.Add("{{$f.SQL.Column}}", p.{{$f.Name}})
	{{- end}}
	return u
}

{{range $_, $f := .Type.Fields}}
// Set{{$f.Name}} sets value for column {{$f.SQL.Column}} in the UPDATE statement
func (u *Update) Set{{$f.Name}}(value {{$f.Type}}) *Update {
	u.internal.Assignments.Add("{{$f.SQL.Column}}", value)
	return u
}
{{end}}
