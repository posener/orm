package example_test

import (
	"testing"
	"time"

	"github.com/posener/orm"
	"github.com/posener/orm/example"
	aorm "github.com/posener/orm/example/allsqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	t.Parallel()
	db := prepareAll(t)

	tm := time.Now().Round(time.Millisecond).UTC()

	res, err := db.Insert().SetTime(tm).SetString("not null").Exec()
	require.Nil(t, err)
	affected, err := res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	alls, err := db.Select().Query()
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))

	assert.Equal(t, tm, alls[0].Time)
}

func TestAutoIncrement(t *testing.T) {
	t.Parallel()
	db := prepareAll(t)

	res, err := db.Insert().SetString("1").Exec()
	require.Nil(t, err)
	affected, err := res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	res, err = db.Insert().SetString("2").Exec()
	require.Nil(t, err)
	affected, err = res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	alls, err := db.Select().OrderByInt(orm.Asc).Query()
	require.Nil(t, err)
	require.Equal(t, 2, len(alls))

	assert.Equal(t, 1, alls[0].Int)
	assert.Equal(t, 2, alls[1].Int)
}

// TestNotNull tests that given inserting an empty not null field causes an error
func TestNotNull(t *testing.T) {
	t.Parallel()
	orm := prepareAll(t)

	_, err := orm.Insert().SetInt(1).Exec()
	require.NotNil(t, err)
}

func TestFieldReservedName(t *testing.T) {
	t.Parallel()
	db := prepareAll(t)

	res, err := db.InsertAll(&example.All{Select: 42}).Exec()
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)

	query := db.Select().
		SelectSelect().
		Where(aorm.WhereSelect(orm.OpEq, 42).
			Or(aorm.WhereSelectBetween(10, 50)).
			Or(aorm.WhereSelectIn(11, 12))).
		OrderBySelect(orm.Desc).
		GroupBySelect()

	alls, err := query.Query()
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))
	assert.Equal(t, 42, alls[0].Select)

	res, err = db.Update().SetSelect(11).Where(aorm.WhereSelect(orm.OpEq, 42)).Exec()
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
}

func prepareAll(t *testing.T) aorm.API {
	t.Helper()
	db, err := aorm.Open(":memory:")
	require.Nil(t, err)
	if testing.Verbose() {
		db.Logger(t.Logf)
	}
	_, err = db.Create().Exec()
	require.Nil(t, err)

	return db
}
