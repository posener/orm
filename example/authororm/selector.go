// Package authororm was auto-generated by github.com/posener/orm; DO NOT EDIT
package authororm

import (
	"database/sql/driver"
	"fmt"

	"github.com/posener/orm/example"

	"github.com/posener/orm/common"
)

const errMsg = "converting %s: column %d with value %v (type %T) to %s"

// selector selects columns for SQL queries and for parsing SQL rows
type selector struct {
	SelectID      bool
	SelectName    bool
	SelectHobbies bool
	JoinBooks     BooksScanner
	count         bool // used for sql COUNT(*) column
}

// Columns are the names of selected columns
func (s *selector) Columns() []string {
	var cols []string
	if s.SelectID {
		cols = append(cols, "id")
	}
	if s.SelectName {
		cols = append(cols, "name")
	}
	if s.SelectHobbies {
		cols = append(cols, "hobbies")
	}
	return cols
}

// Joins are join options of the query
func (s *selector) Joins() []common.Join {
	var joins []common.Join
	if selector := s.JoinBooks; selector != nil {
		joins = append(joins, common.Join{
			Column:        "id",
			RefTable:      "book",
			RefColumn:     "authorid",
			SelectColumns: selector.Columns(),
		})
	}

	return joins
}

// Count is true when a COUNT(*) column should be added to the query
func (s *selector) Count() bool {
	return s.count
}

// FirstCount scans an SQL row to a AuthorCount struct
func (s *selector) FirstCount(dialect string, vals []driver.Value) (*AuthorCount, error) {
	switch dialect {
	case "mysql":
		return s.scanmysql(vals)

	case "sqlite3":
		return s.scansqlite3(vals)
	default:
		return nil, fmt.Errorf("unsupported dialect %s", dialect)
	}
}

// First scans an SQL row to a Author struct
func (s *selector) First(dialect string, vals []driver.Value) (*example.Author, error) {
	item, err := s.FirstCount(dialect, vals)
	if err != nil {
		return nil, err
	}
	return &item.Author, nil
}

// scanmysql scans mysql row to a Author struct
func (s *selector) scanmysql(vals []driver.Value) (*AuthorCount, error) {
	var (
		row AuthorCount
		all = s.selectAll()
		i   int
	)

	if all || s.SelectID {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int64(parseInt(val))
				row.ID = tmp
			case int64:
				tmp := int64(val)
				row.ID = tmp
			default:
				return nil, fmt.Errorf(errMsg, "ID", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

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

	if all || s.SelectHobbies {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := string(val)
				row.Hobbies = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Hobbies", i, vals[i], vals[i], "[]byte, []byte")
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

	if j := s.JoinBooks; j != nil {
		tmp, err := j.First("mysql", vals[i:])
		if err != nil {
			return nil, err
		}
		row.Books = append(row.Books, tmp)
	}

	return &row, nil
}

// scansqlite3 scans sqlite3 row to a Author struct
func (s *selector) scansqlite3(vals []driver.Value) (*AuthorCount, error) {
	var (
		row AuthorCount
		all = s.selectAll()
		i   int
	)

	if all || s.SelectID {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "ID", i, vals[i], vals[i], "int64")
			}
			tmp := int64(val)
			row.ID = tmp
		}
		i++
	}

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

	if all || s.SelectHobbies {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Hobbies", i, vals[i], vals[i], "string")
			}
			tmp := string(val)
			row.Hobbies = tmp
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

	if j := s.JoinBooks; j != nil {
		tmp, err := j.First("sqlite3", vals[i:])
		if err != nil {
			return nil, err
		}
		row.Books = append(row.Books, tmp)
	}

	return &row, nil
}

// selectAll returns true if no column was specifically selected
func (s *selector) selectAll() bool {
	return !s.SelectID && !s.SelectName && !s.SelectHobbies && !s.count
}
