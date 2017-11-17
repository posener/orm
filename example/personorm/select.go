package personorm

import (
	"strings"

	"github.com/posener/orm/example"
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
	names := make([]string, 0, 2) // according to number of type fields
	if c.Name {
		names = append(names, "name")
	}
	if c.Age {
		names = append(names, "age")
	}

	return strings.Join(names, ", ")
}

// ScanArgs are list of fields to be given to the sql Scan command
func (c *Select) scanArgs(p *example.Person) []interface{} {
	all := c == nil

	// add to args all the fields of p
	args := make([]interface{}, 0, 2) // according to number of type fields
	if all || c.Name {
		args = append(args, &p.Name)
	}
	if all || c.Age {
		args = append(args, &p.Age)
	}

	return args
}
