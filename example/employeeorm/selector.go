// Package employeeorm was auto-generated by github.com/posener/orm; DO NOT EDIT
package employeeorm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"unsafe"
)

const errMsg = "converting %s: column %d with value %v (type %T) to %s"

// selector selects columns for SQL queries and for parsing SQL rows
type selector struct {
	SelectName   bool
	SelectAge    bool
	SelectSalary bool

	count bool // used for sql COUNT(*) column
}

// Columns are the names of selected columns
func (s *selector) Columns() []string {
	var cols []string
	if s.SelectName {
		cols = append(cols, "name")
	}
	if s.SelectAge {
		cols = append(cols, "age")
	}
	if s.SelectSalary {
		cols = append(cols, "salary")
	}

	return cols
}

// Count is true when a COUNT(*) column should be added to the query
func (s *selector) Count() bool {
	return s.count
}

// scan an SQL row to a EmployeeCount struct
func (s *selector) scan(dialect string, rows *sql.Rows) (*EmployeeCount, error) {
	switch dialect {
	case "mysql":
		return s.scanmysql(rows)

	case "sqlite3":
		return s.scansqlite3(rows)
	default:
		return nil, fmt.Errorf("unsupported dialect %s", dialect)
	}
}

// scanmysql scans mysql row to a Employee struct
func (s *selector) scanmysql(rows *sql.Rows) (*EmployeeCount, error) {
	var (
		vals = values(*rows)
		row  EmployeeCount
		all  = s.selectAll()
		i    = 0
	)

	if all || s.SelectName {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := string(val)
				row.Name = tmp
			default:
				return nil, fmt.Errorf(errMsg, "string", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectAge {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int(parseInt(val))
				row.Age = tmp
			case int64:
				tmp := int(val)
				row.Age = tmp
			default:
				return nil, fmt.Errorf(errMsg, "int", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectSalary {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int(parseInt(val))
				row.Salary = tmp
			case int64:
				tmp := int(val)
				row.Salary = tmp
			default:
				return nil, fmt.Errorf(errMsg, "int", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if s.count {
		switch val := vals[i].(type) {
		case int64:
			row.Count = val
		case []byte:
			row.Count = parseInt(val)
		default:
			return nil, fmt.Errorf(errMsg, "COUNT(*)", i, vals[i], vals[i], "int64, []byte")
		}
	}

	return &row, nil
}

// scansqlite3 scans sqlite3 row to a Employee struct
func (s *selector) scansqlite3(rows *sql.Rows) (*EmployeeCount, error) {
	var (
		vals = values(*rows)
		row  EmployeeCount
		all  = s.selectAll()
		i    = 0
	)

	if all || s.SelectName {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "string", i, vals[i], vals[i], "string")
			}
			tmp := string(val)
			row.Name = tmp
		}
		i++
	}

	if all || s.SelectAge {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "int", i, vals[i], vals[i], "int")
			}
			tmp := int(val)
			row.Age = tmp
		}
		i++
	}

	if all || s.SelectSalary {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "int", i, vals[i], vals[i], "int")
			}
			tmp := int(val)
			row.Salary = tmp
		}
		i++
	}

	if s.count {
		switch val := vals[i].(type) {
		case int64:
			row.Count = val
		case []byte:
			row.Count = parseInt(val)
		default:
			return nil, fmt.Errorf(errMsg, "COUNT(*)", i, vals[i], vals[i], "int64, []byte")
		}
	}

	return &row, nil
}

// selectAll returns true if no column was specifically selected
func (s *selector) selectAll() bool {
	return !s.SelectName && !s.SelectAge && !s.SelectSalary && !s.count
}

// values is a hack to the sql.Rows struct
// since the rows struct does not expose it's lastcols values, or a way to give
// a custom scanner to the Scan method.
// See issue https://github.com/golang/go/issues/22544
func values(r sql.Rows) []driver.Value {
	// some ugly hack to access lastcols field
	rs := reflect.ValueOf(&r).Elem()
	rf := rs.FieldByName("lastcols")

	// overcome panic reflect.Value.Interface: cannot return value obtained from unexported field or method
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	return rf.Interface().([]driver.Value)
}
