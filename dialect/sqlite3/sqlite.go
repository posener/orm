package sqlite3

import (
	"bytes"
	"fmt"

	"github.com/posener/orm"
	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/sqltypes"
)

func New(tp common.Type) orm.Dialect {
	return &sqlite3{tp: tp}
}

type sqlite3 struct {
	tp common.Type
}

// Insert returns an SQL INSERT statement and arguments
func Insert(i *orm.Insert) (string, []interface{}) {
	stmt := fmt.Sprintf(`INSERT INTO '%s' (%s) VALUES (%s)`,
		i.Table,
		assignColumns(i.Assignments),
		common.QMarks(len(i.Assignments)),
	)

	var args []interface{}
	if i.Assignments != nil {
		args = append(args, i.Assignments.Args()...)
	}

	return stmt, args
}

// Select returns an SQL SELECT statement and arguments
func Select(s *orm.Select) (string, []interface{}) {
	stmt := fmt.Sprintf("SELECT %s FROM '%s' %s %s %s %s",
		columns(s.Columns),
		s.Table,
		whereStatement(s.Where),
		groups(s.Groups),
		orders(s.Orders),
		page(s.Page),
	)

	var args []interface{}
	if s.Where != nil {
		args = append(args, s.Where.Args()...)
	}

	return stmt, args
}

// Delete returns an SQL DELETE statement and arguments
func Delete(d *orm.Delete) (string, []interface{}) {
	stmt := fmt.Sprintf("DELETE FROM '%s' %s", d.Table, whereStatement(d.Where))

	var args []interface{}
	if d.Where != nil {
		args = append(args, d.Where.Args()...)
	}

	return stmt, args
}

// Update returns an SQL UPDATE statement and arguments
func Update(u *orm.Update) (string, []interface{}) {
	stmt := fmt.Sprintf(`UPDATE '%s' SET %s %s`, u.Table, assignSets(u.Assignments), whereStatement(u.Where))

	var args []interface{}
	if u.Assignments != nil {
		args = append(args, u.Assignments.Args()...)
	}
	if u.Where != nil {
		args = append(args, u.Where.Args()...)
	}

	return stmt, args
}

func defaultSQLTypes(tp string) sqltypes.Type {
	switch tp {
	case "int", "int32", "int64":
		return sqltypes.Integer
	case "float", "float32", "float64":
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

func (s *sqlite3) Name() string {
	return "sqlite3"
}

func sqlType(f *common.Field) sqltypes.Type {
	if f.SQL.CustomType != "" {
		return f.SQL.CustomType
	}
	return defaultSQLTypes(f.Type)
}

func columns(c orm.Columner) string {
	if c == nil {
		return "*"
	}
	cols := c.Columns()
	if len(cols) == 0 && !c.Count() {
		return "*"
	}
	b := bytes.NewBuffer(nil)
	for i := range cols {
		b.WriteString("`" + cols[i] + "`, ")
	}

	if c.Count() {
		b.WriteString("COUNT(*), ")
	}

	s := b.String()
	return s[:len(s)-2]
}

func whereStatement(c orm.StatementArger) string {
	if c == nil {
		return ""
	}
	where := c.Statement()
	if len(where) == 0 {
		return ""
	}
	return "WHERE " + where
}

func groups(groups []orm.Group) string {
	if len(groups) == 0 {
		return ""
	}
	b := bytes.NewBufferString("GROUP BY ")
	for i := range groups {
		b.WriteString(fmt.Sprintf("`%s`, ", groups[i].Column))
	}

	s := b.String()
	return s[:len(s)-2]
}

func orders(orders []orm.Order) string {
	if len(orders) == 0 {
		return ""
	}

	b := bytes.NewBufferString("ORDER BY ")
	for i := range orders {
		b.WriteString(fmt.Sprintf("`%s` %s, ", orders[i].Column, orders[i].Dir))
	}

	s := b.String()
	return s[:len(s)-2]
}

func page(p orm.Page) string {
	if p.Limit == 0 { // why would someone ask for a page of zero size?
		return ""
	}
	stmt := fmt.Sprintf("LIMIT %d", p.Limit)
	if p.Offset != 0 {
		stmt += fmt.Sprintf(" OFFSET %d", p.Offset)
	}
	return stmt
}

func assignColumns(a orm.Assignments) string {
	if len(a) == 0 {
		return ""
	}
	b := bytes.NewBuffer(nil)
	for i := range a {
		b.WriteString("`" + a[i].Column + "`, ")
	}

	s := b.String()
	return s[:len(s)-2]
}

func assignSets(a orm.Assignments) string {
	if len(a) == 0 {
		return ""
	}
	b := bytes.NewBuffer(nil)
	for i := range a {
		b.WriteString("`" + a[i].Column + "` = ?, ")
	}

	s := b.String()
	return s[:len(s)-2]
}
