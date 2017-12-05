package sqlite3

import "fmt"

// Dialect represents the sqlite3 dialect
type Dialect struct{}

// Name returns the name of the dialect
func (*Dialect) Name() string {
	return "sqlite3"
}

// TableQuote takes table name and return the name quoted, according to the syntax
func (*Dialect) TableQuote(name string) string {
	return fmt.Sprintf("'%s'", name)
}
