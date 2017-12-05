package dialect

import (
	"fmt"
	"log"
	"strings"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/format"
	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/sqlite3"
	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/load"
)

// Generator is API for different dialects
type Generator interface {
	// Name is the dialect name
	Name() string
	// ColumnsStatement returns the fields parts of SQL CREATE TABLE statement
	// for a specific struct and specific dialect.
	// It is used by the generation tool.
	ColumnsStatement(tp *load.Type) string
	// ConvertValueCode returns go code for converting value returned from the
	// database to the given field.
	ConvertValueCode(field *load.Field) string
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
	GoTypeToColumnType(*load.Type) sqltypes.Type
	ColumnCreateString(*load.Field, sqltypes.Type) string
	ConvertValueCode(*load.Field, sqltypes.Type) string
}

// ColumnsStatement returns the fields parts of SQL CREATE TABLE statement
func (g *gen) ColumnsStatement(tp *load.Type) string {
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
			stmts = append(stmts, fmt.Sprintf("`%s` %s", fk.Column, string(g.columnType(refType.PrimaryKey))))
			refs = append(refs, fk)
		default:
			stmts = append(stmts, g.ColumnCreateString(f, g.columnType(f)))
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

func (g *gen) ConvertValueCode(field *load.Field) string {
	return g.GenImplementer.ConvertValueCode(field, g.columnType(field))
}

func (g *gen) columnType(field *load.Field) sqltypes.Type {
	if custom := field.CustomType; custom != "" {
		return custom
	}
	return g.GoTypeToColumnType(field.SetType())
}
