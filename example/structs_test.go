package example_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
	"github.com/posener/orm/example"
	"github.com/posener/orm/example/allorm"
	"github.com/posener/orm/example/bookorm"
	"github.com/posener/orm/example/loanerorm"
	"github.com/posener/orm/example/personorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mySQLAddr = os.Getenv("MYSQL_ADDR")

var dbNames = []string{"sqlite3", "mysql"}

func TestTypes(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := allDB(t, conn)
		defer db.Close()

		a := example.All{
			Int:   1,
			Int8:  int8(2),
			Int16: int16(3),
			Int32: int32(4),
			Int64: int64(5),

			UInt:   uint(6),
			UInt8:  uint8(7),
			UInt16: uint16(8),
			UInt32: uint32(9),
			UInt64: uint64(10),

			String:        "11",
			Bytes:         []byte("12"),
			VarCharString: "13",
			VarCharByte:   []byte("14"),

			Bool: true,

			Time: time.Now().Round(time.Millisecond).UTC(),

			NotNil: "not-nil",
		}

		a.PInt = &a.Int
		a.PInt8 = &a.Int8
		a.PInt16 = &a.Int16
		a.PInt32 = &a.Int32
		a.PInt64 = &a.Int64

		a.PUInt = &a.UInt
		a.PUInt8 = &a.UInt8
		a.PUInt16 = &a.UInt16
		a.PUInt32 = &a.UInt32
		a.PUInt64 = &a.UInt64

		a.PString = &a.String
		a.PBytes = &a.Bytes
		a.PVarCharString = &a.VarCharString
		a.PVarCharByte = &a.VarCharByte

		a.PBool = &a.Bool
		a.PTime = &a.Time

		res, err := db.Insert().InsertAll(&a).Exec()
		require.Nil(t, err)
		affected, err := res.RowsAffected()
		require.Nil(t, err)
		require.Equal(t, int64(1), affected)

		alls, err := db.Select().Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(alls))

		a.Auto = 1 // auto-increment field is set to 1
		assert.Equal(t, a, alls[0])
	})
}

func TestAutoIncrement(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := allDB(t, conn)
		defer db.Close()

		res, err := db.Insert().SetNotNil("1").Exec()
		require.Nil(t, err)
		affected, err := res.RowsAffected()
		require.Nil(t, err)
		require.Equal(t, int64(1), affected)

		res, err = db.Insert().SetNotNil("2").Exec()
		require.Nil(t, err)
		affected, err = res.RowsAffected()
		require.Nil(t, err)
		require.Equal(t, int64(1), affected)

		alls, err := db.Select().OrderByAuto(orm.Asc).Query()
		require.Nil(t, err)
		require.Equal(t, 2, len(alls))

		assert.Equal(t, 1, alls[0].Auto)
		assert.Equal(t, 2, alls[1].Auto)

		alls, err = db.Select().OrderByAuto(orm.Desc).Query()
		require.Nil(t, err)
		require.Equal(t, 2, len(alls))

		assert.Equal(t, 2, alls[0].Auto)
		assert.Equal(t, 1, alls[1].Auto)
	})
}

// TestNotNull tests that given inserting an empty not null field causes an error
func TestNotNull(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := allDB(t, conn)
		defer db.Close()

		_, err := db.Insert().SetInt(1).Exec()
		require.NotNil(t, err)
	})
}

func TestFieldReservedName(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := allDB(t, conn)
		defer db.Close()

		res, err := db.Insert().SetSelect(42).SetNotNil("not-nil").Exec()
		require.Nil(t, err)
		assertRowsAffected(t, 1, res)

		query := db.Select().
			SelectSelect().
			Where(allorm.WhereSelect(orm.OpEq, 42).
				Or(allorm.WhereSelectBetween(10, 50)).
				Or(allorm.WhereSelectIn(11, 12))).
			OrderBySelect(orm.Desc).
			GroupBySelect()

		alls, err := query.Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(alls))
		assert.Equal(t, 42, alls[0].Select)

		res, err = db.Update().SetSelect(11).Where(allorm.WhereSelect(orm.OpEq, 42)).Exec()
		require.Nil(t, err)
		assertRowsAffected(t, 1, res)

		alls, err = db.Select().SelectSelect().Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(alls))
		assert.Equal(t, 11, alls[0].Select)

		res, err = db.Delete().Exec()
		require.Nil(t, err)
		assertRowsAffected(t, 1, res)

		alls, err = db.Select().SelectSelect().Query()
		require.Nil(t, err)
		require.Equal(t, 0, len(alls))
	})
}

