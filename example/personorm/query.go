package personorm

import (
	"database/sql"
	"log"

	"github.com/posener/orm/example"
)

type Query struct {
	Columns *SelectColumns
	Where   *WhereOptions
}

func (q *Query) String() string {
	return "SELECT " + q.Columns.String() + " FROM person " + q.Where.String()
}

func (q *Query) Args() []interface{} {
	return q.Where.Args()
}

func (q *Query) Exec(db *sql.DB) ([]example.Person, error) {
	// create select statement
	stmt := q.String()
	args := q.Args()
	log.Printf("Query: '%v' %v", stmt, args)
	rows, err := db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []example.Person
	for rows.Next() {
		var i example.Person
		if err := rows.Scan(q.Columns.ScanArgs(&i)...); err != nil {
			return nil, err
		}
		all = append(all, i)
	}
	return all, rows.Err()
}
