package example_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/posener/orm/example"
	aorm "github.com/posener/orm/example/allorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	assert.Equal(t,
		`CREATE TABLE 'all' ( 'int' INT PRIMARY KEY, 'string' VARCHAR(100) NOT NULL, 'bool' BOOLEAN, 'auto' INT AUTO_INCREMENT, 'time' TIMESTAMP, 'select' INT )`,
		new(aorm.ORM).Create().String(),
	)
}

func TestTime(t *testing.T) {
	t.Parallel()
	orm := prepare(t)

	tm := time.Now().Round(time.Millisecond).UTC()

	res, err := orm.Insert().SetTime(tm).SetString("not null").Exec()
	require.Nil(t, err)
	affected, err := res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	alls, err := orm.Select().Query()
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))

	assert.Equal(t, tm, alls[0].Time)
}

func TestAutoIncrement(t *testing.T) {
	t.Parallel()
	orm := prepare(t)

	res, err := orm.Insert().SetString("1").Exec()
	require.Nil(t, err)
	affected, err := res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	res, err = orm.Insert().SetString("2").Exec()
	require.Nil(t, err)
	affected, err = res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	alls, err := orm.Select().OrderByAuto(aorm.Asc).Query()
	require.Nil(t, err)
	require.Equal(t, 2, len(alls))

	assert.Equal(t, 1, alls[0].Auto)
	assert.Equal(t, 2, alls[1].Auto)
}

// TestNotNull tests that given inserting an empty not null field causes an error
func TestNotNull(t *testing.T) {
	t.Parallel()
	orm := prepare(t)

	_, err := orm.Insert().SetInt(1).Exec()
	require.NotNil(t, err)
}

func TestFieldReservedName(t *testing.T) {
	t.Parallel()
	orm := prepare(t)

	res, err := orm.InsertAll(&example.All{Select: 42}).Exec()
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)

	query := orm.Select().
		SelectSelect().
		Where(aorm.WhereSelect(aorm.OpEq, 42).
			Or(aorm.WhereSelectBetween(10, 50)).
			Or(aorm.WhereSelectIn(11, 12))).
		OrderBySelect(aorm.Desc).
		GroupBySelect()

	alls, err := query.Query()
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))
	assert.Equal(t, 42, alls[0].Select)

	res, err = orm.Update().SetSelect(11).Where(aorm.WhereSelect(aorm.OpEq, 42)).Exec()
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)

	alls, err = orm.Select().SelectSelect().Query()
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))
	assert.Equal(t, 11, alls[0].Select)

	res, err = orm.Delete().Exec()
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)

	alls, err = orm.Select().SelectSelect().Query()
	require.Nil(t, err)
	require.Equal(t, 0, len(alls))
}

func prepare(t *testing.T) aorm.API {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err)
	orm := aorm.New(db)
	if testing.Verbose() {
		orm.Logger(t.Logf)
	}
	_, err = orm.Create().Exec()
	require.Nil(t, err)

	return orm
}
