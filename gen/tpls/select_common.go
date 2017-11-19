package tpls

import "strings"

// Query returns an object to create a SELECT statement
func Query() *Select {
	return &Select{}
}

// Select is the struct that holds the SELECT data
type Select struct {
	columns []string
	where   *Where
	page    Page
}

// Where applies where conditions on the query
func (s *Select) Where(where *Where) *Select {
	s.where = where
	return s
}

// Limit applies rows limit on the query response
func (s *Select) Limit(limit int64) *Select {
	s.page.limit = limit
	return s
}

// Page applies rows offset and limit on the query response
func (s *Select) Page(offset, limit int64) *Select {
	s.page.offset = offset
	s.page.limit = limit
	return s
}

// selectString returns the columns to select for the SELECT statement
func (s *Select) selectString() string {
	if len(s.columns) == 0 {
		return "*"
	}
	return strings.Join(s.columns, ", ")
}

// add adds a column to the select statement
func (s *Select) add(column string) *Select {
	s.columns = append(s.columns, column)
	return s
}

// columnsMap is a map the maps column name to it's (list index + 1)
// if columnMap[column] == 0, the column does not exists in the select columns
func (s *Select) columnsMap() map[string]int {
	m := make(map[string]int, len(s.columns))
	for i, col := range s.columns {
		m[col] = i + 1
	}
	return m
}
