package {{.Package}}

{{range $_, $f := .Type.Fields}}
// Where{{$f.Name}} adds a condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}(op Op, val {{$f.Type}}) Where {
	return newWhere(op, "{{$f.ColumnName}}", val)
}

// Where{{$f.Name}}In adds an IN condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}In(vals ...{{$f.Type}}) Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return newMulWhere(OpIn, "{{$f.ColumnName}}", args...)
}

// Where{{$f.Name}}Between adds a BETWEEN condition on {{$f.Name}} to the WHERE statement
func Where{{$f.Name}}Between(low, high {{$f.Type}}) Where {
	return newMulWhere(OpBetween, "{{$f.ColumnName}}", low, high)
}
{{end}}
