package tpls

import "strings"

type groupBy []string

// String is the SQL representation of GROUP BY
func (g groupBy) String() string {
	if len(g) == 0 {
		return ""
	}
	return "GROUP BY " + strings.Join(g, ", ")
}

// Add adds a column to the grouping
func (g *groupBy) add(column string) {
	*g = append(*g, column)
}
