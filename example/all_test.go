package example_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
	"github.com/posener/orm/example"
	aorm "github.com/posener/orm/example/allorm"
	porm "github.com/posener/orm/example/personorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mySQLAddr = os.Getenv("MYSQL_ADDR")

func TestMySQL(t *testing.T) {
	if mySQLAddr == "" {
		t.Skipf("mysql environment is not set")
	}
	db, err := sql.Open("mysql", mySQLAddr)
	require.Nil(t, err)
	defer db.Close()

	t.Run("all", func(t *testing.T) { testAllTable(t, "mysql", db) })
	t.Run("person", func(t *testing.T) { testPersonTable(t, "mysql", db) })
}

func TestTypesSQLite(t *testing.T) {
	t.Parallel()
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err)
	defer db.Close()

	t.Run("all", func(t *testing.T) { testAllTable(t, "sqlite3", db) })
	t.Run("person", func(t *testing.T) { testPersonTable(t, "sqlite3", db) })
}

func testAllTable(t *testing.T, dialect string, db aorm.DB) {
	t.Helper()

	all, err := aorm.New(dialect, db)
	require.Nil(t, err)

	if testing.Verbose() {
		all.Logger(t.Logf)
	}

	resetAllTable(t, db, all)
	t.Run("types", func(t *testing.T) { testTypes(t, all) })
	resetAllTable(t, db, all)
	t.Run("auto increment", func(t *testing.T) { testAutoIncrement(t, all) })
	resetAllTable(t, db, all)
	t.Run("not null", func(t *testing.T) { testNotNull(t, all) })
	resetAllTable(t, db, all)
	t.Run("field with reserved name", func(t *testing.T) { testFieldReservedName(t, all) })
}

func resetAllTable(t *testing.T, db aorm.DB, orm aorm.API) {
	t.Helper()
	_, err := db.ExecContext(context.Background(), "DROP TABLE IF EXISTS `all`")
	require.Nil(t, err)
	_, err = orm.Create().Exec(context.Background())
	require.Nil(t, err)
}

func testTypes(t *testing.T, db aorm.API) {
	t.Helper()

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

	ctx := context.Background()

	res, err := db.InsertAll(&a).Exec(ctx)
	require.Nil(t, err)
	affected, err := res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	alls, err := db.Select().Query(ctx)
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))

	a.Auto = 1 // auto-increment field is set to 1
	assert.Equal(t, a, alls[0])
}

func testAutoIncrement(t *testing.T, db aorm.API) {
	t.Helper()
	ctx := context.Background()

	res, err := db.Insert().SetNotNil("1").Exec(ctx)
	require.Nil(t, err)
	affected, err := res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	res, err = db.Insert().SetNotNil("2").Exec(ctx)
	require.Nil(t, err)
	affected, err = res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	alls, err := db.Select().OrderByAuto(orm.Asc).Query(ctx)
	require.Nil(t, err)
	require.Equal(t, 2, len(alls))

	assert.Equal(t, 1, alls[0].Auto)
	assert.Equal(t, 2, alls[1].Auto)

	alls, err = db.Select().OrderByAuto(orm.Desc).Query(ctx)
	require.Nil(t, err)
	require.Equal(t, 2, len(alls))

	assert.Equal(t, 2, alls[0].Auto)
	assert.Equal(t, 1, alls[1].Auto)
}

// TestNotNull tests that given inserting an empty not null field causes an error
func testNotNull(t *testing.T, db aorm.API) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Insert().SetInt(1).Exec(ctx)
	require.NotNil(t, err)
}

