package main_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/posener/orm/example"
	aorm "github.com/posener/orm/example/allorm"
	porm "github.com/posener/orm/example/personorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BenchmarkORMInsert tests inserts with posener/orm package
func BenchmarkORMInsert(b *testing.B) {
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

// BenchmarkGORMInsert tests inserts with jinzhu/gorm package
func BenchmarkGORMInsert(b *testing.B) {
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

// BenchmarkRawInsert tests inserts with raw SQL commands
func BenchmarkRawInsert(b *testing.B) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(b, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE 'person' ( 'name' VARCHAR (255), 'age' INTEGER )`)
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = db.Exec(`INSERT INTO 'person' ('name', 'age') VALUES (?, ?)`, "xxx", i)
		assert.Nil(b, err)
	}
}

// datasetSize is the number of columns that are used in querying benchmarks
const datasetSize = 1000

// BenchmarkORMQuery tests queries with posener/orm package
func BenchmarkORMQuery(b *testing.B) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(b, err)
	defer db.Close()

	orm := porm.New(db)

	_, err = orm.Create().Exec()
	require.Nil(b, err)

	for i := 0; i < datasetSize; i++ {
		_, err = orm.InsertPerson(&example.Person{Name: "xxx", Age: i}).Exec()
		require.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps, err := orm.Select().Query()
		require.Nil(b, err)
		assert.Equal(b, datasetSize, len(ps))
	}
}

// BenchmarkGORMQuery tests queries with jinzhu/gorm package
func BenchmarkGORMQuery(b *testing.B) {
	db, err := gorm.Open("sqlite3", ":memory:")
	require.Nil(b, err)
	defer db.Close()

	err = db.AutoMigrate(&example.Person{}).Error
	require.Nil(b, err)

	for i := 0; i < datasetSize; i++ {
		err = db.Create(&example.Person{Name: "xxx", Age: i}).Error
		require.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var ps []example.Person
		err = db.Find(&ps).Error
		require.Nil(b, err)
		assert.Equal(b, datasetSize, len(ps))
	}
}

// BenchmarkRawQuery tests queries with raw SQL commands
func BenchmarkRawQuery(b *testing.B) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(b, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE 'person' ( 'name' VARCHAR (255), 'age' INTEGER )`)
	require.Nil(b, err)

	for i := 0; i < datasetSize; i++ {
		_, err = db.Exec(`INSERT INTO 'person' ('name', 'age') VALUES (?, ?)`, "xxx", i)
		require.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := db.Query(`SELECT * FROM 'person'`)
		require.Nil(b, err)
		var ps []example.Person
		for rows.Next() {
			var p example.Person
			err := rows.Scan(&p.Name, &p.Age)
			require.Nil(b, err)
			ps = append(ps, p)
		}
		assert.Equal(b, datasetSize, len(ps))
		rows.Close()
	}
}

// BenchmarkORMQueryLargeStruct tests large struct queries with posener/orm package
func BenchmarkORMQueryLargeStruct(b *testing.B) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(b, err)
	defer db.Close()

	orm := aorm.New(db)

	_, err = orm.Create().Exec()
	require.Nil(b, err)

	tm := time.Now().Round(time.Millisecond).UTC()

	for i := 0; i < datasetSize; i++ {
		_, err = orm.InsertAll(&example.All{String: "xxx", Select: i, Auto: i, Int: i, Time: tm, Bool: true}).Exec()
		require.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alls, err := orm.Select().Query()
		require.Nil(b, err)
		assert.Equal(b, datasetSize, len(alls))
	}
}

// BenchmarkGORMQueryLargeStruct tests large struct  queries with jinzhu/gorm package
func BenchmarkGORMQueryLargeStruct(b *testing.B) {
	db, err := gorm.Open("sqlite3", ":memory:")
	require.Nil(b, err)
	defer db.Close()

	err = db.AutoMigrate(&example.All{}).Error
	require.Nil(b, err)

	tm := time.Now().Round(time.Millisecond).UTC()

	for i := 0; i < datasetSize; i++ {
		err = db.Create(&example.All{String: "xxx", Select: i, Auto: i, Int: i, Time: tm, Bool: true}).Error
		require.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var alls []example.All
		err = db.Find(&alls).Error
		require.Nil(b, err)
		assert.Equal(b, datasetSize, len(alls))
	}
}

// BenchmarkRawQueryLargeStruct tests large struct  queries with raw SQL commands
func BenchmarkRawQueryLargeStruct(b *testing.B) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(b, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE 'all' ( 'int' INT PRIMARY KEY, 'string' VARCHAR(100) NOT NULL, 'bool' BOOLEAN, 'auto' INT AUTO_INCREMENT, 'time' TIMESTAMP, 'select' INT )`)
	require.Nil(b, err)

	tm := time.Now().Round(time.Millisecond).UTC()

	for i := 0; i < datasetSize; i++ {
		_, err = db.Exec(`INSERT INTO 'all' ('int', 'string', 'bool', 'auto', 'time', 'select') VALUES (?, ?, ?, ?, ?, ?)`, i, "xxx", true, i, tm, i)
		require.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := db.Query(`SELECT * FROM 'all'`)
		require.Nil(b, err)
		var ps []example.All
		for rows.Next() {
			var p example.All
			err := rows.Scan(&p.Int, &p.String, &p.Bool, &p.Auto, &p.Time, &p.Select)
			require.Nil(b, err)
			ps = append(ps, p)
		}
		assert.Equal(b, datasetSize, len(ps))
		rows.Close()
	}
}
