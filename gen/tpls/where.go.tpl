package {{.Package}}

import (
    {{ range $_, $import := .Type.FieldsImports -}}
    "{{$import}}"
    {{ end }}

	"github.com/posener/orm/common"
)

{{range $_, $f := .Type.Fields}}
// Where{{$f.VarName}} adds a condition on {{$f.VarName}} to the WHERE statement
func Where{{$f.VarName}}(op common.Op, val {{$f.ExtTypeName}}) common.Where {
	return common.NewWhere(op, "{{$f.SQL.Column}}", val)
}

// Where{{$f.VarName}}In adds an IN condition on {{$f.VarName}} to the WHERE statement
func Where{{$f.VarName}}In(vals ...{{$f.ExtTypeName}}) common.Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return common.NewWhereIn("{{$f.SQL.Column}}", args...)
}

// Where{{$f.VarName}}Between adds a BETWEEN condition on {{$f.VarName}} to the WHERE statement
func Where{{$f.VarName}}Between(low, high {{$f.ExtTypeName}}) common.Where {
	return common.NewWhereBetween("{{$f.SQL.Column}}", low, high)
}
{{end}}
