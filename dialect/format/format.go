package format

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/posener/orm/common"
)

// Columns extract SQL columns list statement
func Columns(table string, c common.Selector) string {
	var parts []string

	if c == nil {
		return fmt.Sprintf("`%s`.*", table)
	}
	cols := c.Columns()
	parts = append(parts, formatColumns(table, cols, c.Count())...)
	if joins := c.Joins(); joins != nil {
		for _, join := range joins {
			parts = append(parts, formatColumns(join.RefTable, join.SelectColumns, c.Count())...)
		}
	}
	if c.Count() {
		parts = append(parts, "COUNT(*)")
	}
	return strings.Join(parts, ", ")
}

func formatColumns(table string, cols []string, isCount bool) []string {
	if len(cols) == 0 && !isCount {
		return []string{fmt.Sprintf("`%s`.*", table)}
	}
	var parts []string
	for _, col := range cols {
		parts = append(parts, fmt.Sprintf("`%s`.`%s`", table, col))
	}
	return parts
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
func Join(table string, c common.Selector) string {
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
		tables = append(tables, fmt.Sprintf("`%s`", j.RefTable))
		conds = append(conds, fmt.Sprintf("`%s`.`%s` = `%s`.`%s`", table, j.Column, j.RefTable, j.RefColumn))
	}
	return fmt.Sprintf("JOIN (%s) ON (%s)",
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)
}

// IfNotExists formats an SQL IF NOT EXISTS statement
func IfNotExists(ifNotExists bool) string {
	if ifNotExists {
		return "IF NOT EXISTS"
	}
	return ""
}

func ForeignKey(foreignKeys []common.ForeignKey) []string {
	var stmts []string
	for _, fk := range foreignKeys {
		stmts = append(stmts, fmt.Sprintf("FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`)",
			fk.Column, fk.RefTable, fk.RefColumn))
	}
	return stmts
}
