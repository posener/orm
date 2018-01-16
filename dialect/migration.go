package dialect

import (
	"encoding/json"
	"fmt"
	"sort"
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
func (t Table) Marshal() string {
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

// Column describes an SQL table
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

func (f *ForeignKey) hash() string {
	cols := make([]string, len(f.Columns))
	copy(cols, f.Columns)
	sort.Strings(cols)
	return strings.Join(cols, "/")
}

// NewTable returns table structure to be used for generated code
func NewTable(gr *graph.Graph) Table {
	t := Table{}
	for _, f := range gr.Fields {
		if !f.IsReference() {
			sqlColumn := f.Columns()[0]
			t.Columns = append(t.Columns, Column{
				Name:    sqlColumn.Name,
				GoType:  f.Type.Naked.Ext(""),
				Options: options(f),
				SQLType: f.CustomType.String(),
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

func RelationTables(gr *graph.Graph) map[string]Table {
	var tables = make(map[string]Table)
	for _, e := range gr.RelTable {
		tables[e.Field.RelationTable()] = relationTable(e.Field)
	}
	return tables
}

func relationTable(f *load.Field) Table {
	var t1, t2 = f.RelationTypes()

	table := Table{}

	for _, t := range []*load.Naked{t1, t2} {
		var foreignColumns, columns []string
		for _, pk := range t.PrimaryKeys {
			columnName := fmt.Sprintf("%s_%s", t.Table(), pk.Column().Name)
			table.Columns = append(table.Columns, Column{
				Name:    columnName,
				GoType:  pk.Type.Naked.Ext(""),
				SQLType: pk.CustomType.String(),
			})
			columns = append(columns, columnName)
			foreignColumns = append(foreignColumns, pk.Column().Name)
		}

		table.ForeignKeys = append(table.ForeignKeys, ForeignKey{
			Table:          t.Table(),
			Columns:        columns,
			ForeignColumns: foreignColumns,
		})
	}
	return table
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
	fk.Table = outEdge.RelationType.Table()
	dstFields := outEdge.RelationType.PrimaryKeys
	for i, col := range outEdge.Field.Columns() {
		cols = append(cols, Column{
			Name:   col.Name,
			GoType: dstFields[i].Type.Naked.Ext(""),
		})
		fk.Columns = append(fk.Columns, col.Name)
		fk.ForeignColumns = append(fk.ForeignColumns, dstFields[i].Column().Name)
	}
	return
}

// Diff calculate difference between a 'got' table and a 'want' table
// and returns a table containing the columns, primary keys and foreign keys that
// we need to add the that 'got' table.
// It returns an error in case that there is a change that can't be applied.
func (t *Table) Diff(want *Table) (*Table, error) {
	var (
		got    = t
		gotMap = make(map[string]bool)
		diff   = new(Table)
	)
	// see what we got first
	for _, c := range got.Columns {
		gotMap["col/"+c.Name] = true
	}
	for _, fk := range got.ForeignKeys {
		gotMap["fk/"+fk.hash()] = true
	}
	for _, pk := range got.PrimaryKeys {
		gotMap["pk/"+pk] = true
	}

	// append to 'diff' table all missing things
	for _, col := range want.Columns {
		if !gotMap["col/"+col.Name] {
			diff.Columns = append(diff.Columns, col)
		}
	}
	for _, fk := range want.ForeignKeys {
		if !gotMap["fk/"+fk.hash()] {
			diff.ForeignKeys = append(diff.ForeignKeys, fk)
		}
	}

	for _, pk := range want.PrimaryKeys {
		if !gotMap["pk/"+pk] {
			diff.PrimaryKeys = append(diff.PrimaryKeys, pk)
		}
	}
	return diff, nil
}
