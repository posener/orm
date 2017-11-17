package personorm

import (
	"github.com/posener/orm/where"
)


func WhereName(op where.Op, val string) where.Options {
	return where.New(op, "name", val)
}

func WhereAge(op where.Op, val int) where.Options {
	return where.New(op, "age", val)
}

