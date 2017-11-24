package gen

import (
	"strings"

	"github.com/posener/orm/tags"
)

const tagSQLType = "sql"

func defaultSQLTypes(tp string) string {
	switch tp {
	case "string":
		return "VARCHAR(255)"
	case "int", "int32":
		return "INT"
	case "int64":
		return "BIGINT"
	case "float", "float32":
		return "FLOAT"
	case "float64":
		return "DOUBLE"
	case "bool":
		return "BOOLEAN"
	case "time.Time":
		return "TIMESTAMP"
	default:
		return ""
	}
}

// SQL hold the SQL tags for a field in a struct
type SQL struct {
	// Type matches the 'sql.type' tag: the SQL type of the field
	Type          string
	PrimaryKey    bool
	NotNull       bool
	AutoIncrement bool
	Unique        bool
	// Default value of column
	Default string
}

// ConvertType is the type of the field when returned by sql/driver from database
func (s *SQL) ConvertType() string {
	switch s.typeFamily() {
	case "INT", "BIGINT", "INTEGER":
		return "int64"
	case "VARCHAR":
		return "[]byte"
	case "BOOLEAN":
		return "bool"
	case "DOUBLE":
		return "float64"
	case "FLOAT":
		return "float32"
	case "TIMESTAMP":
		return "time.Time"
	default:
		return "interface{}"
	}
}

// ParseTags parses tags from a struct tags into a SQL struct.
func ParseTags(tag string) SQL {
	var t SQL
	if tag == "" {
		return t
	}

	tagsMap := tags.Parse(tag)
	for key, value := range tagsMap[tagSQLType] {
		switch key {
		case "type":
			t.Type = value
		case "primary_key", "primary key":
			t.PrimaryKey = true
		case "not null", "not_null":
			t.NotNull = true
		case "auto_increment", "auto increment":
			t.AutoIncrement = true
		case "unique":
			t.Unique = true
		case "default":
			t.Default = value
		}
	}

	return t
}

// typeFamily returns the family of the SQL type.
// for example, VARCHAR(255) will return only VARCHAR
func (s *SQL) typeFamily() string {
	prefixEnds := strings.Index(s.Type, "(")
	if prefixEnds == -1 {
		prefixEnds = len(s.Type)
	}
	return s.Type[:prefixEnds]
}
