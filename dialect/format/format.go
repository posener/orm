package format

import (
	"bytes"
	"fmt"

	"github.com/posener/orm/common"
)

// Columns extract SQL columns list statement
func Columns(c common.Columner) string {
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

// Where formats an SQL WHERE statement
func Where(c common.StatementArger) string {
	if c == nil {
		return ""
	}
	where := c.Statement()
	if len(where) == 0 {
		return ""
	}
	return "WHERE " + where
}

// GroupBy formats an SQL GROUP BY statement
func GroupBy(groups []common.Group) string {
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

// OrderBy formats an SQL ORDER BY statement
func OrderBy(orders []common.Order) string {
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

// Page formats an SQL LIMIT...OFFSET statement
func Page(p common.Page) string {
	if p.Limit == 0 { // why would someone ask for a page of zero size?
		return ""
	}
	stmt := fmt.Sprintf("LIMIT %d", p.Limit)
	if p.Offset != 0 {
		stmt += fmt.Sprintf(" OFFSET %d", p.Offset)
	}
	return stmt
}

// AssignSets formats a list of assignments for SQL UPDATE SET statements
func AssignSets(a common.Assignments) string {
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

// AssignColumns gets an assignment list and formats the assign column names
// for an SQL INSERT STATEMENT
func AssignColumns(a common.Assignments) string {
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

// IfNotExists formats an SQL IF NOT EXISTS statement
func IfNotExists(ifNotExists bool) interface{} {
	if ifNotExists {
		return "IF NOT EXISTS"
	}
	return ""
}
