package dialect

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/posener/orm/common"
)

// columns extract SQL wantCols list statement
func columns(p *common.SelectParams) string {
	parts := columnsParts(p)

	// we can add COUNT(*) only once and it work only on the most upper level
	if p.Columns.Count() {
		parts = append(parts, "COUNT(*)")
	}

	return strings.Join(parts, ", ")
}

func columnsParts(p *common.SelectParams) []string {
	var (
		parts  []string
		exists = make(map[string]bool)
	)

	parts = append(parts, columnsCollect(p.Table, p.Columns.Columns(), p.Columns.Count())...)

	for _, join := range p.Columns.Joins() {
		for _, part := range columnsParts(&join.SelectParams) {
			if exists[part] {
				continue
			}
			parts = append(parts, part)
			exists[part] = true
		}
	}
	return parts
}

func columnsCollect(table string, cols []string, isCount bool) []string {
	if len(cols) == 0 && !isCount {
		return []string{fmt.Sprintf("`%s`.*", table)}
	}
	var parts []string
	for _, col := range cols {
		parts = append(parts, fmt.Sprintf("`%s`.`%s`", table, col))
	}
	return parts
}

// whereJoin takes SelectParams and traverse all the join options
// it concat all the conditions with an AND operator
func whereJoin(p *common.SelectParams) common.Where {
	w := p.Where
	for _, join := range p.Columns.Joins() {
		if w != nil {
			w = w.And(whereJoin(&join.SelectParams))
		} else {
			w = join.SelectParams.Where
		}
	}
	return w
}

func where(c common.StatementArger) string {
	if c == nil {
		return ""
	}
	where := c.Statement()
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

// join extract SQL join list statement
func join(p *common.SelectParams) string {
	return strings.Join(joinParts(p), " ")
}

func joinParts(p *common.SelectParams) []string {
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
		recursive = append(recursive, joinParts(&j.SelectParams)...)
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
