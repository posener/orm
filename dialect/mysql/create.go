package mysql

import (
	"fmt"
	"strings"

	"github.com/posener/orm/common"
)

func (g *Gen) ColumnsStatement() string {
	fieldsStmt := make([]string, len(g.Tp.Fields))
	for i, f := range g.Tp.Fields {
		fieldsStmt[i] = g.fieldCreateString(&f)
	}
	return strings.Join(fieldsStmt, ", ")
}

func (g *Gen) fieldCreateString(f *common.Field) string {
	sqlType := g.sqlType(f)
	stmt := []string{fmt.Sprintf("`%s` %s", f.SQL.Column, sqlType)}
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
