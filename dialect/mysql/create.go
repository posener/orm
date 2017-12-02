package mysql

import (
	"fmt"
	"log"
	"strings"

	"github.com/posener/orm/load"
)

// ColumnsStatement returns the fields parts of SQL CREATE TABLE statement
func (g *Gen) ColumnsStatement(tp *load.Type) string {
	var (
		stmts []string
		refs  []*load.Field
	)
	for _, f := range tp.Fields {
		if f.IsReference() {
			refs = append(refs, &f)
		} else {
			stmts = append(stmts, g.fieldCreateString(&f))
		}
	}
	for _, f := range refs {
		stmts = append(stmts, g.referenceString(f)...)
	}
	return strings.Join(stmts, ", ")
}

func (g *Gen) fieldCreateString(f *load.Field) string {
	sqlType := g.sqlType(f)
	stmt := []string{fmt.Sprintf("`%s` %s", f.Column(), sqlType)}
	if f.SQL.NotNull {
		stmt = append(stmt, "NOT NULL")
	}
	if f.SQL.Default != "" {
		stmt = append(stmt, "DEFAULT", f.SQL.Default)
	}
	if f.SQL.PrimaryKey || f.SQL.AutoIncrement {
		stmt = append(stmt, "PRIMARY KEY")
	}
	if f.SQL.AutoIncrement {
		stmt = append(stmt, "AUTO_INCREMENT")
	}
	if f.SQL.Unique {
		stmt = append(stmt, " UNIQUE")
	}
	return strings.Join(stmt, " ")
}

func (g *Gen) referenceString(f *load.Field) []string {
	refType := f.Type
	pk := refType.PrimaryKey
	if pk == nil {
		log.Fatalf("reference %s does not have primary key defined", refType)
	}
	var (
		columnName = f.Column()
		columnType = g.sqlType(refType.PrimaryKey)
		refTable   = refType.Table()
		refColumn  = pk.Column()
	)
	return []string{
		fmt.Sprintf("%s %s", columnName, columnType),
		fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", columnName, refTable, refColumn),
	}
}
