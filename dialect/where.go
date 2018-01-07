package dialect

import (
	"fmt"
	"strings"

	"github.com/posener/orm"
)

// StatementArger is interface for queries.
// The statement and the args are given to the SQL query.
type StatementArger interface {
	Statement(table string, d Dialect) string
	Args() []interface{}
}

// Where is an API for creating WHERE statements
type Where interface {
	StatementArger
	// Or applies OR condition with another Where statement
	Or(Where) Where
	// And applies AND condition with another Where statement
	And(Where) Where
}

// Where are options for SQL WHERE statement
type where struct {
	stmtr func(table string, d Dialect) string
	args  []interface{}
}

// NewWhere returns a new WHERE statement
func NewWhere(op orm.Op, variable string, value interface{}) Where {
	return &where{
		stmtr: func(table string, d Dialect) string {
			return fmt.Sprintf("%s.%s %s ?", d.Quote(table), d.Quote(variable), op)
		},
		args: []interface{}{value},
	}
}

// NewWhereIn returns a new 'WHERE variable IN (...)' statement
func NewWhereIn(variable string, values ...interface{}) Where {
	return &where{
		stmtr: func(table string, d Dialect) string {
			return fmt.Sprintf("%s.%s IN (%s)", d.Quote(table), d.Quote(variable), QMarks(len(values)))
		},
		args: values,
	}
}

// NewWhereBetween returns a new 'WHERE variable BETWEEN low AND high' statement
func NewWhereBetween(variable string, low, high interface{}) Where {
	return &where{
		stmtr: func(table string, d Dialect) string {
			return fmt.Sprintf("%s.%s BETWEEN ? AND ?", d.Quote(table), d.Quote(variable))
		},
		args: []interface{}{low, high},
	}
}

// String returns the WHERE SQL statement
func (w *where) Statement(table string, d Dialect) string {
	if w == nil || w.stmtr == nil {
		return ""
	}
	return w.stmtr(table, d)
}

// Or applies an or condition between two where conditions
func (w *where) Or(right Where) Where {
	if w == nil {
		return right
	}
	return binary(w, right, "OR")
}

// And applies an and condition between two where conditions
func (w *where) And(right Where) Where {
	if w == nil {
		return right
	}
	return binary(w, right, "AND")
}

// Args are the arguments for executing a SELECT query with a WHERE condition
func (w *where) Args() []interface{} {
	if w == nil {
		return nil
	}
	return w.args
}

func binary(l *where, r StatementArger, op string) Where {
	return &where{
		stmtr: func(table string, d Dialect) string {
			return fmt.Sprintf("(%s) %s (%s)", l.Statement(table, d), op, r.Statement(table, d))
		},
		args: append(l.Args(), r.Args()...),
	}
}

// QMarks is a helper function for concatenating question mark for an SQL statement
func QMarks(n int) string {
	if n == 0 {
		return ""
	}
	qMark := strings.Repeat("?, ", n)
	qMark = qMark[:len(qMark)-2] // remove last ", "
	return qMark
}
