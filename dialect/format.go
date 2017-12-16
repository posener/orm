package dialect

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/posener/orm/common"
)

// columns extract SQL wantCols list statement
func columns(p *common.SelectParams) string {
	parts := columnsParts(p.Table, p)

	// we can add COUNT(*) only once and it work only on the most upper level
	if p.Columns.Count() {
		parts = append(parts, "COUNT(*)")
	}

	return strings.Join(parts, ", ")
}

func columnsParts(table string, p *common.SelectParams) []string {
	var (
		parts  []string
		exists = make(map[string]bool)
	)

	parts = append(parts, columnsCollect(table, p.Columns.Columns())...)

	for _, join := range p.Columns.Joins() {
		for _, part := range columnsParts(join.TableName(table), &join.SelectParams) {
			if exists[part] {
				continue
			}
			parts = append(parts, part)
			exists[part] = true
		}
	}
	return parts
}

func columnsCollect(table string, cols []string) []string {
	var parts []string
	for _, col := range cols {
		parts = append(parts, fmt.Sprintf("`%s`.`%s`", table, col))
	}
	return parts
}

// whereJoin takes SelectParams and traverse all the join options
// it concat all the conditions with an AND operator
func whereJoin(table string, p *common.SelectParams) string {
	stmt := whereJoinRec(table, p)
	if stmt == "" {
		return ""
	}
	return "WHERE " + stmt
}

func whereJoinRec(table string, p *common.SelectParams) string {
	var parts []string
	if p.Where != nil {
		if w := p.Where.Statement(table); w != "" {
			parts = append(parts, w)
		}
	}
	for _, join := range p.Columns.Joins() {
		joinCond := whereJoinRec(join.TableName(table), &join.SelectParams)
		if joinCond != "" {
			parts = append(parts, joinCond)
		}
	}
	return strings.Join(parts, " AND ")
}

func where(table string, c common.StatementArger) string {
	if c == nil {
		return ""
	}
	where := c.Statement(table)
	if len(where) == 0 {
		return ""
	}
	return "WHERE " + where
}

// groupBy formats an SQL GROUP BY statement
func groupBy(table string, groups []common.Group) string {
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

// orderBy formats an SQL ORDER BY statement
func orderBy(table string, orders []common.Order) string {
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

// page formats an SQL LIMIT...OFFSET statement
func page(p common.Page) string {
	if p.Limit == 0 { // why would someone ask for a page of zero size?
		return ""
	}
	stmt := fmt.Sprintf("LIMIT %d", p.Limit)
	if p.Offset != 0 {
		stmt += fmt.Sprintf(" OFFSET %d", p.Offset)
	}
	return stmt
}

// assignSets formats a list of assignments for SQL UPDATE SET statements
func assignSets(a common.Assignments) string {
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

// assignColumns gets an assignment list and formats the assign column names
// for an SQL INSERT STATEMENT
func assignColumns(a common.Assignments) string {
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

// IfNotExists formats an SQL IF NOT EXISTS statement
func IfNotExists(ifNotExists bool) string {
	if ifNotExists {
		return "IF NOT EXISTS"
	}
	return ""
}
