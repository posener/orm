package personorm

import (
	"github.com/posener/orm/where"
)


func WhereName(op where.Op, val string) where.Options {
	return where.New(op, "name", val)
}

func WhereNameIn(vals ...string) where.Options {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return where.NewMul(where.OpIn, "name", args...)
}

func WhereNameBetween(low, high string) where.Options {
	return where.NewMul(where.OpBetween, "name", low, high)
}

func WhereAge(op where.Op, val int) where.Options {
	return where.New(op, "age", val)
}

func WhereAgeIn(vals ...int) where.Options {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return where.NewMul(where.OpIn, "age", args...)
}

func WhereAgeBetween(low, high int) where.Options {
	return where.NewMul(where.OpBetween, "age", low, high)
}

