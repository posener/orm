package examples

import (
	"os"

	"github.com/posener/orm"
)

func conn() orm.DB {
	addr := os.Getenv("MYSQL_ADDR")
	if addr == "" {
		return nil
	}
	conn, err := orm.Open("mysql", addr)
	if err != nil {
		panic(err)
	}
	return conn
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
