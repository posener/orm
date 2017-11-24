package sqlite3

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/posener/orm/def"
	"github.com/posener/orm/dialect/api"
	"github.com/posener/orm/dialect/sqltypes"
)

func New(tp def.Type) api.Dialect {
	return &sqlite3{tp: tp}
}

type sqlite3 struct {
	tp def.Type
}

func Insert(tabler api.Tabler, assignments []api.Assignment) string {
	return fmt.Sprintf(`INSERT INTO '%s' (%s) VALUES (%s)`,
		tabler.Table(),
		assignColumns(assignments),
		qMarks(len(assignments)),
	)
}

// Select returns SQL SELECT string
func Select(tabler api.Tabler, columner api.Columner, wherer api.Wherer, g []api.Group, o []api.Order, pager api.Pager) string {
	return fmt.Sprintf("SELECT %s FROM '%s' %s %s %s %s",
		columns(columner),
		tabler.Table(),
		where(wherer),
		groups(g),
		orders(o),
		page(pager),
	)
}

// Delete returns SQL DELETE query string
func Delete(tabler api.Tabler, wherer api.Wherer) string {
	return fmt.Sprintf("DELETE FROM '%s' %s", tabler.Table(), where(wherer))
}

// Update returns SQL UPDATE query string
func Update(tabler api.Tabler, assignments []api.Assignment, wherer api.Wherer) string {
	return fmt.Sprintf(`UPDATE '%s' SET %s %s`, tabler.Table(), assignSets(assignments), where(wherer))
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

func sqlType(f *def.Field) sqltypes.Type {
	if f.SQL.CustomType != "" {
		return f.SQL.CustomType
	}
	return defaultSQLTypes(f.Type)
}

func columns(c api.Columner) string {
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

func where(c api.Wherer) string {
	if c == nil {
		return ""
	}
	where := c.Where()
	if len(where) == 0 {
		return ""
	}
	return "WHERE " + where
}

func groups(groups []api.Group) string {
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

func orders(orders []api.Order) string {
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

func page(p api.Pager) string {
	if p == nil {
		return ""
	}
	limit, offset := p.Page()
	if limit == 0 && offset == 0 {
		return ""
	}
	stmt := fmt.Sprintf("LIMIT %d", limit)
	if offset != 0 {
		stmt += fmt.Sprintf(" OFFSET %d", offset)
	}
	return stmt
}

func assignColumns(a []api.Assignment) string {
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

func assignSets(a []api.Assignment) string {
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

func qMarks(n int) string {
	if n == 0 {
		return ""
	}
	qMark := strings.Repeat("?, ", n)
	qMark = qMark[:len(qMark)-2] // remove last ", "
	return qMark
}
