package external

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExternal tests an ORM generated with -out flag
func TestExternal(t *testing.T) {
	db, err := orm.Open("sqlite3", ":memory:")
	require.Nil(t, err)
	person, err := NewPersonORM(db)
	require.Nil(t, err)
	assert.Nil(t, person.Create().IfNotExists().Exec())
}
