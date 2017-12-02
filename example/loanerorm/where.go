// Package loanerorm was auto-generated by github.com/posener/orm; DO NOT EDIT
package loanerorm

import (
	"github.com/posener/orm/example"

	"github.com/posener/orm/common"
)

// WhereID adds a condition on ID to the WHERE statement
func WhereID(op common.Op, val int64) common.Where {
	return common.NewWhere(op, "loaner", "id", val)
}

// WhereIDIn adds an IN condition on ID to the WHERE statement
func WhereIDIn(vals ...int64) common.Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return common.NewWhereIn("loaner", "id", args...)
}

// WhereIDBetween adds a BETWEEN condition on ID to the WHERE statement
func WhereIDBetween(low, high int64) common.Where {
	return common.NewWhereBetween("loaner", "id", low, high)
}

// WhereName adds a condition on Name to the WHERE statement
func WhereName(op common.Op, val string) common.Where {
	return common.NewWhere(op, "loaner", "name", val)
}

// WhereNameIn adds an IN condition on Name to the WHERE statement
func WhereNameIn(vals ...string) common.Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return common.NewWhereIn("loaner", "name", args...)
}

// WhereNameBetween adds a BETWEEN condition on Name to the WHERE statement
func WhereNameBetween(low, high string) common.Where {
	return common.NewWhereBetween("loaner", "name", low, high)
}

// WhereBook adds a condition on Book to the WHERE statement
func WhereBook(op common.Op, val *example.Book) common.Where {
	return common.NewWhere(op, "loaner", "book_id", val)
}

// WhereBookIn adds an IN condition on Book to the WHERE statement
func WhereBookIn(vals ...*example.Book) common.Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return common.NewWhereIn("loaner", "book_id", args...)
}

// WhereBookBetween adds a BETWEEN condition on Book to the WHERE statement
func WhereBookBetween(low, high *example.Book) common.Where {
	return common.NewWhereBetween("loaner", "book_id", low, high)
}