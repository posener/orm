package example_test

import (
	"database/sql"
	"testing"

	"github.com/posener/orm/example"
	porm "github.com/posener/orm/example/personorm"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BenchmarkORM tests posener/orm package
func BenchmarkORM(b *testing.B) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(b, err)
	defer db.Close()

	orm := porm.New(db)

	_, err = orm.Create().Exec()
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = orm.InsertPerson(&example.Person{Name: "xxx", Age: i}).Exec()
		assert.Nil(b, err)
	}
}

// BenchmarkORM tests jinzhu/gorm package
func BenchmarkGORM(b *testing.B) {
	db, err := gorm.Open("sqlite3", ":memory:")
	require.Nil(b, err)
	defer db.Close()

	err = db.AutoMigrate(&example.Person{}).Error
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = db.Create(&example.Person{Name: "xxx", Age: i}).Error
		assert.Nil(b, err)
	}
}
