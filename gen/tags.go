package gen

import (
	"github.com/posener/orm/tags"
)

const tagSQLType = "sql"

var defaultSQLTypes = map[string]string{
	"string":  "VARCHAR(255)",
	"int":     "BIGINT",
	"int32":   "INT",
	"int64":   "BIGINT",
	"float":   "DOUBLE",
	"float32": "FLOAT",
	"float64": "DOUBLE",
	"bool":    "BOOLEAN",
}

// Tags hold the SQL tags for a field in a struct
type Tags struct {
	// Type matches the 'sql.type' tag: the SQL type of the field
	Type string

	PrimaryKey    bool
	NotNull       bool
	AutoIncrement bool
	Unique        bool
	// Default value of column
	Default string
}

// ParseTags parses tags from a struct tags into a Tags struct.
func ParseTags(tag string) Tags {
	var t Tags
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
