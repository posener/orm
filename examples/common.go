package examples

import (
	"database/sql"
	"os"
)

func conn() *sql.DB {
	addr := os.Getenv("MYSQL_ADDR")
	if addr == "" {
		return nil
	}
	db, err := sql.Open("mysql", addr)
	if err != nil {
		panic(err)
	}
	return db
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
