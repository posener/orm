package personorm

import (
	"strings"

	"github.com/posener/orm/where"
)

type WhereOptions struct {
	Name *where.String
	Age  *where.Int
}

func (w *WhereOptions) Args() []interface{} {
	if w == nil {
		return nil
	}
	args := make([]interface{}, 0, 1)
	if w.Name != nil {
		args = append(args, w.Name.Val)
	}
	if w.Age != nil {
		args = append(args, w.Age.Val)
	}
	return args
}

func (c *WhereOptions) String() string {
	if c == nil {
		return ""
	}
	names := make([]string, 0, 2) // 2 is according to number of type fields
	if c.Name != nil {
		names = append(names, "name=?")
	}
	if c.Age != nil {
		names = append(names, "age=?")
	}
	return "WHERE " + strings.Join(names, " AND ")
}
