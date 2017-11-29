package load

import (
	"go/types"
	"strings"

	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/tags"
)

const tagSQLType = "sql"

// SQL hold the SQL tags for a field in a struct
type SQL struct {
	Column        string
	CustomType    sqltypes.Type
	PrimaryKey    bool
	NotNull       bool
	AutoIncrement bool
	Unique        bool
	Default       string // Default value of column
}

// Auto returns whether the column should be set
func (s *SQL) Auto() bool {
	return s.PrimaryKey || s.AutoIncrement
}

func newSQL(name string, st *types.Struct, fieldIndex int) (*SQL, error) {
	var sql = new(SQL)
	sql.Column = strings.ToLower(name)
	sql.parseTags(st.Tag(fieldIndex))
	return sql, nil
}

// parseTags parses tags from a struct tags into a SQL struct.
func (s *SQL) parseTags(tag string) {
	if tag == "" {
		return
	}

	tagsMap := tags.Parse(tag)
	for key, value := range tagsMap[tagSQLType] {
		switch key {
		case "type":
			s.CustomType = sqltypes.Type(value)
		case "primary_key", "primary key":
			s.PrimaryKey = true
		case "not null", "not_null":
			s.NotNull = true
		case "auto_increment", "auto increment", "autoincrement":
			s.AutoIncrement = true
		case "unique":
			s.Unique = true
		case "default":
			s.Default = value
		}
	}
}
