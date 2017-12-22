package examples

import (
	"database/sql"
	"fmt"
	"os"
)

func conn() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_ADDR"))
	if err != nil {
		panic(err)
	}
	for _, table := range []string{"simple", "othermany", "one", "otherone"} {
		_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			panic(err)
		}
	}
	return db
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
