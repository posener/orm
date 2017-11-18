package {{.Package}}

{{$n := len .Type.Fields}}

func (c TCreate) String() string {
    // Create statement has a line for each variable with it's name and it's type.
	return `CREATE TABLE {{.Table}} (
	{{- range $i, $f := .Type.Fields }}
		{{$f.ColumnName}} {{$f.ColumnSQLType}}{{if lt (plus1 $i) $n}},{{end}}
    {{- end }}
	)`
}
