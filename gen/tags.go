package gen

import (
	"fmt"
	"go/types"
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
	Column string
	// Type matches the 'sql.type' tag: the SQL type of the field
	Type          string
	PrimaryKey    bool
	NotNull       bool
	AutoIncrement bool
	Unique        bool
	// Default value of column
	Default string
}

func newSQL(name string, st *types.Struct, fieldIndex int) (*SQL, error) {
	var sql = new(SQL)
	sql.Column = strings.ToLower(name)
	sql.parseTags(st.Tag(fieldIndex))
	field := st.Field(fieldIndex)
	fieldType := field.Type().String()
	if sql.Type == "" {
		sql.Type = defaultSQLTypes(fieldType)
	}
	if sql.Type == "" {
		return nil, fmt.Errorf("unsupported field type: %s", fieldType)
	}
	return sql, nil
}

// ConvertType is the type of the field when returned by sql/driver from database
func (s *SQL) ConvertType() string {
	switch s.typeFamily() {
	case "INT", "BIGINT", "INTEGER":
		return "int64"
	case "SMALLINT":
		return "int32"
	case "TINYINT":
		return "byte"
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

// parseTags parses tags from a struct tags into a SQL struct.
func (s *SQL) parseTags(tag string) {
	if tag == "" {
		return
	}

	tagsMap := tags.Parse(tag)
	for key, value := range tagsMap[tagSQLType] {
		switch key {
		case "type":
			s.Type = value
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

// typeFamily returns the family of the SQL type.
// for example, VARCHAR(255) will return only VARCHAR
func (s *SQL) typeFamily() string {
	prefixEnds := strings.Index(s.Type, "(")
	if prefixEnds == -1 {
		prefixEnds = len(s.Type)
	}
	return s.Type[:prefixEnds]
}
