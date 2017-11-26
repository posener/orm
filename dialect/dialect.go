package dialect

import (
	"fmt"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/sqlite3"
)

type Dialect interface {
	Name() string
	Create(*common.CreateParams) (string, []interface{})
	Insert(*common.InsertParams) (string, []interface{})
	Select(*common.SelectParams) (string, []interface{})
	Delete(*common.DeleteParams) (string, []interface{})
	Update(*common.UpdateParams) (string, []interface{})
}

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
	ConvertValueCode(field *common.Field) string
}

// NewGen returns all known DialectGenerators
func NewGen(tp common.Type) []Generator {
	return []Generator{
		&mysql.Gen{Tp: tp},
		&sqlite3.Gen{Tp: tp},
	}
}