var (
	p1 = example.Person{Name: "moshe", Age: 1}
	p2 = example.Person{Name: "haim", Age: 2}
	p3 = example.Person{Name: "zvika", Age: 3}
)

func TestPersonSelect(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := personDB(t, conn)
		defer db.Close()

		res, err := db.Insert().SetName(p1.Name).SetAge(p1.Age).Exec()
		require.Nil(t, err, "Failed inserting")
		assertRowsAffected(t, 1, res)
		res, err = db.Insert().SetName(p2.Name).SetAge(p2.Age).Exec()
		require.Nil(t, err, "Failed inserting")
		assertRowsAffected(t, 1, res)
		res, err = db.Insert().InsertPerson(&p3).Exec()
		require.Nil(t, err, "Failed inserting")
		assertRowsAffected(t, 1, res)

		tests := []struct {
			q    personorm.Querier
			want []example.Person
		}{
			{
				q:    db.Select(),
				want: []example.Person{p1, p2, p3},
			},
			{
				q:    db.Select().SelectName(),
				want: []example.Person{{Name: "moshe"}, {Name: "haim"}, {Name: "zvika"}},
			},
			{
				q:    db.Select().SelectAge(),
				want: []example.Person{{Age: 1}, {Age: 2}, {Age: 3}},
			},
			{
				q:    db.Select().SelectAge().SelectName(),
				want: []example.Person{p1, p2, p3},
			},
			{
				q:    db.Select().SelectName().SelectAge(),
				want: []example.Person{p1, p2, p3},
			},
			{
				q:    db.Select().Where(personorm.WhereName(orm.OpEq, "moshe")),
				want: []example.Person{p1},
			},
			{
				q:    db.Select().Where(personorm.WhereName(orm.OpEq, "moshe").Or(personorm.WhereAge(orm.OpEq, 2))),
				want: []example.Person{p1, p2},
			},
			{
				q:    db.Select().Where(personorm.WhereName(orm.OpEq, "moshe").And(personorm.WhereAge(orm.OpEq, 1))),
				want: []example.Person{p1},
			},
			{
				q: db.Select().Where(personorm.WhereName(orm.OpEq, "moshe").And(personorm.WhereAge(orm.OpEq, 2))),
			},
			{
				q:    db.Select().Where(personorm.WhereAge(orm.OpGE, 2)),
				want: []example.Person{p2, p3},
			},
			{
				q:    db.Select().Where(personorm.WhereAge(orm.OpGt, 2)),
				want: []example.Person{p3},
			},
			{
				q:    db.Select().Where(personorm.WhereAge(orm.OpLE, 2)),
				want: []example.Person{p1, p2},
			},
			{
				q:    db.Select().Where(personorm.WhereAge(orm.OpLt, 2)),
				want: []example.Person{p1},
			},
			{
				q:    db.Select().Where(personorm.WhereName(orm.OpNe, "moshe")),
				want: []example.Person{p2, p3},
			},
			{
				q:    db.Select().Where(personorm.WhereNameIn("moshe", "haim")),
				want: []example.Person{p1, p2},
			},
			{
				q:    db.Select().Where(personorm.WhereAgeBetween(0, 2)),
				want: []example.Person{p1, p2},
			},
			{
				q:    db.Select().Limit(2),
				want: []example.Person{p1, p2},
			},
			{
				q:    db.Select().Page(1, 1),
				want: []example.Person{p2},
			},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%+v", tt.q), func(t *testing.T) {
				p, err := tt.q.Query()
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, tt.want, p)
			})
		}
	})
}

