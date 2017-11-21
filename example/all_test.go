package example_test

import (
	"testing"

	"database/sql"

	"time"

	"github.com/posener/orm/example"
	aorm "github.com/posener/orm/example/allorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	assert.Equal(t,
		`CREATE TABLE 'all' ( 'int' BIGINT PRIMARY KEY, 'string' VARCHAR(100) NOT NULL, 'bool' BOOLEAN, 'time' TIMESTAMP )`,
		new(aorm.ORM).Create().String(),
	)
}

func TestTime(t *testing.T) {
	t.Parallel()
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err)
	defer db.Close()

	orm := aorm.New(db)

	tm := time.Now().Round(time.Millisecond).UTC()

	_, err = orm.Create().Exec()
	require.Nil(t, err)
	res, err := orm.InsertAll(&example.All{Time: tm}).Exec()
	require.Nil(t, err, "Failed inserting")
	affected, err := res.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, int64(1), affected)

	alls, err := orm.Select().Query()
	require.Nil(t, err)
	require.Equal(t, 1, len(alls))

	assert.Equal(t, tm, alls[0].Time)
}
