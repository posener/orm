package dialect

import (
	"fmt"
	"strings"

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
	GoTypeToColumnType(*load.Type) sqltypes.Type
	ColumnCreateString(*load.Field, sqltypes.Type) string
	ConvertValueCode(*load.Type, *load.Field, sqltypes.Type) string
}

// ColumnsStatement returns the fields parts of SQL CREATE TABLE statement
func (g *gen) ColumnsStatement(tp *load.Type) string {
	var (
		stmts []string
		refs  []*load.ForeignKey
	)
	for _, f := range tp.Fields {
		if !f.Type.Slice {
			stmts = append(stmts, g.ColumnCreateString(f, g.columnType(f)))
			if fk := f.ForeignKey; fk != nil {
				refs = append(refs, fk)
			}
		}
	}
	stmts = append(stmts, foreignKeys(refs)...)
	return strings.Join(stmts, ", ")
}

func (g *gen) ConvertValueCode(tp *load.Type, field *load.Field) string {
	return g.GenImplementer.ConvertValueCode(tp, field, g.columnType(field))
}

func (g *gen) columnType(field *load.Field) sqltypes.Type {
	if custom := field.CustomType; custom != "" {
		return custom
	}
	return g.GoTypeToColumnType(field.SetType())
}

func foreignKeys(foreignKeys []*load.ForeignKey) []string {
	var stmts []string
	for _, fk := range foreignKeys {
		stmts = append(stmts, fmt.Sprintf("FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`)",
			fk.Src.Column(), fk.Dst.ParentType.Table(), fk.Dst.ParentType.PrimaryKey.Column()))
	}
	return stmts
}
