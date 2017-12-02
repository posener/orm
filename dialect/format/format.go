package format

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/posener/orm/common"
)

// Columns extract SQL columns list statement
func Columns(table string, c common.Columner) string {
	if c == nil {
		return "*"
	}
	cols := c.Columns()
	if len(cols) == 0 && !c.Count() {
		return "*"
	}
	b := bytes.NewBuffer(nil)
	for i := range cols {
		b.WriteString(fmt.Sprintf("`%s`.`%s`, ", table, cols[i]))
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
func GroupBy(table string, groups []common.Group) string {
	if len(groups) == 0 {
		return ""
	}
	b := bytes.NewBufferString("GROUP BY ")
	for i := range groups {
		b.WriteString(fmt.Sprintf("`%s`.`%s`, ", table, groups[i].Column))
	}

	s := b.String()
	return s[:len(s)-2]
}

// OrderBy formats an SQL ORDER BY statement
func OrderBy(table string, orders []common.Order) string {
	if len(orders) == 0 {
		return ""
	}

	b := bytes.NewBufferString("ORDER BY ")
	for i := range orders {
		b.WriteString(fmt.Sprintf("`%s`.`%s` %s, ", table, orders[i].Column, orders[i].Dir))
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
		b.WriteString(fmt.Sprintf("`%s` = ?, ", a[i].Column))
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
		b.WriteString(fmt.Sprintf("`%s`, ", a[i].Column))
	}

	s := b.String()
	return s[:len(s)-2]
}

// Join extract SQL join list statement
func Join(table string, c common.Columner) string {
	if c == nil {
		return ""
	}
	joins := c.Joins()
	if len(joins) == 0 {
		return ""
	}
	var (
		tables []string
		conds  []string
	)
	for _, j := range joins {
		tables = append(tables, fmt.Sprintf("'%s'", j.RefTable))
		conds = append(conds, fmt.Sprintf("`%s`.`%s` = `%s`.`%s`", table, j.Column, j.RefTable, j.RefColumn))
	}
	return fmt.Sprintf("JOIN (%s) ON (%s)",
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)
}

// IfNotExists formats an SQL IF NOT EXISTS statement
func IfNotExists(ifNotExists bool) interface{} {
	if ifNotExists {
		return "IF NOT EXISTS"
	}
	return ""
}
