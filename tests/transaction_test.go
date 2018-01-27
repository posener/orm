package tests

import (
	"context"
	"testing"

	"github.com/posener/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactions(t *testing.T) {
	t.Parallel()
	testDBs(t, func(t *testing.T, conn orm.Conn) {
		person, err := NewPersonORM(conn)
		require.Nil(t, err)

		require.Nil(t, person.Drop().IfExists().Exec())
		require.Nil(t, person.Create().Exec())

		ps, err := person.Select().Query()
		assert.Equal(t, 0, len(ps))

		ctx := context.Background()

		t.Logf("Transaction commit")

		personTx, err := person.Begin(ctx, nil)
		require.Nil(t, err)

		_, err = personTx.Insert().SetName("John").SetAge(23).Exec()
		require.Nil(t, err)

		ps, err = personTx.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(ps))

		ps, err = person.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, 0, len(ps))

		err = personTx.Commit()
		require.Nil(t, err)

		ps, err = person.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(ps))

		t.Logf("Transaction rollback")

		personTx, err = person.Begin(ctx, nil)
		require.Nil(t, err)

		_, err = personTx.Insert().SetName("Bill").SetAge(46).Exec()

		ps, err = person.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(ps))

		require.Nil(t, personTx.Rollback())

		ps, err = person.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(ps))

		t.Logf("Transaction context cancel")
		ctx, cancel := context.WithCancel(ctx)

		personTx, err = person.Begin(ctx, nil)
		require.Nil(t, err)

		_, err = personTx.Insert().SetName("Bill").SetAge(46).Exec()

		ps, err = person.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(ps))

		cancel()
		require.NotNil(t, personTx.Commit())

		ps, err = person.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(ps))
	})
}
