package dialect

import (
	"context"
	"fmt"

	"github.com/posener/orm"
)

// Load loads a table from an existing database
func (d *dialect) loadTable(ctx context.Context, conn orm.Conn, tableName string) (*Table, error) {
	exists, err := d.tableExists(ctx, conn, tableName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	columns, err := d.columns(ctx, conn, tableName)
	if err != nil {
		return nil, fmt.Errorf("get table %s columns information: %s", tableName, err)
	}
	indices, err := d.indices(ctx, conn, tableName)
	if err != nil {
		return nil, fmt.Errorf("get table %s indicies information: %s", tableName, err)
	}
	t := new(Table)
	for _, col := range columns {
		newCol := Column{
			Name:    col.Name,
			SQLType: col.Type,
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
		if index.Primary {
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

func (d *dialect) tableExists(ctx context.Context, conn orm.Conn, tableName string) (bool, error) {
	b := newBuilder(d, "SELECT COUNT(*) FROM information_schema.tables WHERE table_name =")
	b.Var(tableName)
	var count int
	err := conn.QueryRow(ctx, b.Statement(), b.Args()...).Scan(&count)
	return count > 0, err
}

// columns returns all columns of a table by doing an SQL query
func (d *dialect) columns(ctx context.Context, conn orm.Conn, tableName string) ([]column, error) {
	var stmt string
	switch d.Name() {
	case Postgres:
		stmt = "SELECT column_name, CONCAT(udt_name, '(', character_maximum_length, datetime_precision, ')'), is_nullable, column_default, '' FROM information_schema.columns WHERE table_name ="
	default:
		stmt = "SELECT column_name, column_type, is_nullable, column_default, extra FROM information_schema.columns WHERE table_name ="
	}

	b := newBuilder(d, stmt)
	b.Var(tableName)

	rows, err := conn.Query(ctx, b.Statement(), b.Args()...)
	if err != nil {
		return nil, err
	}
	var cols []column
	for rows.Next() {
		var col column
		err = rows.Scan(&col.Name, &col.Type, &col.Null, &col.Default, &col.Extra)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, nil
}

func (d *dialect) indices(ctx context.Context, conn orm.Conn, tableName string) ([]index, error) {
	var stmt string
	switch d.Name() {
	case Postgres:
		stmt = "SELECT a.attname AS column_name, i.relname AS index_name, ix.indisprimary as primary FROM pg_class t, pg_class i, pg_index ix, pg_attribute a WHERE t.oid = ix.indrelid AND i.oid = ix.indexrelid AND a.attrelid = t.oid AND a.attnum = ANY(ix.indkey) AND t.relkind = 'r' AND t.relname ="
	default:
		stmt = "SELECT column_name, index_name FROM information_schema.statistics where table_name ="
	}
	b := newBuilder(d, stmt)
	b.Var(tableName)
	rows, err := conn.Query(ctx, b.Statement(), b.Args()...)
	if err != nil {
		return nil, err
	}
	var indices []index
	for rows.Next() {
		var i index
		switch d.Name() {
		case Postgres:
			var primary string
			err = rows.Scan(&i.KeyName, &i.ColumnName, &primary)
			if err != nil {
				return nil, err
			}
			if primary == "t" {
				i.Primary = true
			}
		default:
			err = rows.Scan(&i.KeyName, &i.ColumnName)
			if err != nil {
				return nil, err
			}
			if i.KeyName == "PRIMARY" {
				i.Primary = true
			}
		}
		indices = append(indices, i)
	}
	return indices, nil
}

// column describes an SQL column from an SQL DESCRIBE command
type column struct {
	Name    string
	Type    string
	Null    *string
	Default *string
	Extra   *string
}

// index describes an SQL index from an SQL SHOW INDEX FROM command
type index struct {
	KeyName    string
	ColumnName string
	Primary    bool
}
