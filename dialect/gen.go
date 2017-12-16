package dialect

import (
	"fmt"
	"log"
	"strings"

	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/sqlite3"
	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/graph"
	"github.com/posener/orm/load"
)

// Generator is API for different dialects
type Generator interface {
	// Name is the dialect name
	Name() string
	// ConvertValueCode returns go code for converting value returned from the
	// database to the given field.
	ConvertValueCode(tp *load.Type, field *load.Field) string
}

// NewGen returns all known Generators
func NewGen() []Generator {
	return []Generator{
		&gen{GenImplementer: new(mysql.Gen)},
		&gen{GenImplementer: new(sqlite3.Gen)},
	}
}

type gen struct {
	GenImplementer
}

type GenImplementer interface {
	Name() string
	GoTypeToColumnType(*load.Type) *sqltypes.Type
	Translate(string) string
	PreProcess(f *load.Field, sqlType *sqltypes.Type) error
	ConvertValueCode(*load.Type, *load.Field, *sqltypes.Type) string
}

// ColumnsStatement returns the fields parts of SQL CREATE TABLE statement
func (g *gen) ColumnsStatement(gr *graph.Graph) string {
	var (
		colStmts []string
		fkStmts  []string
	)
	for _, f := range gr.Fields {
		if !f.IsReference() {
			sqlColumn := f.Columns()[0]
			sqlType := g.columnType(&sqlColumn)
			err := g.PreProcess(f, sqlType)
			if err != nil {
				log.Fatal(err)
			}
			colStmts = append(colStmts, g.columnCreateString(sqlColumn.Name, f, sqlType))
		}
	}

	// define foreign keys for the outgoing references
	for _, e := range gr.Out {
		eColStmts, eFKStmts := g.foreignKeys(e)
		colStmts = append(colStmts, eColStmts...)
		fkStmts = append(fkStmts, eFKStmts...)
	}
	stmts := append(colStmts, fkStmts...)
	return strings.Join(stmts, ", ")
}

func (g *gen) columnCreateString(name string, f *load.Field, sqlType *sqltypes.Type) string {
	stmt := []string{fmt.Sprintf("`%s` %s", name, sqlType)}
	if f.NotNull {
		stmt = append(stmt, g.Translate("NOT NULL"))
	}
	if f.Null {
		stmt = append(stmt, g.Translate("NULL"))
	}
	if f.Default != "" {
		stmt = append(stmt, g.Translate("DEFAULT"), f.Default)
	}
	if f.PrimaryKey {
		stmt = append(stmt, g.Translate("PRIMARY KEY"))
	}
	if f.AutoIncrement {
		stmt = append(stmt, g.Translate("AUTO_INCREMENT"))
	}
	if f.Unique {
		stmt = append(stmt, g.Translate("UNIQUE"))
	}
	return strings.Join(stmt, " ")
}

func (g *gen) ConvertValueCode(tp *load.Type, field *load.Field) string {
	if field.IsReference() {
		return ""
	}
	return g.GenImplementer.ConvertValueCode(tp, field, g.columnType(&field.Columns()[0]))
}

func (g *gen) columnType(col *load.SQLColumn) *sqltypes.Type {
	if custom := col.CustomType; custom != nil {
		return custom
	}
	return g.GoTypeToColumnType(col.SetType)
}

func (g *gen) foreignKeys(outEdge graph.Edge) (colStmts []string, fkStmts []string) {
	cols := outEdge.SrcField.Columns()
	dstFields := outEdge.RelationType().PrimaryKeys
	for i := range cols {
		colStmts = append(colStmts,
			fmt.Sprintf("`%s` %s", cols[i].Name, g.GoTypeToColumnType(&dstFields[i].Type)))
		fkStmts = append(fkStmts,
			fmt.Sprintf("FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`)",
				cols[i].Name, outEdge.RelationType().Table(), dstFields[i].Column().Name))
	}
	return
}
