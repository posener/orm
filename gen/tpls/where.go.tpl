package {{.Package}}

import (
    {{ range $_, $import := .Type.FieldsImports -}}
    "{{$import}}"
    {{ end }}

	"github.com/posener/orm/common"
)

{{range $_, $f := .Type.Fields}}
// Where{{$f.Name}} adds a condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}(op common.Op, val {{$f.ExtTypeName}}) common.Where {
	return common.NewWhere(op, "{{$f.SQL.Column}}", val)
}

// Where{{$f.Name}}In adds an IN condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}In(vals ...{{$f.ExtTypeName}}) common.Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return common.NewWhereIn("{{$f.SQL.Column}}", args...)
}

// Where{{$f.Name}}Between adds a BETWEEN condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}Between(low, high {{$f.ExtTypeName}}) common.Where {
	return common.NewWhereBetween("{{$f.SQL.Column}}", low, high)
}
{{end}}
