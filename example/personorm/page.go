// Autogenerated by github.com/posener/orm
package personorm

import "fmt"

// Page represents an SQL LIMIT statement
type Page struct {
	limit  int64
	offset int64
}

// String returns the SQL query representation of the Page
func (p *Page) String() string {
	if p.limit == 0 && p.offset == 0 {
		return ""
	}
	if p.offset == 0 {
		return fmt.Sprintf("LIMIT %d", p.limit)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.limit, p.offset)
}
