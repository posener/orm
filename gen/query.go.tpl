package {{.PackageName}}

import (
	"database/sql"
	"log"

	"github.com/posener/orm/where"

	"{{.Type.ImportPath}}"
)

type Query struct {
	Select *Select
	Where   where.Options
}

func (q *Query) String() string {
	return "SELECT " + q.Select.String() + " FROM {{.Table}} " + q.Where.String()
}

func (q *Query) Exec(db *sql.DB) ([]example.Person, error) {
	// create select statement
	stmt := q.String()
	args := q.Where.Args()
	log.Printf("Query: '%v' %v", stmt, args)
	rows, err := db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.Name}}
	for rows.Next() {
		var i {{.Type.Name}}
		if err := rows.Scan(q.Select.scanArgs(&i)...); err != nil {
			return nil, err
		}
		all = append(all, i)
	}
	return all, rows.Err()
}
