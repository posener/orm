package tests

import (
	"context"
	"testing"

	"github.com/posener/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactions(t *testing.T) {
	testDBs(t, func(t *testing.T, conn orm.DB) {
		person, err := NewPersonORM(conn)
		require.Nil(t, err)

		require.Nil(t, person.Create().Exec())

		ctx := context.Background()

		t.Logf("Transaction commit")

		personTx, err := person.Begin(ctx, nil)
		require.Nil(t, err)

		_, err = personTx.Insert().SetName("John").SetAge(23).Exec()

		ps, err := person.Select().Query()
		assert.Equal(t, 0, len(ps))

		require.Nil(t, personTx.Commit())

		ps, err = person.Select().Query()
		assert.Equal(t, 1, len(ps))

		t.Logf("Transaction rollback")

		personTx, err = person.Begin(ctx, nil)
		require.Nil(t, err)

		_, err = personTx.Insert().SetName("Bill").SetAge(46).Exec()

		ps, err = person.Select().Query()
		assert.Equal(t, 1, len(ps))

		require.Nil(t, personTx.Rollback())

		ps, err = person.Select().Query()
		assert.Equal(t, 1, len(ps))

		t.Logf("Transaction context cancel")

		ctx, cancel := context.WithCancel(ctx)

		personTx, err = person.Begin(ctx, nil)
		require.Nil(t, err)

		_, err = personTx.Insert().SetName("Bill").SetAge(46).Exec()

		ps, err = person.Select().Query()
		assert.Equal(t, 1, len(ps))

		cancel()
		require.NotNil(t, personTx.Commit())

		ps, err = person.Select().Query()
		assert.Equal(t, 1, len(ps))

	})
}
