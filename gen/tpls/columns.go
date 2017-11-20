package tpls

import "strings"

type columns []string

// String is the SQL representation of columns
func (c columns) String() string {
	if len(c) == 0 {
		return "*"
	}
	return strings.Join(c, ", ")
}

// Add adds a column to the statement
func (c *columns) add(column string) {
	*c = append(*c, column)
}

// indexMap maps column name to it's (list index + 1)
// if columnMap[column] == 0, the column does not exists in the select columns
func (c columns) indexMap() map[string]int {
	m := make(map[string]int, len(c))
	for i, col := range c {
		m[col] = i + 1
	}
	return m
}
