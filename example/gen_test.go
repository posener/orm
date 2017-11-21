package example_test

import (
	"testing"

	"github.com/posener/orm/example/allorm"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	assert.Equal(t,
		`CREATE TABLE all ( int BIGINT PRIMARY KEY, string VARCHAR(100) NOT NULL, bool BOOLEAN, time TIMESTAMP )`,
		new(allorm.ORM).Create().String(),
	)
}
