package dialect

import (
	"context"
	"fmt"
	"log"

	"github.com/posener/orm"
)

type column struct {
	Field   string
	Type    string
	Null    *string
	Key     *string
	Default *string
	Extra   *string
}

func describeTable(ctx context.Context, db orm.DB, tableName string) ([]column, error) {
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
