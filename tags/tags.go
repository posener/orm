package tags

import "strings"

// Parse parsers a struct tags
// A tag is a 2 levels deep mapping.
// * First level is separated by space, each item is of the pattern 'key:"value"' (the value is quoted)
// * Second level (value from first level) is separated by ';', each item is of the pattern 'key:value'
func Parse(tag string) map[string]map[string]string {
	if tag == "" {
		return nil
	}

	m := make(map[string]map[string]string)
	for key, val := range firstLevel(tag) {
		m[key] = make(map[string]string)
		for _, part := range strings.Split(val, ";") {
			key2, val2 := split(part)
			m[key][key2] = val2
		}
	}
	return m
}

func firstLevel(tag string) map[string]string {
	m := make(map[string]string)
	parts := Fields(tag)
	for _, part := range parts {
		key, val := split(part)
		m[key] = strings.Trim(val, `"`)
	}
	return m
}

// Fields splits a string on a space characters
// * Following spaces are treated as one space.
// * Quoted spaces are ignored.
// * Escaped quotes are not consider as quotes.
func Fields(s string) []string {
	var (
		fields  []string
		quoted  = false
		escaped = false
		begin   = -1
	)
	for i := range s {
		switch s[i] {
		case '"':
			if !escaped {
				quoted = !quoted
			}
		case '\\':
			escaped = !escaped
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
		default:
			escaped = false // if the last character was not escaped, we are no longer escaped
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
