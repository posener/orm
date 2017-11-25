package orm

import (
	"fmt"
	"strings"

	"github.com/posener/orm/common"
)

// Where are options for SQL WHERE statement
type where struct {
	stmt []string
	args []interface{}
}

// NewWhere returns a new WHERE statement
func NewWhere(op Op, variable string, value interface{}) Where {
	var w where
	w.stmt = append(w.stmt, fmt.Sprintf("`%s` %s ?", variable, op))
	w.args = append(w.args, value)
	return &w
}

// NewWhereIn returns a new 'WHERE variable IN (...)' statement
func NewWhereIn(variable string, values ...interface{}) Where {
	var w where
	w.stmt = append(w.stmt, fmt.Sprintf("`%s` IN (%s)", variable, common.QMarks(len(values))))
	w.args = append(w.args, values...)
	return &w
}

// NewWhereBetween returns a new 'WHERE variable BETWEEN low AND high' statement
func NewWhereBetween(variable string, low, high interface{}) Where {
	var w where
	w.stmt = append(w.stmt, fmt.Sprintf("`%s` BETWEEN ? AND ?", variable))
	w.args = append(w.args, low, high)
	return &w
}

// String returns the WHERE SQL statement
func (w *where) Statement() string {
	if w == nil || len(w.stmt) == 0 {
		return ""
	}
	return strings.Join(w.stmt, " ")
}

// Or applies an or condition between two where conditions
func (w *where) Or(right Where) Where {
	return binary(w, right, "OR")
}

// And applies an and condition between two where conditions
func (w *where) And(right Where) Where {
	return binary(w, right, "AND")
}

// Args are the arguments for executing a SELECT query with a WHERE condition
func (w *where) Args() []interface{} {
	if w == nil {
		return nil
	}
	return w.args
}

func binary(l *where, r Where, op string) Where {
	l.stmt = []string{"(", l.Statement(), ")", op, "(", r.Statement(), ")"}
	l.args = append(l.args, r.Args()...)
	return l
}
