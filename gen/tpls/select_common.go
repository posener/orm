package tpls

import "strings"

// Select returns an object to create a SELECT statement
func Select() TSelect {
	return TSelect{}
}

// TSelect is the struct that holds the SELECT data
type TSelect []string

func (s TSelect) String() string {
	if len(s) == 0 {
		return "*"
	}
	return strings.Join(s, ", ")
}

func (s TSelect) add(column string) TSelect {
	return append(s, column)
}
