package tests

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
	"github.com/stretchr/testify/require"
)

var (
	mySQLAddr    = os.Getenv("MYSQL_ADDR")
	postgresAddr = os.Getenv("POSTGRES_ADDR")
	skipCleanup  = os.Getenv("SKIP_CLEANUP") != ""
)

func testDBs(t *testing.T, testFunc func(t *testing.T, conn orm.Conn)) {
	t.Helper()
	testNameFixer := strings.NewReplacer(
		"/", "_",
		`"`, "",
		"'", "",
		" ", "_",
	)
	var options []orm.Option
	if testing.Verbose() {
		options = append(options, orm.OptLogger(t.Logf))
	}

	for _, name := range []string{
		"sqlite3",
		"mysql",
		"postgres",
	} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			databaseName := "orm_test_" + testNameFixer.Replace(t.Name())
			var address string
			var (
				conn    orm.Conn
				err     error
				cleanup func()
			)
			switch name {
			case "mysql":
				if mySQLAddr == "" {
					t.Skipf("Mysql environment is not set")
				}
				cleanup = createMysqlTestDatabase(t, databaseName)
				address = mySQLAddr + databaseName

			case "postgres":
				if postgresAddr == "" {
					t.Skipf("Postgres environment is not set")
				}
				cleanup = createPostgresTestDatabase(t, databaseName)
				u, err := url.Parse(postgresAddr)
				require.Nil(t, err)
				u.Path = databaseName
				address = u.String()

			case "sqlite3":
				address = dbFileName(databaseName)
				os.Remove(address)
				cleanup = func() {
					os.Remove(address)
				}

			default:
				t.Fatalf("Unknown database name %s", name)
			}

			conn, err = orm.Open(name, address, options...)
			require.Nil(t, err)

			defer func() {
				if skipCleanup {
					return
				}
				cleanup()
			}()

			defer func() {
				require.Nil(t, conn.Close())
			}()
			testFunc(t, conn)
		})
	}
}

func createMysqlTestDatabase(t *testing.T, databaseName string) func() {
	ctx := context.Background()
	conn, err := orm.Open("mysql", mySQLAddr)
	require.Nil(t, err)
	defer conn.Close()
	_, err = conn.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", databaseName))
	require.Nil(t, err)
	_, err = conn.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s`", databaseName))
	require.Nil(t, err)
	return func() {
		conn, err := sql.Open("mysql", mySQLAddr)
		require.Nil(t, err)
		defer func() {
			require.Nil(t, conn.Close())
		}()
		_, err = conn.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", databaseName))
		require.Nil(t, err)
	}
}

func createPostgresTestDatabase(t *testing.T, databaseName string) func() {
	ctx := context.Background()
	conn, err := orm.Open("postgres", postgresAddr)
	require.Nil(t, err)
	defer conn.Close()
	conn.ExecContext(ctx, fmt.Sprintf(`DROP DATABASE "%s"`, databaseName))
	conn.ExecContext(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, databaseName))
	return func() {
		conn, err := sql.Open("postgres", postgresAddr)
		require.Nil(t, err)
		defer func() {
			require.Nil(t, conn.Close())
		}()
		_, err = conn.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, databaseName))
		require.Nil(t, err)
	}
}

func dbFileName(databaseName string) string {
	return fmt.Sprintf("/tmp/%s.db", databaseName)
}
