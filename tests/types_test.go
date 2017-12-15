package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mySQLAddr = os.Getenv("MYSQL_ADDR")

var dbNames = []string{"sqlite3", "mysql"}

func TestTypes(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := allDB(t, conn)

		a := &All{
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

		aGot, err := db.Insert().InsertAll(a).Exec()
		require.Nil(t, err)
		assert.Equal(t, 1, aGot.Auto)
		a.Auto = aGot.Auto
		assert.Equal(t, a, aGot)

		alls, err := db.Select().Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(alls))

		assert.Equal(t, a, &alls[0])
	})
}

func TestAutoIncrement(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := allDB(t, conn)

		a1, err := db.Insert().SetNotNil("1").Exec()
		require.Nil(t, err)
		assert.Equal(t, 1, a1.Auto)

		a2, err := db.Insert().SetNotNil("2").Exec()
		require.Nil(t, err)
		assert.Equal(t, 2, a2.Auto)

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

		_, err := db.Insert().SetInt(1).Exec()
		require.NotNil(t, err)
	})
}

func TestFieldReservedName(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := allDB(t, conn)

		_, err := db.Insert().SetSelect(42).SetNotNil("not-nil").Exec()
		require.Nil(t, err)

		query := db.Select().
			SelectSelect().
			Where(db.Where().Select(orm.OpEq, 42).
				Or(db.Where().SelectBetween(10, 50)).
				Or(db.Where().SelectIn(11, 12))).
			OrderBySelect(orm.Desc).
			GroupBySelect()

		alls, err := query.Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(alls))
		assert.Equal(t, 42, alls[0].Select)

		res, err := db.Update().SetSelect(11).Where(db.Where().Select(orm.OpEq, 42)).Exec()
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
	p1 = Person{Name: "moshe", Age: 1}
	p2 = Person{Name: "haim", Age: 2}
	p3 = Person{Name: "zvika", Age: 3}
)

func TestPersonSelect(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := personDB(t, conn)

		_, err := db.Insert().SetName(p1.Name).SetAge(p1.Age).Exec()
		require.Nil(t, err)
		_, err = db.Insert().SetName(p2.Name).SetAge(p2.Age).Exec()
		require.Nil(t, err)
		_, err = db.Insert().InsertPerson(&p3).Exec()
		require.Nil(t, err)

		tests := []struct {
			q    *PersonSelectBuilder
			want []Person
		}{
			{
				q:    db.Select(),
				want: []Person{p1, p2, p3},
			},
			{
				q:    db.Select().SelectName(),
				want: []Person{{Name: "moshe"}, {Name: "haim"}, {Name: "zvika"}},
			},
			{
				q:    db.Select().SelectAge(),
				want: []Person{{Age: 1}, {Age: 2}, {Age: 3}},
			},
			{
				q:    db.Select().SelectAge().SelectName(),
				want: []Person{p1, p2, p3},
			},
			{
				q:    db.Select().SelectName().SelectAge(),
				want: []Person{p1, p2, p3},
			},
			{
				q:    db.Select().Where(db.Where().Name(orm.OpEq, "moshe")),
				want: []Person{p1},
			},
			{
				q:    db.Select().Where(db.Where().Name(orm.OpEq, "moshe").Or(db.Where().Age(orm.OpEq, 2))),
				want: []Person{p1, p2},
			},
			{
				q:    db.Select().Where(db.Where().Name(orm.OpEq, "moshe").And(db.Where().Age(orm.OpEq, 1))),
				want: []Person{p1},
			},
			{
				q: db.Select().Where(db.Where().Name(orm.OpEq, "moshe").And(db.Where().Age(orm.OpEq, 2))),
			},
			{
				q:    db.Select().Where(db.Where().Age(orm.OpGE, 2)),
				want: []Person{p2, p3},
			},
			{
				q:    db.Select().Where(db.Where().Age(orm.OpGt, 2)),
				want: []Person{p3},
			},
			{
				q:    db.Select().Where(db.Where().Age(orm.OpLE, 2)),
				want: []Person{p1, p2},
			},
			{
				q:    db.Select().Where(db.Where().Age(orm.OpLt, 2)),
				want: []Person{p1},
			},
			{
				q:    db.Select().Where(db.Where().Name(orm.OpNe, "moshe")),
				want: []Person{p2, p3},
			},
			{
				q:    db.Select().Where(db.Where().NameIn("moshe", "haim")),
				want: []Person{p1, p2},
			},
			{
				q:    db.Select().Where(db.Where().AgeBetween(0, 2)),
				want: []Person{p1, p2},
			},
			{
				q:    db.Select().Limit(2),
				want: []Person{p1, p2},
			},
			{
				q:    db.Select().Page(1, 1),
				want: []Person{p2},
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

		// prepareAll dataset
		for _, p := range []Person{p1, p2, p3} {
			_, err := db.Insert().InsertPerson(&p).Exec()
			require.Nil(t, err)
		}

		ps, err := db.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, []Person{p1, p2, p3}, ps)

		// Test delete
		delete := db.Delete().Where(db.Where().Name(orm.OpEq, "moshe"))
		res, err := delete.Exec()
		require.Nil(t, err)
		assertRowsAffected(t, 1, res)
		ps, err = db.Select().Query()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []Person{p2, p3}, ps)
		ps, err = db.Select().Where(db.Where().Name(orm.OpEq, "moshe")).Query()
		require.Nil(t, err)
		assert.Equal(t, []Person(nil), ps)

		// Test Update
		update := db.Update().SetName("Jonney").Where(db.Where().Name(orm.OpEq, "zvika"))
		res, err = update.Exec()
		require.Nil(t, err)
		assertRowsAffected(t, 1, res)

		ps, err = db.Select().Where(db.Where().Name(orm.OpEq, "Jonney")).Query()
		require.Nil(t, err)
		assert.Equal(t, []Person{{Name: "Jonney", Age: 3}}, ps)
	})
}

