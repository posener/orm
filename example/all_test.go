package example_test

import (
	"testing"
	"time"

	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
	"github.com/posener/orm/example"
	aorm "github.com/posener/orm/example/allsqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypes(t *testing.T) {
	t.Parallel()
	db := prepareAll(t)
	ctx := context.Background()

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

	res, err := db.InsertAll(&a).Exec(ctx)
	require.Nil(t, err)
	affected, err := res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	alls, err := db.Select().Query(ctx)
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))

	assert.Equal(t, a, alls[0])
}

func TestAutoIncrement(t *testing.T) {
	t.Parallel()
	db := prepareAll(t)
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
func TestNotNull(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	orm := prepareAll(t)

	_, err := orm.Insert().SetInt(1).Exec(ctx)
	require.NotNil(t, err)
}

func TestFieldReservedName(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := prepareAll(t)

	res, err := db.InsertAll(&example.All{Select: 42}).Exec(ctx)
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

func prepareAll(t *testing.T) aorm.API {
	t.Helper()
	ctx := context.Background()
	db, err := aorm.Open(":memory:")
	require.Nil(t, err)
	if testing.Verbose() {
		db.Logger(t.Logf)
	}
	_, err = db.Create().Exec(ctx)
	require.Nil(t, err)

	return db
}
