package {{.PackageName}}

import (
	"github.com/posener/orm/where"
)

{{range $_, $f := .Type.Fields}}
func Where{{$f.Name}}(op where.Op, val {{$f.Type}}) where.Options {
	return where.New(op, "{{$f.ColumnName}}", val)
}
{{end}}
