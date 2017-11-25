package sqlite3

import (
	"fmt"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/format"
	"github.com/posener/orm/dialect/sqltypes"
)

func New(tp common.Type) common.Dialect {
	return &sqlite3{tp: tp}
}

type sqlite3 struct {
	tp common.Type
}

func (s *sqlite3) Name() string {
	return "sqlite3"
}

// ConvertType is the type of the field when returned by sql/driver from database
func (s *sqlite3) ConvertType(f *common.Field) string {
	switch sqltypes.Family(s.sqlType(f)) {
	case sqltypes.Integer:
		return "int64"
	case sqltypes.Float:
		return "float64"
	case sqltypes.Text, sqltypes.Blob, sqltypes.VarChar:
		return "[]byte"
	case sqltypes.Boolean:
		return "bool"
	default:
		return f.NonPointerType()
	}
}

func (sqlite3) sqlType(f *common.Field) sqltypes.Type {
	if f.SQL.CustomType != "" {
		return f.SQL.CustomType
	}
	switch f.NonPointerType() {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return sqltypes.Integer
	case "float", "float8", "float16", "float32", "float64":
		return sqltypes.Float
	case "bool":
		return sqltypes.Boolean
	case "string":
		return sqltypes.Text
	case "[]byte":
		return sqltypes.Blob
	case "time.Time":
		return sqltypes.TimeStamp
	default:
		return sqltypes.NA
	}
}

// Insert returns an SQL CREATE statement and arguments
func Create(c *common.Create) (string, []interface{}) {
	stmt := fmt.Sprintf(`CREATE TABLE %s '%s' ( %s )`,
		format.IfNotExists(c.IfNotExists),
		c.Table,
		c.ColumnsStatement,
	)

	return stmt, nil
}

// Insert returns an SQL INSERT statement and arguments
func Insert(i *common.Insert) (string, []interface{}) {
	stmt := fmt.Sprintf(`INSERT INTO '%s' (%s) VALUES (%s)`,
		i.Table,
		format.AssignColumns(i.Assignments),
		common.QMarks(len(i.Assignments)),
	)

	var args []interface{}
	if i.Assignments != nil {
		args = append(args, i.Assignments.Args()...)
	}

	return stmt, args
}

// Select returns an SQL SELECT statement and arguments
func Select(s *common.Select) (string, []interface{}) {
	stmt := fmt.Sprintf("SELECT %s FROM '%s' %s %s %s %s",
		format.Columns(s.Columns),
		s.Table,
		format.Where(s.Where),
		format.GroupBy(s.Groups),
		format.OrderBy(s.Orders),
		format.Page(s.Page),
	)

	var args []interface{}
	if s.Where != nil {
		args = append(args, s.Where.Args()...)
	}

	return stmt, args
}

// Delete returns an SQL DELETE statement and arguments
func Delete(d *common.Delete) (string, []interface{}) {
	stmt := fmt.Sprintf("DELETE FROM '%s' %s",
		d.Table,
		format.Where(d.Where),
	)

	var args []interface{}
	if d.Where != nil {
		args = append(args, d.Where.Args()...)
	}

	return stmt, args
}

// Update returns an SQL UPDATE statement and arguments
func Update(u *common.Update) (string, []interface{}) {
	stmt := fmt.Sprintf(`UPDATE '%s' SET %s %s`,
		u.Table,
		format.AssignSets(u.Assignments),
		format.Where(u.Where),
	)

	var args []interface{}
	if u.Assignments != nil {
		args = append(args, u.Assignments.Args()...)
	}
	if u.Where != nil {
		args = append(args, u.Where.Args()...)
	}

	return stmt, args
}
