package dialect

import (
	"fmt"
	"strings"

	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/sqlite3"
	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/graph"
	"github.com/posener/orm/load"
)

// Generator generates static content of model for table creation and row decoding
type Generator interface {
	// Name is the dialect name
	Name() string
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
	GoTypeToColumnType(string) *sqltypes.Type
	Translate(string) string
	ConvertValueCode(*load.Field, *sqltypes.Type) string
}

type Table struct {
	Columns     []Column
	PrimaryKeys []string
	ForeignKeys []ForeignKey
}

type Column struct {
	Name    string
	Type    string
	Options string
}

func (c *Column) String() string {
	s := fmt.Sprintf("`%s` %s", c.Name, c.Type)
	if c.Options != "" {
		s += " " + c.Options
	}
	return s
}

type ForeignKey struct {
	Columns        []string
	Table          string
	ForeignColumns []string
}

func (fk *ForeignKey) String() string {
	return fmt.Sprintf(
		"FOREIGN KEY (%s) REFERENCES %s(%s)",
		strings.Join(quote(fk.Columns), ", "),
		fk.Table,
		strings.Join(quote(fk.ForeignColumns), ", "),
	)
}

func quote(s []string) []string {
	for i := range s {
		s[i] = fmt.Sprintf("`%s`", s[i])
	}
	return s
}

func (fk *ForeignKey) Hash() string {
	return strings.Join(fk.Columns, ",")
}

// Table returns table structure for generated code
func (g *gen) Table(gr *graph.Graph) *Table {
	t := new(Table)
	for _, f := range gr.Fields {
		if !f.IsReference() {
			sqlColumn := f.Columns()[0]
			sqlType := g.columnType(&sqlColumn)
			t.Columns = append(t.Columns, Column{
				Name:    sqlColumn.Name,
				Type:    sqlType.String(),
				Options: g.constructColumnStmt(f),
			})
		}
		if f.PrimaryKey {
			t.PrimaryKeys = append(t.PrimaryKeys, f.Column().Name)
		}
	}

	// define foreign keys for the outgoing references
	for _, e := range gr.Out {
		cols, fk := g.foreignKey(e)
		t.Columns = append(t.Columns, cols...)
		t.ForeignKeys = append(t.ForeignKeys, fk)
	}
	return t
}

func (g *gen) constructColumnStmt(f *load.Field) string {
	var stmt []string
	if f.NotNull {
		stmt = append(stmt, g.Translate("NOT NULL"))
	}
	if f.Null {
		stmt = append(stmt, g.Translate("NULL"))
	}
	if f.Default != "" {
		stmt = append(stmt, g.Translate("DEFAULT"), f.Default)
	}
	if f.AutoIncrement {
		stmt = append(stmt, g.Translate("AUTO_INCREMENT"))
	}
	if f.Unique {
		stmt = append(stmt, g.Translate("UNIQUE"))
	}
	return strings.Join(stmt, " ")
}

func (g *gen) ConvertValueCode(field *load.Field) string {
	if field.IsReference() {
		return ""
	}
	return g.GenImplementer.ConvertValueCode(field, g.columnType(&field.Columns()[0]))
}

func (g *gen) columnType(col *load.SQLColumn) *sqltypes.Type {
	if custom := col.CustomType; custom != nil {
		return custom
	}
	return g.GoTypeToColumnType(col.SetType.Naked.Ext(""))
}

func (g *gen) foreignKey(outEdge graph.Edge) (cols []Column, fk ForeignKey) {
	fk.Table = outEdge.RelationType().Table()
	dstFields := outEdge.RelationType().PrimaryKeys
	for i, col := range outEdge.SrcField.Columns() {
		cols = append(cols, Column{
			Name: col.Name,
			Type: g.GoTypeToColumnType(dstFields[i].Type.Naked.Ext("")).String(),
		})
		fk.Columns = append(fk.Columns, col.Name)
		fk.ForeignColumns = append(fk.ForeignColumns, dstFields[i].Column().Name)
	}
	return
}