func testFieldReservedName(t *testing.T, db aorm.API) {
	t.Helper()
	ctx := context.Background()

	res, err := db.Insert().SetSelect(42).SetNotNil("not-nil").Exec(ctx)
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)

	query := db.Select().
		SelectSelect().
		Where(aorm.WhereSelect(orm.OpEq, 42).
			Or(aorm.WhereSelectBetween(10, 50)).
			Or(aorm.WhereSelectIn(11, 12))).
		OrderBySelect(orm.Desc).
		GroupBySelect()

	alls, err := query.Query(ctx)
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))
	assert.Equal(t, 42, alls[0].Select)

	res, err = db.Update().SetSelect(11).Where(aorm.WhereSelect(orm.OpEq, 42)).Exec(ctx)
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)

	alls, err = db.Select().SelectSelect().Query(ctx)
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))
	assert.Equal(t, 11, alls[0].Select)

	res, err = db.Delete().Exec(ctx)
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)

	alls, err = db.Select().SelectSelect().Query(ctx)
	require.Nil(t, err)
	require.Equal(t, 0, len(alls))
}

var (
	p1 = example.Person{Name: "moshe", Age: 1}
	p2 = example.Person{Name: "haim", Age: 2}
	p3 = example.Person{Name: "zvika", Age: 3}
)

func testPersonTable(t *testing.T, dialect string, db porm.DB) {
	t.Helper()

	person, err := porm.New(dialect, db)
	require.Nil(t, err)

	if testing.Verbose() {
		person.Logger(t.Logf)
	}

	resetPersonTable(t, db, person)
	t.Run("select", func(t *testing.T) { testPersonSelect(t, person) })
	resetPersonTable(t, db, person)
	t.Run("CRUD", func(t *testing.T) { testCRUD(t, person) })
	resetPersonTable(t, db, person)
	t.Run("count", func(t *testing.T) { testCount(t, person) })
	resetPersonTable(t, db, person)
	t.Run("create if not exists", func(t *testing.T) { testCreateIfNotExists(t, person) })
}

func resetPersonTable(t *testing.T, db porm.DB, orm porm.API) {
	t.Helper()
	_, err := db.ExecContext(context.Background(), "DROP TABLE IF EXISTS `person`")
	require.Nil(t, err)
	_, err = orm.Create().Exec(context.Background())
	require.Nil(t, err)
}

func testPersonSelect(t *testing.T, db porm.API) {
	t.Parallel()
	ctx := context.Background()

	res, err := db.Insert().SetName(p1.Name).SetAge(p1.Age).Exec(ctx)
	require.Nil(t, err, "Failed inserting")
	assertRowsAffected(t, 1, res)
	res, err = db.Insert().SetName(p2.Name).SetAge(p2.Age).Exec(ctx)
	require.Nil(t, err, "Failed inserting")
	assertRowsAffected(t, 1, res)
	res, err = db.InsertPerson(&p3).Exec(ctx)
	require.Nil(t, err, "Failed inserting")
	assertRowsAffected(t, 1, res)

	tests := []struct {
		q    porm.Querier
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
			q:    db.Select().Where(porm.WhereName(orm.OpEq, "moshe")),
			want: []example.Person{p1},
		},
		{
			q:    db.Select().Where(porm.WhereName(orm.OpEq, "moshe").Or(porm.WhereAge(orm.OpEq, 2))),
			want: []example.Person{p1, p2},
		},
		{
			q:    db.Select().Where(porm.WhereName(orm.OpEq, "moshe").And(porm.WhereAge(orm.OpEq, 1))),
			want: []example.Person{p1},
		},
		{
			q: db.Select().Where(porm.WhereName(orm.OpEq, "moshe").And(porm.WhereAge(orm.OpEq, 2))),
		},
		{
			q:    db.Select().Where(porm.WhereAge(orm.OpGE, 2)),
			want: []example.Person{p2, p3},
		},
		{
			q:    db.Select().Where(porm.WhereAge(orm.OpGt, 2)),
			want: []example.Person{p3},
		},
		{
			q:    db.Select().Where(porm.WhereAge(orm.OpLE, 2)),
			want: []example.Person{p1, p2},
		},
		{
			q:    db.Select().Where(porm.WhereAge(orm.OpLt, 2)),
			want: []example.Person{p1},
		},
		{
			q:    db.Select().Where(porm.WhereName(orm.OpNe, "moshe")),
			want: []example.Person{p2, p3},
		},
		{
			q:    db.Select().Where(porm.WhereNameIn("moshe", "haim")),
			want: []example.Person{p1, p2},
		},
		{
			q:    db.Select().Where(porm.WhereAgeBetween(0, 2)),
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
			p, err := tt.q.Query(ctx)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.want, p)
		})
	}
}

