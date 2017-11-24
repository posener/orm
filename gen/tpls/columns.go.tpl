package {{.Package}}

type columns struct {
    {{ range $i, $f := .Type.Fields -}}
    Select{{$f.Name}} bool
    {{ end }}
    count bool // used for sql COUNT(*) column
}

// String is the SQL representation of columns
func (c *columns) Columns() []string {
	var cols []string
    {{ range $i, $f := .Type.Fields -}}
    if c.Select{{$f.Name}} {
        cols = append(cols, "{{$f.SQL.Column}}")
    }
    {{ end }}
	return cols
}

func (c *columns) Count() bool {
    return c.count
}

// selectAll returns true if no column was specifically selected
func (c *columns) selectAll() bool {
    return {{ range $i, $f := .Type.Fields -}} !c.Select{{$f.Name}} && {{end}} !c.count
}