func TestCRUD(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := personDB(t, conn)
		defer db.Close()

		// prepareAll dataset
		for _, p := range []example.Person{p1, p2, p3} {
			res, err := db.Insert().InsertPerson(&p).Exec()
			require.Nil(t, err, "Failed inserting")
			assertRowsAffected(t, 1, res)
		}

		ps, err := db.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, []example.Person{p1, p2, p3}, ps)

		// Test delete
		delete := db.Delete().Where(personorm.WhereName(orm.OpEq, "moshe"))
		res, err := delete.Exec()
		require.Nil(t, err)
		assertRowsAffected(t, 1, res)
		ps, err = db.Select().Query()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []example.Person{p2, p3}, ps)
		ps, err = db.Select().Where(personorm.WhereName(orm.OpEq, "moshe")).Query()
		require.Nil(t, err)
		assert.Equal(t, []example.Person(nil), ps)

		// Test Update
		update := db.Update().SetName("Jonney").Where(personorm.WhereName(orm.OpEq, "zvika"))
		res, err = update.Exec()
		require.Nil(t, err)
		assertRowsAffected(t, 1, res)

		ps, err = db.Select().Where(personorm.WhereName(orm.OpEq, "Jonney")).Query()
		require.Nil(t, err)
		assert.Equal(t, []example.Person{{Name: "Jonney", Age: 3}}, ps)
	})
}

func TestCount(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := personDB(t, conn)
		defer db.Close()

		for i := 0; i < 100; i++ {
			res, err := db.Insert().InsertPerson(&example.Person{Name: fmt.Sprintf("Jim %d", i), Age: i / 5}).Exec()
			require.Nil(t, err, "Failed inserting")
			assertRowsAffected(t, 1, res)
		}

		tests := []struct {
			q    personorm.Counter
			want []personorm.PersonCount
		}{
			{
				q:    db.Select(),
				want: []personorm.PersonCount{{Count: 100}},
			},
			{
				q:    db.Select().Where(personorm.WhereAge(orm.OpLt, 10)),
				want: []personorm.PersonCount{{Count: 50}},
			},
			{
				q: db.Select().SelectAge().GroupByAge().Where(personorm.WhereAgeIn(1, 3, 12)),
				want: []personorm.PersonCount{
					{Person: example.Person{Age: 1}, Count: 5},
					{Person: example.Person{Age: 3}, Count: 5},
					{Person: example.Person{Age: 12}, Count: 5},
				},
			},
			{
				q: db.Select().SelectAge().GroupByAge().Where(personorm.WhereAgeIn(1, 3, 12)).OrderByAge(orm.Desc),
				want: []personorm.PersonCount{
					{Person: example.Person{Age: 12}, Count: 5},
					{Person: example.Person{Age: 3}, Count: 5},
					{Person: example.Person{Age: 1}, Count: 5},
				},
			},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%+v", tt.q), func(t *testing.T) {
				got, err := tt.q.Count()
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, tt.want, got)
			})
		}
	})
}

func TestCreateIfNotExists(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := personDB(t, conn)
		defer db.Close()

		_, err := db.Create().IfNotExists().Exec()
		require.Nil(t, err)

		_, err = db.Create().Exec()
		require.NotNil(t, err)
	})
}

func TestFirst(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := personDB(t, conn)
		defer db.Close()

		_, err := db.Select().First()
		assert.Equal(t, orm.ErrNotFound, err)

		smith := example.Person{Name: "Smith", Age: 99}
		_, err = db.Insert().InsertPerson(&smith).Exec()
		require.Nil(t, err)

		got, err := db.Select().First()
		assert.Nil(t, err)
		assert.Equal(t, &smith, got)

		john := example.Person{Name: "John", Age: 12}
		_, err = db.Insert().InsertPerson(&john).Exec()
		require.Nil(t, err)

		got, err = db.Select().First()
		assert.Nil(t, err)
		assert.Equal(t, &smith, got)

		got, err = db.Select().Where(personorm.WhereAge(orm.OpEq, 12)).First()
		assert.Nil(t, err)
		assert.Equal(t, &john, got)

		got, err = db.Select().OrderByAge(orm.Asc).First()
		assert.Nil(t, err)
		assert.Equal(t, &john, got)

		got, err = db.Select().OrderByAge(orm.Desc).First()
		assert.Nil(t, err)
		assert.Equal(t, &smith, got)
	})
}

