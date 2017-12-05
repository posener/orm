package sqlite3

import (
	"fmt"
	"log"
	"strings"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/format"
	"github.com/posener/orm/dialect/sqltypes"
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
			stmts = append(stmts, g.fieldCreateString(f))
			if fk := f.ForeignKey; fk != nil {
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
	stmt := []string{fmt.Sprintf("'%s' %s", f.Column(), sqlType)}
	if f.NotNull {
		stmt = append(stmt, "NOT NULL")
	}
	if f.Null {
		stmt = append(stmt, "NULL")
	}
	if f.Default != "" {
		stmt = append(stmt, "DEFAULT", f.Default)
	}
	if f.PrimaryKey || f.AutoIncrement {
		stmt = append(stmt, "PRIMARY KEY")
	}
	if f.AutoIncrement {
		if !f.PrimaryKey || sqlType != sqltypes.Integer {
			log.Fatalf("Gen supports autoincrement only for 'INTEGER PRIMARY KEY' columns")
		}
		stmt = append(stmt, "AUTOINCREMENT")
	}
	if f.Unique {
		stmt = append(stmt, " UNIQUE")
	}
	return strings.Join(stmt, " ")
}
