package {{.Package}}

import (
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
	{{- end }}

    "{{.Type.ImportPath}}"
)

// Insert{{.Type.Name}} creates an INSERT statement according to the given object
func (o *ORM) Insert{{.Type.Name}}(p *{{.Type.FullName}}) *Insert {
	i := o.Insert()
	{{- range $_, $f := .Type.Fields }}
	i.Assignments.Add("{{$f.SQL.Column}}", p.{{$f.Name}})
	{{- end }}
	return i
}

{{range $_, $f := .Type.Fields}}
// Set{{$f.Name}} sets value for column {{$f.SQL.Column}} in the INSERT statement
func (i *Insert) Set{{$f.Name}}(value {{$f.Type}}) *Insert {
	i.Assignments.Add("{{$f.SQL.Column}}", value)
	return i
}
{{end}}
