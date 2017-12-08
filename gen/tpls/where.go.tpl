{{range $_, $f := .Type.Fields}}
// {{$.Type.PrefixPublic}}Where{{$f.Name}} adds a condition on {{$f.Name}} to the WHERE statement
func {{$.Type.PrefixPublic}}Where{{$f.Name}}(op common.Op, val {{$f.Type.ExtName $.Type.Package}}) common.Where {
	return common.NewWhere(op, "{{$.Type.Table}}", "{{$f.Column}}", val)
}

// {{$.Type.PrefixPublic}}Where{{$f.Name}}In adds an IN condition on {{$f.Name}} to the WHERE statement
func {{$.Type.PrefixPublic}}Where{{$f.Name}}In(vals ...{{$f.Type.ExtName $.Type.Package}}) common.Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return common.NewWhereIn("{{$.Type.Table}}", "{{$f.Column}}", args...)
}

// {{$.Type.PrefixPublic}}Where{{$f.Name}}Between adds a BETWEEN condition on {{$f.Name}} to the WHERE statement
func {{$.Type.PrefixPublic}}Where{{$f.Name}}Between(low, high {{$f.Type.ExtName $.Type.Package}}) common.Where {
	return common.NewWhereBetween("{{$.Type.Table}}", "{{$f.Column}}", low, high)
}
{{end}}
