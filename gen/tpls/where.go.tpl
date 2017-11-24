package {{.Package}}

import (
	{{ range $_, $f := .Type.Fields -}}
	{{ if $f.ImportPath }}"{{$f.ImportPath}}"{{ end }}
	{{- end }}
)

{{range $_, $f := .Type.Fields}}
// Where{{$f.Name}} adds a condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}(op Op, val {{$f.Type}}) *Where {
	return newWhere(op, "{{$f.SQL.Column}}", val)
}

// Where{{$f.Name}}In adds an IN condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}In(vals ...{{$f.Type}}) *Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return newWhereIn("{{$f.SQL.Column}}", args...)
}

// Where{{$f.Name}}Between adds a BETWEEN condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}Between(low, high {{$f.Type}}) *Where {
	return newWhereBetween("{{$f.SQL.Column}}", low, high)
}
{{end}}
