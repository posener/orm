package sqlite3

import (
	"fmt"
	"log"
	"strings"

	"github.com/posener/orm/def"
	"github.com/posener/orm/dialect/sqltypes"
)

func (s *sqlite3) Create() string {
	return fmt.Sprintf("CREATE TABLE '%s' ( %s )", s.tp.Table(), s.fieldsCreateString())
}
func (s *sqlite3) fieldsCreateString() string {
	fieldsStmt := make([]string, len(s.tp.Fields))
	for i, f := range s.tp.Fields {
		fieldsStmt[i] = fieldCreateString(&f)
	}
	return strings.Join(fieldsStmt, ", ")
}

func fieldCreateString(f *def.Field) string {
	stmt := []string{fmt.Sprintf("'%s' %s", f.SQL.Column, sqlType(f))}
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
		if !f.SQL.PrimaryKey || defaultSQLTypes(f.Type) != sqltypes.Integer {
			log.Fatalf("sqlite3 supports autoincrement only for 'INTEGER PRIMARY KEY' columns")
		}
		stmt = append(stmt, "AUTOINCREMENT")
	}
	if f.SQL.Unique {
		stmt = append(stmt, " UNIQUE")
	}
	return strings.Join(stmt, " ")
}
