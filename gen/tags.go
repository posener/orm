package gen

import (
	"github.com/posener/orm/tags"
)

const tagSQLType = "sql"

var defaultSQLTypes = map[string]string{
	"string": "VARCHAR(255)",
	"int":    "INT",
	"bool":   "BOOLEAN",
	"float":  "REAL",
}

// Tags hold the SQL tags for a field in a struct
type Tags struct {
	// Type matches the 'sql.type' tag: the SQL type of the field
	Type string
	// PrimaryKey matches the 'sql.primary key' tag: the field is the primary key of the struct
	PrimaryKey bool
	// PrimaryKey matches the 'sql.not null' tag: the field is of type "NOT NULL"
	NotNull bool
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
		}
	}

	return t
}
