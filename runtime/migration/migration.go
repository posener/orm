package migration

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/posener/orm/graph"
	"github.com/posener/orm/load"
)

// Table represents an SQL table for marshaling to go code
type Table struct {
	Columns     []Column
	PrimaryKeys []string
	ForeignKeys []ForeignKey
}

// Marshal returns a string of Table
func (t *Table) Marshal() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// UnMarshal takes a string and set the table content
func (t *Table) UnMarshal(s string) error {
	return json.Unmarshal([]byte(s), t)
}

// Column describe an SQL table
type Column struct {
	Name    string
	GoType  string
	SQLType string
	Options []string
}

// ForeignKey describes an SQL foreign key
type ForeignKey struct {
	Columns        []string
	Table          string
	ForeignColumns []string
}

// Table returns table structure to be used for generated code
func NewTable(gr *graph.Graph) *Table {
	t := new(Table)
	for _, f := range gr.Fields {
		if !f.IsReference() {
			sqlColumn := f.Columns()[0]
			t.Columns = append(t.Columns, Column{
				Name:    sqlColumn.Name,
				GoType:  f.Type.Naked.Ext(""),
				Options: options(f),
			})
		}
		if f.PrimaryKey {
			t.PrimaryKeys = append(t.PrimaryKeys, f.Column().Name)
		}
	}

	// define foreign keys for the outgoing references
	for _, e := range gr.Out {
		cols, fk := foreignKey(e)
		t.Columns = append(t.Columns, cols...)
		t.ForeignKeys = append(t.ForeignKeys, fk)
	}
	return t
}

func options(f *load.Field) []string {
	var stmt []string
	if f.NotNull {
		stmt = append(stmt, "NOT NULL")
	}
	if f.Null {
		stmt = append(stmt, "NULL")
	}
	if f.Default != "" {
		stmt = append(stmt, fmt.Sprintf("DEFAULT %s", f.Default))
	}
	if f.AutoIncrement {
		stmt = append(stmt, "AUTO_INCREMENT")
	}
	if f.Unique {
		stmt = append(stmt, "UNIQUE")
	}
	return stmt
}

func foreignKey(outEdge graph.Edge) (cols []Column, fk ForeignKey) {
	fk.Table = outEdge.RelationType().Table()
	dstFields := outEdge.RelationType().PrimaryKeys
	for i, col := range outEdge.SrcField.Columns() {
		cols = append(cols, Column{
			Name:   col.Name,
			GoType: dstFields[i].Type.Naked.Ext(""),
		})
		fk.Columns = append(fk.Columns, col.Name)
		fk.ForeignColumns = append(fk.ForeignColumns, dstFields[i].Column().Name)
	}
	return
}

// Hash returns a unique identifier for the foreign key
func (fk *ForeignKey) Hash() string {
	return strings.Join(fk.Columns, ",")
}
