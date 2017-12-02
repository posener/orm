package load

import (
	"go/types"

	"encoding/json"

	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/tags"
)

const tagSQLType = "sql"

// SQL hold the SQL tags for a field in a struct
type SQL struct {
	// CustomType can be defined to a column
	CustomType sqltypes.Type `json:",omitempty"`
	// PrimaryKey defines a column as a table's primary key
	PrimaryKey bool `json:",omitempty"`
	// NotNull defines that this coulmn value can't be null
	NotNull bool `json:",omitempty"`
	// AutoIncrement defines this column as auto-increment column
	AutoIncrement bool `json:",omitempty"`
	// Unique defines that 2 rows can't have the same value of this column
	Unique bool `json:",omitempty"`
	// Default sets a default value for this column
	Default string `json:",omitempty"`
}

func newSQL(st *types.Struct, fieldIndex int) (*SQL, error) {
	var sql = new(SQL)
	err := sql.parseTags(st.Tag(fieldIndex))
	return sql, err
}

// parseTags parses tags from a struct tags into a SQL struct.
func (s *SQL) parseTags(tag string) error {
	if tag == "" {
		return nil
	}

	tagsMap := tags.Parse(tag)
	for key, value := range tagsMap[tagSQLType] {
		switch key {
		case "type":
			s.CustomType = sqltypes.Type(value)
		case "primary key", "primarykey", "primary_key":
			s.PrimaryKey = true
		case "not null", "not_null":
			s.NotNull = true
		case "auto increment", "autoincrement", "auto_increment":
			s.AutoIncrement = true
		case "unique":
			s.Unique = true
		case "default":
			s.Default = value
		}
	}
	return nil
}

func (s *SQL) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}
