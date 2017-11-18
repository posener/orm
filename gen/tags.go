package gen

import (
	"strings"
)

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

	sqlField := getSqlField(tag)
	for _, part := range strings.Split(sqlField, ";") {
		key, value := split(part)
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

func getSqlField(tag string) string {
	parts := fields(tag)
	for _, part := range parts {
		key, val := split(part)
		if key != tagSQLType {
			continue
		}
		return strings.Trim(val, `"`)
	}
	return ""
}

func fields(s string) []string {
	var (
		fields []string
		quoted = false
		begin  = -1
	)
	for i := range s {
		switch s[i] {
		case '"':
			quoted = !quoted
		case ' ':
			// if quoted, do not split
			if quoted {
				continue
			}
			// double space
			if i > begin+1 {
				fields = append(fields, s[begin+1:i])
			}
			begin = i
		}
	}
	if len(s) > begin+1 {
		fields = append(fields, s[begin+1:])
	}
	return fields
}

// split splits key:value to a (key, value)
func split(joined string) (string, string) {
	parts := strings.SplitN(joined, ":", 2)
	key := parts[0]
	value := ""
	if len(parts) == 2 {
		value = parts[1]
	}
	return key, value
}
