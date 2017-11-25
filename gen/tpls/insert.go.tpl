package {{.Package}}

import (
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
	{{- end }}

    "{{.Type.ImportPath}}"
)

// Insert{{.Type.Name}} creates an INSERT statement according to the given object
func (o *orm) Insert{{.Type.Name}}(p *{{.Type.FullName}}) *Insert {
	i := o.Insert()
	{{- range $_, $f := .Type.Fields }}
	{{- if not $f.SQL.Auto }}
	i.internal.Assignments.Add("{{$f.SQL.Column}}", p.{{$f.Name}})
	{{- end -}}
	{{- end }}
	return i
}

{{- range $_, $f := .Type.Fields }}
{{ if not $f.SQL.Auto -}}
// Set{{$f.Name}} sets value for column {{$f.SQL.Column}} in the INSERT statement
func (i *Insert) Set{{$f.Name}}(value {{$f.Type}}) *Insert {
	i.internal.Assignments.Add("{{$f.SQL.Column}}", value)
	return i
}
{{ end -}}
{{ end -}}
