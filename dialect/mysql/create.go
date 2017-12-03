package mysql

import (
	"fmt"
	"log"
	"strings"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/format"
	"github.com/posener/orm/load"
)

// ColumnsStatement returns the fields parts of SQL CREATE TABLE statement
func (g *Gen) ColumnsStatement(tp *load.Type) string {
	var (
		stmts []string
		refs  []common.ForeignKey
	)
	for _, f := range tp.Fields {
		switch {
		case f.Type.Slice:
		case f.IsReference():
			var (
				refType = f.Type
				pk      = refType.PrimaryKey
			)
			if pk == nil {
				log.Fatalf("Field %s (%s) is reference and does not have primary key defined", f.Name, refType.Name)
			}
			fk := common.ForeignKey{
				Column:    f.Column(),
				RefTable:  refType.Table(),
				RefColumn: pk.Column(),
			}
			stmts = append(stmts, fmt.Sprintf("`%s` %s", fk.Column, string(g.sqlType(refType.PrimaryKey))))
			refs = append(refs, fk)
		default:
			stmts = append(stmts, g.fieldCreateString(&f))
			if fk := f.SQL.ForeignKey; fk != nil {
				refs = append(refs, common.ForeignKey{
					Column:    f.Column(),
					RefTable:  fk.Type.Table(),
					RefColumn: fk.Field.Column(),
				})
			}
		}
	}
	stmts = append(stmts, format.ForeignKey(refs)...)
	return strings.Join(stmts, ", ")
}

func (g *Gen) fieldCreateString(f *load.Field) string {
	sqlType := g.sqlType(f)
	stmt := []string{fmt.Sprintf("`%s` %s", f.Column(), sqlType)}
	if f.SQL.NotNull {
		stmt = append(stmt, "NOT NULL")
	}
	if f.SQL.Null {
		stmt = append(stmt, "NULL")
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
