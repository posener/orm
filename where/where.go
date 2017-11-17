package where

import "strings"

type Options struct {
	stmt []string
	args []interface{}
}

func New(op Op, variable string, value interface{}) Options {
	var o Options
	o.stmt = append(o.stmt, variable+string(op)+"?")
	o.args = append(o.args, value)
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
