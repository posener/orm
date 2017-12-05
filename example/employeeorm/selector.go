// Package employeeorm was auto-generated by github.com/posener/orm; DO NOT EDIT
package employeeorm

import (
	"database/sql/driver"
	"fmt"

	"github.com/posener/orm/example"

	"github.com/posener/orm/common"
)

const errMsg = "converting %s: column %d with value %v (type %T) to %s"

// selector selects columns for SQL queries and for parsing SQL rows
type selector struct {
	SelectName   bool
	SelectAge    bool
	SelectSalary bool
	count        bool // used for sql COUNT(*) column
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

// Joins are join options of the query
func (s *selector) Joins() []common.JoinParams {
	var joins []common.JoinParams

	return joins
}

// Count is true when a COUNT(*) column should be added to the query
func (s *selector) Count() bool {
	return s.count
}

// FirstCount scans an SQL row to a EmployeeCount struct
func (s *selector) FirstCount(dialect string, vals []driver.Value) (*EmployeeCount, error) {
	switch dialect {
	case "mysql":
		return s.scanmysql(vals)

	case "sqlite3":
		return s.scansqlite3(vals)
	default:
		return nil, fmt.Errorf("unsupported dialect %s", dialect)
	}
}

// First scans an SQL row to a Employee struct
func (s *selector) First(dialect string, vals []driver.Value) (*example.Employee, error) {
	item, err := s.FirstCount(dialect, vals)
	if err != nil {
		return nil, err
	}
	return &item.Employee, nil
}

// scanmysql scans mysql row to a Employee struct
func (s *selector) scanmysql(vals []driver.Value) (*EmployeeCount, error) {
	var (
		row EmployeeCount
		all = s.selectAll()
		i   int
	)

	if all || s.SelectName {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := string(val)
				row.Name = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Name", i, vals[i], vals[i], "[]byte, []byte")
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
				return nil, fmt.Errorf(errMsg, "Age", i, vals[i], vals[i], "[]byte, int64")
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
				return nil, fmt.Errorf(errMsg, "Salary", i, vals[i], vals[i], "[]byte, int64")
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
		i++
	}

	return &row, nil
}

// scansqlite3 scans sqlite3 row to a Employee struct
func (s *selector) scansqlite3(vals []driver.Value) (*EmployeeCount, error) {
	var (
		row EmployeeCount
		all = s.selectAll()
		i   int
	)

	if all || s.SelectName {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Name", i, vals[i], vals[i], "string")
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
				return nil, fmt.Errorf(errMsg, "Age", i, vals[i], vals[i], "int")
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
				return nil, fmt.Errorf(errMsg, "Salary", i, vals[i], vals[i], "int")
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
		i++
	}

	return &row, nil
}

// selectAll returns true if no column was specifically selected
func (s *selector) selectAll() bool {
	return !s.SelectName && !s.SelectAge && !s.SelectSalary && !s.count
}
