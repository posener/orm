package example_test

import (
	"testing"

	"github.com/posener/orm/example/allorm"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	assert.Equal(t,
		`CREATE TABLE all ( int INT, text VARCHAR(100) NOT NULL, bool BOOLEAN, PRIMARY KEY (int) )`,
		allorm.Create().String(),
	)
}
