package {{.PackageName}}

import (
	"github.com/posener/orm/where"
)

{{range $_, $f := .Type.Fields}}
func Where{{$f.Name}}(op where.Op, val {{$f.Type}}) where.Options {
	return where.New(op, "{{$f.ColumnName}}", val)
}

func Where{{$f.Name}}In(vals ...{{$f.Type}}) where.Options {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return where.NewMul(where.OpIn, "{{$f.ColumnName}}", args...)
}

func Where{{$f.Name}}Between(low, high {{$f.Type}}) where.Options {
	return where.NewMul(where.OpBetween, "{{$f.ColumnName}}", low, high)
}
{{end}}
