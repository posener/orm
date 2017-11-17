package {{.PackageName}}

import (
	"strings"

    "{{.Type.ImportPath}}"
)

type Select struct {
	Name bool
	Age  bool
}

func (c *Select) String() string {
    if c == nil {
        return "*"
    }
    // collect all fields which their Select.FieldName is true
	names := make([]string, 0, {{len .Type.Fields}}) // according to number of type fields
	{{ range $_, $f := .Type.Fields -}}
	if c.{{$f.Name}} {
	    names = append(names, "{{$f.ColumnName}}")
	}
	{{ end }}
	return strings.Join(names, ", ")
}

// ScanArgs are list of fields to be given to the sql Scan command
func (c *Select) scanArgs(p *{{.Type.Name}}) []interface{} {
	all := c == nil

	// add to args all the fields of p
	args := make([]interface{}, 0, {{len .Type.Fields}}) // according to number of type fields
	{{ range $_, $f := .Type.Fields -}}
	if all || c.{{$f.Name}} {
		args = append(args, &p.{{$f.Name}})
	}
	{{ end }}
	return args
}
