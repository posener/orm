package gen

import "strings"

// Tags hold the SQL tags for a field in a struct
type Tags struct {
	// Type matches the 'sql.type' tag
	Type string
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
		}
	}

	return t
}

func getSqlField(tag string) string {
	parts := strings.Fields(tag)
	for _, part := range parts {
		key, val := split(part)
		if key != tagSQLType {
			continue
		}
		return strings.Trim(val, `"`)
	}
	return ""
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
