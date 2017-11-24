package {{.Package}}

import (
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
	{{- end }}
	"github.com/posener/orm/dialect/{{.Dialect.Name}}"

    "{{.Type.ImportPath}}"
)

func (i *TInsert) String() string {
    return {{.Dialect.Name}}.Insert(i.orm, i.assign)
}

// Insert{{.Type.Name}} creates an INSERT statement according to the given object
func (o *ORM) Insert{{.Type.Name}}(p *{{.Type.FullName}}) *TInsert {
	i := o.Insert()
	{{- range $_, $f := .Type.Fields }}
	i.add("{{$f.SQL.Column}}", p.{{$f.Name}})
	{{- end }}
	return i
}

{{range $_, $f := .Type.Fields}}
// Set{{$f.Name}} sets value for column {{$f.SQL.Column}} in the INSERT statement
func (i *TInsert) Set{{$f.Name}}(value {{$f.Type}}) *TInsert {
	return i.add("{{$f.SQL.Column}}", value)
}
{{end}}