func TestCount(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := personDB(t, conn)

		for i := 0; i < 100; i++ {
			_, err := db.Insert().InsertPerson(&Person{Name: fmt.Sprintf("Jim %d", i), Age: i / 5}).Exec()
			require.Nil(t, err)
		}

		tests := []struct {
			q    *PersonSelectBuilder
			want []PersonCount
		}{
			{
				q:    db.Select(),
				want: []PersonCount{{Count: 100}},
			},
			{
				q:    db.Select().Where(db.Where().Age(orm.OpLt, 10)),
				want: []PersonCount{{Count: 50}},
			},
			{
				q: db.Select().SelectAge().GroupByAge().Where(db.Where().AgeIn(1, 3, 12)),
				want: []PersonCount{
					{Person: Person{Age: 1}, Count: 5},
					{Person: Person{Age: 3}, Count: 5},
					{Person: Person{Age: 12}, Count: 5},
				},
			},
			{
				q: db.Select().SelectAge().GroupByAge().Where(db.Where().AgeIn(1, 3, 12)).OrderByAge(orm.Desc),
				want: []PersonCount{
					{Person: Person{Age: 12}, Count: 5},
					{Person: Person{Age: 3}, Count: 5},
					{Person: Person{Age: 1}, Count: 5},
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

		err := db.Create().IfNotExists().Exec()
		require.Nil(t, err)

		err = db.Create().Exec()
		require.NotNil(t, err)
	})
}

func TestFirst(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		db := personDB(t, conn)

		_, err := db.Select().First()
		assert.Equal(t, orm.ErrNotFound, err)

		smith := Person{Name: "Smith", Age: 99}
		_, err = db.Insert().InsertPerson(&smith).Exec()
		require.Nil(t, err)

		got, err := db.Select().First()
		assert.Nil(t, err)
		assert.Equal(t, &smith, got)

		john := Person{Name: "John", Age: 12}
		_, err = db.Insert().InsertPerson(&john).Exec()
		require.Nil(t, err)

		got, err = db.Select().First()
		assert.Nil(t, err)
		assert.Equal(t, &smith, got)

		got, err = db.Select().Where(db.Where().Age(orm.OpEq, 12)).First()
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

	orm, err := NewPersonORM("sqlite3", db)
	require.Nil(t, err)

	if testing.Verbose() {
		orm.Logger(t.Logf)
	}
	err = orm.Create().Exec()
	require.Nil(t, err)
}

func personDB(t *testing.T, conn conn) PersonORM {
	t.Helper()
	db, err := NewPersonORM(conn.name, conn)
	require.Nil(t, err)
	if testing.Verbose() {
		db.Logger(t.Logf)
	}
	err = db.Create().Exec()
	require.Nil(t, err)
	return db
}

func allDB(t *testing.T, conn conn) AllORM {
	t.Helper()
	db, err := NewAllORM(conn.name, conn)
	require.Nil(t, err)
	if testing.Verbose() {
		db.Logger(t.Logf)
	}
	err = db.Create().Exec()
	require.Nil(t, err)
	return db
}

type conn struct {
	name string
	orm.DB
}

func testDBs(t *testing.T, testFunc func(t *testing.T, conn conn)) {
	t.Helper()
	for _, name := range dbNames {
		t.Run(name, func(t *testing.T) {
			var (
				s   *sql.DB
				err error
			)
			switch name {
			case "mysql":
				if mySQLAddr == "" {
					t.Skipf("mysql environment is not set")
				}
				s, err = sql.Open(name, mySQLAddr)
				require.Nil(t, err)
				_, err = s.Exec("DROP DATABASE IF EXISTS test")
				require.Nil(t, err)
				_, err = s.Exec("CREATE DATABASE test")
				require.Nil(t, err)
				_, err = s.Exec("USE test")
				require.Nil(t, err)

			case "sqlite3":
				os.Remove("test.db")
				s, err = sql.Open(name, "test.db")
				require.Nil(t, err)

			default:
				panic("unknown db")
			}

			defer s.Close()
			testFunc(t, conn{name: name, DB: s})

		})
	}
}