func testCRUD(t *testing.T, db porm.API) {
	t.Helper()
	ctx := context.Background()

	// prepareAll dataset
	for _, p := range []example.Person{p1, p2, p3} {
		res, err := db.InsertPerson(&p).Exec(ctx)
		require.Nil(t, err, "Failed inserting")
		assertRowsAffected(t, 1, res)
	}

	ps, err := db.Select().Query(ctx)
	require.Nil(t, err)
	assert.Equal(t, []example.Person{p1, p2, p3}, ps)

	// Test delete
	delete := db.Delete().Where(porm.WhereName(orm.OpEq, "moshe"))
	res, err := delete.Exec(ctx)
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)
	ps, err = db.Select().Query(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []example.Person{p2, p3}, ps)
	ps, err = db.Select().Where(porm.WhereName(orm.OpEq, "moshe")).Query(ctx)
	require.Nil(t, err)
	assert.Equal(t, []example.Person(nil), ps)

	// Test Update
	update := db.Update().SetName("Jonney").Where(porm.WhereName(orm.OpEq, "zvika"))
	res, err = update.Exec(ctx)
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)

	ps, err = db.Select().Where(porm.WhereName(orm.OpEq, "Jonney")).Query(ctx)
	require.Nil(t, err)
	assert.Equal(t, []example.Person{{Name: "Jonney", Age: 3}}, ps)
}

func testCount(t *testing.T, db porm.API) {
	t.Helper()
	ctx := context.Background()

	for i := 0; i < 100; i++ {
		res, err := db.InsertPerson(&example.Person{Name: fmt.Sprintf("Jim %d", i), Age: i / 5}).Exec(ctx)
		require.Nil(t, err, "Failed inserting")
		assertRowsAffected(t, 1, res)
	}

	tests := []struct {
		q    porm.Counter
		want []porm.PersonCount
	}{
		{
			q:    db.Select(),
			want: []porm.PersonCount{{Count: 100}},
		},
		{
			q:    db.Select().Where(porm.WhereAge(orm.OpLt, 10)),
			want: []porm.PersonCount{{Count: 50}},
		},
		{
			q: db.Select().SelectAge().GroupByAge().Where(porm.WhereAgeIn(1, 3, 12)),
			want: []porm.PersonCount{
				{Person: example.Person{Age: 1}, Count: 5},
				{Person: example.Person{Age: 3}, Count: 5},
				{Person: example.Person{Age: 12}, Count: 5},
			},
		},
		{
			q: db.Select().SelectAge().GroupByAge().Where(porm.WhereAgeIn(1, 3, 12)).OrderByAge(orm.Desc),
			want: []porm.PersonCount{
				{Person: example.Person{Age: 12}, Count: 5},
				{Person: example.Person{Age: 3}, Count: 5},
				{Person: example.Person{Age: 1}, Count: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%+v", tt.q), func(t *testing.T) {
			got, err := tt.q.Count(ctx)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func testCreateIfNotExists(t *testing.T, db porm.API) {
	t.Helper()
	ctx := context.Background()

	_, err := db.Create().IfNotExists().Exec(ctx)
	require.Nil(t, err)

	_, err = db.Create().Exec(ctx)
	require.NotNil(t, err)
}

func assertRowsAffected(t *testing.T, wantRows int64, result sql.Result) {
	gotRows, err := result.RowsAffected()
	require.Nil(t, err)
	assert.Equal(t, wantRows, gotRows)
}

func TestNew(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err)
	defer db.Close()

	orm, err := porm.New("sqlite3", db)
	require.Nil(t, err)

	if testing.Verbose() {
		orm.Logger(t.Logf)
	}
	_, err = orm.Create().Exec(ctx)
	require.Nil(t, err)
}
