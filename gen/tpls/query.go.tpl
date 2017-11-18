package {{.PackageName}}

import (
	"database/sql"
	"log"

	"{{.Type.ImportPath}}"
)

func (q *Query) String() string {
	return "SELECT " + q.sel.String() + " FROM {{.Table}} " + q.where.String()
}

func (q *Query) Exec(db *sql.DB) ([]example.Person, error) {
	// create select statement
	stmt := q.String()
	args := q.where.Args()
	log.Printf("Query: '%v' %v", stmt, args)
	rows, err := db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.FullName}}
	for rows.Next() {
		var i {{.Type.FullName}}
		if err := rows.Scan(q.sel.scanArgs(&i)...); err != nil {
			return nil, err
		}
		all = append(all, i)
	}
	return all, rows.Err()
}
