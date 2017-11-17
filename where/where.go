package where

import (
	"fmt"
	"log"
	"strings"
)

type Options struct {
	stmt []string
	args []interface{}
}

func New(op Op, variable string, value interface{}) Options {
	switch op {
	case OpEq, OpNe, OpGt, OpGE, OpLt, OpLE, OpLike:
	default:
		log.Panicf("Operation %s is not defined for one value", op)
	}
	var o Options
	o.stmt = append(o.stmt, fmt.Sprintf("%s %s ?", variable, op))
	o.args = append(o.args, value)
	return o
}

func NewMul(op Op, variable string, values ...interface{}) Options {
	var o Options
	switch op {
	case OpBetween:
		if len(values) != 2 {
			log.Panicf("Operation %s accepts only 2 values, got %d values", op, len(values))
		}
		o.stmt = append(o.stmt, fmt.Sprintf("%s %s ? AND ?", variable, op))
	case OpIn:
		if len(values) > 0 {
			qMark := strings.Repeat("?, ", len(values))
			qMark = qMark[:len(qMark)-2] // remove last ", "
			o.stmt = append(o.stmt, fmt.Sprintf("%s %s (%s)", variable, op, qMark))
		}
	default:
		log.Panicf("Operation %s does not support multiple values", op)
	}
	o.args = append(o.args, values...)
	return o
}

func (o *Options) String() string {
	if o == nil || len(o.stmt) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(o.stmt, " ")
}

func (o *Options) Args() []interface{} {
	if o == nil {
		return nil
	}
	return o.args
}

func (o Options) Or(right Options) Options {
	return binary(o, right, "OR")
}

func (o Options) And(right Options) Options {
	return binary(o, right, "AND")
}

func binary(l Options, r Options, op string) Options {
	l.stmt = append([]string{"("}, l.stmt...)
	l.stmt = append(l.stmt, ")", op, "(")
	l.stmt = append(l.stmt, r.stmt...)
	l.stmt = append(l.stmt, ")")
	l.args = append(l.args, r.args...)
	return l
}

type String struct {
	Val string
}

type Int struct {
	Val int
}
