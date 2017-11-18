// Autogenerated by github.com/posener/orm
package personorm

import "github.com/posener/orm/example"

// Name Add Name to the selected column of a query
func (s Select) Name() Select {
	return append(s, "name")
}

// Age Add Age to the selected column of a query
func (s Select) Age() Select {
	return append(s, "age")
}

// scanArgs are list of fields to be given to the sql Scan command
func (s Select) scanArgs(p *example.Person) []interface{} {
	if len(s) == 0 {
		// add to args all the fields of p
		return []interface{}{
			&p.Name,
			&p.Age,
		}
	}

	// select was given, choose only some fields
	m := make(map[string]int, len(s))
	for i, col := range s {
		m[col] = i + 1
	}
	args := make([]interface{}, len(s))
	if i := m["name"]; i != 0 {
		args[i-1] = &p.Name
	}
	if i := m["age"]; i != 0 {
		args[i-1] = &p.Age
	}
	return args
}
