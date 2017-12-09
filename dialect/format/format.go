package format

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/posener/orm/common"
)

// Columns extract SQL wantCols list statement
func Columns(p *common.SelectParams) string {
	parts := columns(p)

	// we can add COUNT(*) only once and it work only on the most upper level
	if p.Columns.Count() {
		parts = append(parts, "COUNT(*)")
	}

	return strings.Join(parts, ", ")
}

func columns(p *common.SelectParams) []string {
	var (
		parts  []string
		exists = make(map[string]bool)
	)

	parts = append(parts, collectColumns(p.Table, p.Columns.Columns(), p.Columns.Count())...)

	for _, join := range p.Columns.Joins() {
		for _, part := range columns(&join.SelectParams) {
			if exists[part] {
				continue
			}
			parts = append(parts, part)
			exists[part] = true
		}
	}
	return parts
}

func collectColumns(table string, cols []string, isCount bool) []string {
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
func Join(p *common.SelectParams) string {
	return strings.Join(join(p), " ")
}

func join(p *common.SelectParams) []string {
	joins := p.Columns.Joins()
	if len(joins) == 0 {
		return nil
	}
	var (
		tables    []string
		conds     []string
		recursive []string
	)
	for _, j := range joins {
		tables = append(tables, fmt.Sprintf("`%s`", j.Table))
		for _, pairing := range j.Pairings {
			conds = append(conds, fmt.Sprintf("`%s`.`%s` = `%s`.`%s`", p.Table, pairing.Column, j.Table, pairing.JoinedColumn))
		}
		recursive = append(recursive, Join(&j.SelectParams))
	}

	joinStmt := fmt.Sprintf("JOIN (%s) ON (%s)",
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return append([]string{joinStmt}, recursive...)
}

// IfNotExists formats an SQL IF NOT EXISTS statement
func IfNotExists(ifNotExists bool) string {
	if ifNotExists {
		return "IF NOT EXISTS"
	}
	return ""
}
