package dialect

import (
	"fmt"

	"github.com/posener/orm"
)

// Where is an API for creating WHERE statements
type Where interface {
	Build(table string, b *builder)
	// Or applies OR condition with another Where statement
	Or(Where) Where
	// And applies AND condition with another Where statement
	And(Where) Where
}

// NewWhere returns a new WHERE statement
func NewWhere(op orm.Op, variable string, value interface{}) Where {
	return whereFunc(func(table string, b *builder) {
		b.Append(fmt.Sprintf("%s.%s", b.Quote(table), b.Quote(variable)))
		b.Append(string(op))
		b.Var(value)
	})
}

// NewWhereIn returns a new 'WHERE variable IN (...)' statement
func NewWhereIn(variable string, values ...interface{}) Where {
	return whereFunc(func(table string, b *builder) {
		b.Append(fmt.Sprintf("%s.%s", b.Quote(table), b.Quote(variable)))
		b.Append("IN")
		b.Open()
		for i, value := range values {
			b.Var(value)
			if i != len(values)-1 {
				b.Comma()
			}
		}
		b.Close()
	})
}

// NewWhereBetween returns a new 'WHERE variable BETWEEN low AND high' statement
func NewWhereBetween(variable string, low, high interface{}) Where {
	return whereFunc(func(table string, b *builder) {
		b.Append(fmt.Sprintf("%s.%s", b.Quote(table), b.Quote(variable)))
		b.Append("BETWEEN")
		b.Var(low)
		b.Append("AND")
		b.Var(high)
	})
}

// Where are options for SQL WHERE statement
type whereFunc func(table string, b *builder)

func (w whereFunc) Build(table string, b *builder) {
	w(table, b)
}

// Or applies an or condition between two where conditions
func (w whereFunc) Or(right Where) Where {
	if w == nil {
		return right
	}
	return binary(w, right, "OR")
}

// And applies an and condition between two where conditions
func (w whereFunc) And(right Where) Where {
	if w == nil {
		return right
	}
	return binary(w, right, "AND")
}

func binary(l Where, r Where, op string) Where {
	return whereFunc(func(table string, b *builder) {
		b.Open()
		l.Build(table, b)
		b.Close()
		b.Append(op)
		b.Open()
		r.Build(table, b)
		b.Close()
	})
}

func Or(w ...Where) Where {
	return multi("OR", w...)
}

func And(w ...Where) Where {
	return multi("AND", w...)
}

func multi(op string, w ...Where) Where {
	return whereFunc(func(table string, b *builder) {
		for i, wi := range w {
			b.Open()
			wi.Build(table, b)
			b.Close()
			if i != len(w)-1 {
				b.Append(op)
			}
		}
	})
}
