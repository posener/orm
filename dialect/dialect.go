package dialect

import (
	"fmt"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/sqlite3"
	"github.com/posener/orm/load"
)

// Dialect is an interface to represent an SQL dialect
// Objects that implement this interface, can convert query params, such as SelectParams or
// UpdateParams, and convert them to an SQL statement and a list of arguments, which can be used
// to invoke SQL Exec or Query functions.
type Dialect interface {
	// Name returns the name of the dialect
	Name() string
	// Create returns the SQL CREATE statement and arguments according to the given parameters
	Create(*common.CreateParams) (string, []interface{})
	// Insert returns the SQL INSERT statement and arguments according to the given parameters
	Insert(*common.InsertParams) (string, []interface{})
	// Select returns the SQL SELECT statement and arguments according to the given parameters
	Select(*common.SelectParams) (string, []interface{})
	// Delete returns the SQL DELETE statement and arguments according to the given parameters
	Delete(*common.DeleteParams) (string, []interface{})
	// Update returns the SQL UPDATE statement and arguments according to the given parameters
	Update(*common.UpdateParams) (string, []interface{})
}

// New returns a new Dialect according to it's name
func New(name string) (Dialect, error) {
	switch name {
	case "mysql":
		return new(mysql.Dialect), nil
	case "sqlite3":
		return new(sqlite3.Dialect), nil
	default:
		return nil, fmt.Errorf("unsupported dialect %s", name)
	}
}

// Generator is API for different dialects
type Generator interface {
	// Name is the dialect name
	Name() string
	// ColumnsStatement returns the fields parts of SQL CREATE TABLE statement
	// for a specific struct and specific dialect.
	// It is used by the generation tool.
	ColumnsStatement() string
	// ConvertValueCode returns go code for converting value returned from the
	// database to the given field.
	ConvertValueCode(field *load.Field) string
}

// NewGen returns all known Generators
func NewGen(tp *load.Type) []Generator {
	return []Generator{
		&mysql.Gen{Tp: tp},
		&sqlite3.Gen{Tp: tp},
	}
}
