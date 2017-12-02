package {{.Package}}

import (
    {{ range $_, $import := .Type.Imports -}}
    "{{$import}}"
    {{ end }}
	"github.com/posener/orm/common"
)

{{range $_, $f := .Type.Fields}}
// Where{{$f.Name}} adds a condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}(op common.Op, val {{$f.Type.ExtName}}) common.Where {
	return common.NewWhere(op, "{{$.Type.Table}}", "{{$f.Column}}", val)
}

// Where{{$f.Name}}In adds an IN condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}In(vals ...{{$f.Type.ExtName}}) common.Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return common.NewWhereIn("{{$.Type.Table}}", "{{$f.Column}}", args...)
}

// Where{{$f.Name}}Between adds a BETWEEN condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}Between(low, high {{$f.Type.ExtName}}) common.Where {
	return common.NewWhereBetween("{{$.Type.Table}}", "{{$f.Column}}", low, high)
}
{{end}}
