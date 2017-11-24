package {{.Package}}

type columns struct {
    {{ range $i, $f := .Type.Fields -}}
    Select{{$f.Name}} bool
    {{ end }}
    count bool // used for sql COUNT(*) column
}

// String is the SQL representation of columns
func (c columns) String() string {
	var s string
    {{ range $i, $f := .Type.Fields -}}
    if c.Select{{$f.Name}} {
        s += "`{{$f.SQL.Column}}`, "
    }
    {{ end }}
    if c.count {
        s += "COUNT(*), "
    }
	if len(s) == 0 { // no specific column was selected
		return "*"
	}
	return s[:len(s)-2]
}

// selectAll returns true if no column was specifically selected
func (c *columns) selectAll() bool {
    return {{ range $i, $f := .Type.Fields -}} !c.Select{{$f.Name}} && {{end}} !c.count
}
