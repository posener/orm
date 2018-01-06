package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/posener/orm"
)

func Conn(dbName string) orm.Conn {
	addr := os.Getenv("MYSQL_ADDR")
	if addr == "" {
		return nil
	}
	ctx := context.Background()
	conn, err := orm.Open("mysql", addr, orm.OptLogger(log.Printf))
	PanicOnErr(err)
	_, err = conn.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName))
	PanicOnErr(err)
	_, err = conn.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s`", dbName))
	PanicOnErr(err)
	_, err = conn.ExecContext(ctx, fmt.Sprintf("USE `%s`", dbName))
	PanicOnErr(err)
	return conn
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
