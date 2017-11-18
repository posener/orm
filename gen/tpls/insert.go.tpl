package {{.PackageName}}

import "{{.Type.ImportPath}}"

{{range $_, $f := .Type.Fields}}
func (i Insert) {{$f.Name}}(value {{$f.Type}}) Insert {
	return i.add("{{$f.ColumnName}}", value)
}
{{end}}

func (i Insert) {{.Type.Name}}(p *{{.Type.FullName}}) Insert {
	var j = i
	{{- range $_, $f := .Type.Fields}}
	j = j.add("{{$f.ColumnName}}", p.{{$f.Name}})
	{{- end}}
	return j
}

