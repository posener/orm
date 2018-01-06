package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/posener/orm"
	"github.com/stretchr/testify/require"
)

var (
	mySQLAddr   = os.Getenv("MYSQL_ADDR")
	skipCleanup = os.Getenv("SKIP_CLEANUP") != ""
)

func testDBs(t *testing.T, testFunc func(t *testing.T, conn orm.Conn)) {
	t.Helper()
	testNameFixer := strings.NewReplacer(
		"/", "_",
		`"`, "",
		"'", "",
		" ", "_",
	)
	for _, name := range []string{"sqlite3", "mysql"} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			databaseName := "orm_test_" + testNameFixer.Replace(t.Name())
			var (
				conn orm.Conn
				err  error
			)
			switch name {
			case "mysql":
				if mySQLAddr == "" {
					t.Skipf("Mysql environment is not set")
				}
				createMysqlTestDatabase(t, databaseName)
				// create a new connection with the database name in the address
				conn, err = orm.Open(name, mySQLAddr+databaseName)
				require.Nil(t, err)

			case "sqlite3":
				fileName := dbFileName(databaseName)
				os.Remove(fileName)
				conn, err = orm.Open(name, fileName)
				require.Nil(t, err)

			default:
				t.Fatalf("Unknown database name %s", name)
			}

			defer cleanUp(t, name, databaseName)

			defer func() {
				require.Nil(t, conn.Close())
			}()
			if testing.Verbose() {
				conn.Logger(t.Logf)
			}
			testFunc(t, conn)
		})
	}
}

func createMysqlTestDatabase(t *testing.T, databaseName string) {
	ctx := context.Background()
	conn, err := orm.Open("mysql", mySQLAddr)
	require.Nil(t, err)
	defer conn.Close()
	_, err = conn.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", databaseName))
	require.Nil(t, err)
	_, err = conn.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s`", databaseName))
	require.Nil(t, err)
	_, err = conn.ExecContext(ctx, fmt.Sprintf("USE `%s`", databaseName))
	require.Nil(t, err)
}

func cleanUp(t *testing.T, dbType, databaseName string) {
	if skipCleanup {
		return
	}
	switch dbType {
	case "mysql":
		conn, err := sql.Open(dbType, mySQLAddr)
		require.Nil(t, err)
		defer func() {
			require.Nil(t, conn.Close())
		}()
		_, err = conn.Exec(fmt.Sprintf("DROP DATABASE `%s`", databaseName))
		require.Nil(t, err)
	case "sqlite3":
		fileName := dbFileName(databaseName)
		os.Remove(fileName)
	}
}

func dbFileName(databaseName string) string {
	return fmt.Sprintf("/tmp/%s.db", databaseName)
}