func assertRowsAffected(t *testing.T, wantRows int64, result sql.Result) {
	gotRows, err := result.RowsAffected()
	require.Nil(t, err)
	assert.Equal(t, wantRows, gotRows)
}

func TestNew(t *testing.T) {
	t.Parallel()
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err)
	defer db.Close()

	orm, err := personorm.New("sqlite3", db)
	require.Nil(t, err)

	if testing.Verbose() {
		orm.Logger(t.Logf)
	}
	_, err = orm.Create().Exec()
	require.Nil(t, err)
}

func TestReferences(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		book, loaner := bookDB(t, conn)
		defer book.Close()

		if conn.name != "sqlite3" { // this is not enforced in sqlite3
			_, err := loaner.Insert().InsertLoaner(&example.Loaner{Name: "James", Book: &example.Book{ID: -10}}).Exec()
			require.NotNil(t, err)
		}

		b := &example.Book{Name: "The Hitchhiker's Guide to the Galaxy", Year: 1979}

		res, err := book.Insert().InsertBook(b).Exec()
		require.Nil(t, err)
		b.ID, err = res.LastInsertId()
		require.Nil(t, err)

		l := &example.Loaner{Name: "James", Book: b}

		res, err = loaner.Insert().InsertLoaner(l).Exec()
		require.Nil(t, err)
		l.ID, err = res.LastInsertId()
		require.Nil(t, err)

		// query without join, Loaner.Book should be nil
		ls, err := loaner.Select().Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(ls)) {
			assert.Equal(t, l.Name, ls[0].Name)
			assert.Nil(t, ls[0].Book)
		}

		// query with join, Loaner.Book should be filled with book's properties
		ls, err = loaner.Select().JoinBook(book.Select().Scanner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(ls)) {
			assert.Equal(t, l, &ls[0])
		}
	})
}

func personDB(t *testing.T, conn conn) personorm.API {
	t.Helper()
	db, err := personorm.New(conn.name, conn)
	require.Nil(t, err)
	if testing.Verbose() {
		db.Logger(t.Logf)
	}
	_, err = db.Create().Exec()
	require.Nil(t, err)
	return db
}

func allDB(t *testing.T, conn conn) allorm.API {
	t.Helper()
	db, err := allorm.New(conn.name, conn)
	require.Nil(t, err)
	if testing.Verbose() {
		db.Logger(t.Logf)
	}
	_, err = db.Create().Exec()
	require.Nil(t, err)
	return db
}

func bookDB(t *testing.T, conn conn) (bookorm.API, loanerorm.API) {
	t.Helper()
	bookDB, err := bookorm.New(conn.name, conn)
	require.Nil(t, err)
	loanerDB, err := loanerorm.New(conn.name, conn)
	require.Nil(t, err)
	if testing.Verbose() {
		bookDB.Logger(t.Logf)
		loanerDB.Logger(t.Logf)
	}

	_, err = bookDB.Create().Exec()
	require.Nil(t, err)

	_, err = loanerDB.Create().Exec()
	require.Nil(t, err)

	return bookDB, loanerDB
}

type conn struct {
	name string
	orm.DB
}

func testDBs(t *testing.T, testFunc func(t *testing.T, conn conn)) {
	t.Helper()
	for _, name := range dbNames {
		t.Run(name, func(t *testing.T) {
			switch name {
			case "mysql":
				if mySQLAddr == "" {
					t.Skipf("mysql environment is not set")
				}
				s, err := sql.Open(name, mySQLAddr)
				require.Nil(t, err)
				for _, table := range []string{"all", "person", "loaner", "book"} {
					_, err = s.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table))
					require.Nil(t, err)
				}
				testFunc(t, conn{name: name, DB: s})

			case "sqlite3":
				s, err := sql.Open(name, ":memory:")
				require.Nil(t, err)
				testFunc(t, conn{name: name, DB: s})

			default:
				panic("unknown db")
			}

		})
	}
}
