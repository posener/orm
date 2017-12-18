package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		if conn.name == "sqlite3" {
			t.Skip("sqlite migrations is not supported")
		}
		m0, err := NewMigration0ORM(conn.name, conn)
		require.Nil(t, err)
		m1, err := NewMigration1ORM(conn.name, conn)
		require.Nil(t, err)
		m2, err := NewMigration2ORM(conn.name, conn)
		require.Nil(t, err)
		m3, err := NewMigration3ORM(conn.name, conn)
		require.Nil(t, err)
		p, err := NewC2ORM(conn.name, conn) // for foreign key constraint
		require.Nil(t, err)
		if testing.Verbose() {
			m0.Logger(t.Logf)
			m1.Logger(t.Logf)
			m2.Logger(t.Logf)
			m3.Logger(t.Logf)
			p.Logger(t.Logf)
		}
		_, err = conn.ExecContext(context.Background(), "DROP TABLE IF EXISTS migrations")
		require.Nil(t, err)

		t.Logf("migration to 0")

		require.Nil(t, m0.Create().AutoMigrate().Exec())

		_, err = m0.Insert().InsertMigration0(&Migration0{A: "A0"}).Exec()
		require.Nil(t, err)

		t.Logf("migration to 1")

		require.Nil(t, m1.Create().AutoMigrate().Exec())

		_, err = m1.Insert().InsertMigration1(&Migration1{A: "A1", B: "B1"}).Exec()
		require.Nil(t, err)

		ms, err := m1.Select().Query()
		require.Nil(t, err)

		if assert.Equal(t, 2, len(ms)) {
			assert.Equal(t, "A0", ms[0].A)
			assert.Equal(t, "", ms[0].B)
			assert.Equal(t, "A1", ms[1].A)
			assert.Equal(t, "B1", ms[1].B)
		}

		require.Nil(t, p.Create().IfNotExists().Exec()) // for foreign key constraint
		c1, err := p.Insert().SetName("C1").Exec()
		require.Nil(t, err)

		t.Logf("migration to 2")

		require.Nil(t, m2.Create().AutoMigrate().Exec())

		_, err = m2.Insert().SetA("A2").SetD("x").SetP1(c1).Exec()
		require.Nil(t, err)

		t.Logf("migration to 3")

		require.Nil(t, m3.Create().AutoMigrate().Exec())
		_, err = m3.Insert().SetA("A3").SetD("x").SetP1(c1).SetP2(c1).Exec()
		require.Nil(t, err)
	})
}
