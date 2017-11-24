package {{.Package}}

import (
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
	{{- end }}
	"github.com/posener/orm/dialect/{{.Dialect.Name}}"

    "{{.Type.ImportPath}}"
)

func (u *TUpdate) String() string {
    return {{.Dialect.Name}}.Update(u.orm, u.assign, u.where)
}

// Insert{{.Type.Name}} creates an UPDATE statement according to the given object
func (o *ORM) Update{{.Type.Name}}(p *{{.Type.FullName}}) *TUpdate {
	u := o.Update()
	{{- range $_, $f := .Type.Fields}}
	u.add("{{$f.SQL.Column}}", p.{{$f.Name}})
	{{- end}}
	return u
}

{{range $_, $f := .Type.Fields}}
// Set{{$f.Name}} sets value for column {{$f.SQL.Column}} in the UPDATE statement
func (u *TUpdate) Set{{$f.Name}}(value {{$f.Type}}) *TUpdate {
	return u.add("{{$f.SQL.Column}}", value)
}
{{end}}
