package tpls

import "strings"

type Select []string

func NewSelect() Select {
	return Select{}
}

func (s Select) String() string {
	if len(s) == 0 {
		return "*"
	}
	return strings.Join(s, ", ")
}

func (s Select) add(column string) Select {
	return append(s, column)
}
