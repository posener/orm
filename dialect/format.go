package dialect

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/posener/orm/runtime"
)

// tableProperties returns all properties of SQL table, as should be given in the table CREATE statement
func (d *dialect) tableProperties(t *runtime.Table) string {
	var stmts []string
	for _, col := range t.Columns {
		stmts = append(stmts, d.createColumn(col))
	}
	if len(t.PrimaryKeys) > 0 {
		stmts = append(stmts, fmt.Sprintf("PRIMARY KEY (%s)",
			strings.Join(d.quoteSlice(t.PrimaryKeys), ", ")))
	}
	for _, fk := range t.ForeignKeys {
		stmts = append(stmts, d.foreignKey(fk))
	}
	return strings.Join(stmts, ", ")
}

// createColumn is an SQL column definition, as given in the SQL CREATE statement
func (d *dialect) createColumn(col runtime.Column) string {
	s := fmt.Sprintf("%s %s", d.Quote(col.Name), d.GoTypeToColumnType(col.GoType))
	for _, opt := range col.Options {
		s += " " + d.Translate(opt)
	}
	return s
}

// foreignKey is teh FOREIGN KEY statement
func (d *dialect) foreignKey(fk runtime.ForeignKey) string {
	return fmt.Sprintf(
		"FOREIGN KEY (%s) REFERENCES %s(%s)",
		strings.Join(d.quoteSlice(fk.Columns), ", "),
		fk.Table,
		strings.Join(d.quoteSlice(fk.ForeignColumns), ", "),
	)
}

// selectColumns returns the columns selected for an SQL SELECT query
func (d *dialect) selectColumns(p *runtime.SelectParams) string {
	parts := d.columnsParts(p.Table, p)

	// we can add COUNT(*) only once and it work only on the most upper level
	if p.Columns.Count() {
		parts = append(parts, "COUNT(*)")
	}

	return strings.Join(parts, ", ")
}

func (d *dialect) columnsParts(table string, p *runtime.SelectParams) []string {
	var (
		parts  []string
		exists = make(map[string]bool)
	)

	parts = append(parts, d.columnsCollect(table, p.Columns.Columns())...)

	for _, join := range p.Columns.Joins() {
		for _, part := range d.columnsParts(join.TableName(table), &join.SelectParams) {
			if exists[part] {
				continue
			}
			parts = append(parts, part)
			exists[part] = true
		}
	}
	return parts
}

func (d *dialect) columnsCollect(table string, cols []string) []string {
	var parts []string
	for _, col := range cols {
		parts = append(parts, fmt.Sprintf("%s.%s", d.Quote(table), d.Quote(col)))
	}
	return parts
}

// whereJoin takes SelectParams and traverse all the join options
// it concat all the conditions with an AND operator
func (d *dialect) whereJoin(table string, p *runtime.SelectParams) string {
	stmt := d.whereJoinRec(table, p)
	if stmt == "" {
		return ""
	}
	return "WHERE " + stmt
}

// whereJoinRec returns a WHERE statement for a recursive join statement
func (d *dialect) whereJoinRec(table string, p *runtime.SelectParams) string {
	var parts []string
	if p.Where != nil {
		if w := p.Where.Statement(table); w != "" {
			parts = append(parts, w)
		}
	}
	for _, join := range p.Columns.Joins() {
		joinCond := d.whereJoinRec(join.TableName(table), &join.SelectParams)
		if joinCond != "" {
			parts = append(parts, joinCond)
		}
	}
	return strings.Join(parts, " AND ")
}

// where returns an SQL WHERE statement
func (d *dialect) where(table string, c runtime.StatementArger) string {
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
func (d *dialect) groupBy(table string, groups []runtime.Group) string {
	if len(groups) == 0 {
		return ""
	}
	b := bytes.NewBufferString("GROUP BY ")
	for i := range groups {
		b.WriteString(fmt.Sprintf("%s.%s, ", d.Quote(table), d.Quote(groups[i].Column)))
	}

	s := b.String()
	return s[:len(s)-2]
}

// orderBy formats an SQL ORDER BY statement
func (d *dialect) orderBy(table string, orders []runtime.Order) string {
	if len(orders) == 0 {
		return ""
	}

	b := bytes.NewBufferString("ORDER BY ")
	for i := range orders {
		b.WriteString(fmt.Sprintf("%s.%s %s, ", d.Quote(table), d.Quote(orders[i].Column), orders[i].Dir))
	}

	s := b.String()
	return s[:len(s)-2]
}

// page formats an SQL LIMIT...OFFSET statement
func (d *dialect) page(p runtime.Page) string {
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
func (d *dialect) assignSets(a runtime.Assignments) string {
	if len(a) == 0 {
		return ""
	}
	b := bytes.NewBuffer(nil)
	for i := range a {
		b.WriteString(fmt.Sprintf("%s = ?, ", d.Quote(a[i].Column)))
	}

	s := b.String()
	return s[:len(s)-2]
}

// assignColumns gets an assignment list and formats the assign createColumn names
// for an SQL INSERT STATEMENT
func (d *dialect) assignColumns(a runtime.Assignments) string {
	if len(a) == 0 {
		return ""
	}
	b := bytes.NewBuffer(nil)
	for i := range a {
		b.WriteString(fmt.Sprintf("%s, ", d.Quote(a[i].Column)))
	}

	s := b.String()
	return s[:len(s)-2]
}

// ifNotExists formats an SQL IF NOT EXISTS statement
func (d *dialect) ifNotExists(ifNotExists bool) string {
	if ifNotExists {
		return "IF NOT EXISTS"
	}
	return ""
}

func (d *dialect) quoteSlice(s []string) []string {
	for i := range s {
		s[i] = d.Quote(s[i])
	}
	return s
}
