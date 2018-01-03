package migration

import (
	"context"
	"fmt"

	"github.com/posener/orm"
)

// Load loads a table from an existing database
func Load(ctx context.Context, conn orm.Conn, tableName string) (*Table, error) {
	columns, err := columns(ctx, conn, tableName)
	if err != nil {
		return nil, fmt.Errorf("get table %s columns information: %s", tableName, err)
	}
	indices, err := indices(ctx, conn, tableName)
	if err != nil {
		return nil, fmt.Errorf("get table %s indicies information: %s", tableName, err)
	}
	t := new(Table)
	for _, col := range columns {
		newCol := Column{
			Name:    col.Field,
			SQLType: col.Type,
		}
		if col.Key != nil {
		}
		if col.Null != nil {
			switch *col.Null {
			case "YES":
				newCol.Options = append(newCol.Options, "NULL")
			case "NO":
				newCol.Options = append(newCol.Options, "NOT NULL")
			}
		}
		if col.Extra != nil {
			switch *col.Extra {
			case "auto_increment":
				newCol.Options = append(newCol.Options, "AUTO_INCREMENT")
			}
		}
		t.Columns = append(t.Columns, newCol)
	}

	// collect foreign keys and set primary keys
	fks := map[string]ForeignKey{}
	for _, index := range indices {
		if index.KeyName == "PRIMARY" {
			t.PrimaryKeys = append(t.PrimaryKeys, index.ColumnName)
		} else {
			fk := fks[index.KeyName]
			fk.Columns = append(fks[index.KeyName].Columns, index.ColumnName)
			fks[index.KeyName] = fk
		}
	}
	for _, fk := range fks {
		t.ForeignKeys = append(t.ForeignKeys, fk)
	}

	return t, nil
}

// columns returns all columns of a table by doing an SQL query
func columns(ctx context.Context, conn orm.Conn, tableName string) ([]column, error) {
	rows, err := conn.QueryContext(ctx, fmt.Sprintf("DESCRIBE `%s`", tableName))
	if err != nil {
		return nil, err
	}
	var cols []column
	for rows.Next() {
		var col column
		err = rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, nil
}

func indices(ctx context.Context, conn orm.Conn, tableName string) ([]index, error) {
	rows, err := conn.QueryContext(ctx, fmt.Sprintf("SHOW INDEX FROM `%s`", tableName))
	if err != nil {
		return nil, err
	}
	var indices []index
	for rows.Next() {
		var i index
		err = rows.Scan(
			&i.Table, &i.NonUnique, &i.KeyName, &i.SeqInIndex, &i.ColumnName, &i.Collation, &i.Cardinality,
			&i.SubPart, &i.Packed, &i.Null, &i.IndexType, &i.Comment, &i.IndexComment,
		)
		if err != nil {
			return nil, err
		}
		indices = append(indices, i)
	}
	return indices, nil
}

// column describes an SQL column from an SQL DESCRIBE command
type column struct {
	Field   string
	Type    string
	Null    *string
	Key     *string
	Default *string
	Extra   *string
}

// index describes an SQL index from an SQL SHOW INDEX FROM command
type index struct {
	Table        string
	NonUnique    int
	KeyName      string
	SeqInIndex   int
	ColumnName   string
	Collation    string
	Cardinality  int
	SubPart      *string
	Packed       *string
	Null         string
	IndexType    string
	Comment      *string
	IndexComment *string
}
