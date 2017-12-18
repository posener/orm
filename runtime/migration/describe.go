package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/posener/orm"
)

// Load loads a table from an existing database
func Load(ctx context.Context, db orm.DB, tableName string) (*Table, error) {
	columns, err := describe(ctx, db, tableName)
	if err != nil {
		return nil, err
	}
	t := new(Table)
	for _, col := range columns {
		t.Columns = append(t.Columns, Column{
			Name:    col.Field,
			SQLType: col.Type,
		})
	}
	return t, nil
}

// describe returns all columns of a table by doing an SQL query
func describe(ctx context.Context, db orm.DB, tableName string) ([]column, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("DESCRIBE `%s`", tableName))
	if err != nil {
		log.Print(err)
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

// column describes an SQL column from an SQL DESCRIBE command
type column struct {
	Field   string
	Type    string
	Null    *string
	Key     *string
	Default *string
	Extra   *string
}
