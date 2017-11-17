package personorm

import (
	"strings"

	"github.com/posener/orm/example"
)

type SelectColumns struct {
	Name bool
	Age  bool
}

func (c *SelectColumns) String() string {
	// no special columns select, select everything
	if c == nil {
		return "*"
	}
	names := make([]string, 0, 2) // 2 is according to number of type fields
	if c.Name {
		names = append(names, "name")
	}
	if c.Age {
		names = append(names, "age")
	}
	return strings.Join(names, ", ")
}

func (c *SelectColumns) ScanArgs(p *example.Person) []interface{} {
	all := c == nil
	args := make([]interface{}, 0, 2) // 2 is according to number of type fields
	if all || c.Name {
		args = append(args, &p.Name)
	}
	if all || c.Age {
		args = append(args, &p.Age)
	}
	return args
}
