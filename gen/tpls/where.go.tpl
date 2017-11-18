package {{.PackageName}}

{{range $_, $f := .Type.Fields}}
func Where{{$f.Name}}(op Op, val {{$f.Type}}) Where {
	return newWhere(op, "{{$f.ColumnName}}", val)
}

func Where{{$f.Name}}In(vals ...{{$f.Type}}) Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return newMulWhere(OpIn, "{{$f.ColumnName}}", args...)
}

func Where{{$f.Name}}Between(low, high {{$f.Type}}) Where {
	return newMulWhere(OpBetween, "{{$f.ColumnName}}", low, high)
}
{{end}}
