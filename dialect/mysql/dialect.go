package mysql

import "fmt"

// Dialect represents the mysql dialect
type Dialect struct{}

// Name returns the name of the dialect
func (*Dialect) Name() string {
	return "mysql"
}

// TableQuote takes table name and return the name quoted, according to the syntax
func (*Dialect) TableQuote(name string) string {
	return fmt.Sprintf("`%s`", name)
}
