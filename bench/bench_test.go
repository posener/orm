package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/posener/orm/example"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BenchmarkORMInsert tests inserts with posener/orm package
func BenchmarkORMInsert(b *testing.B) {
	orm, err := example.OpenPersonORM("sqlite3", ":memory:")
	require.Nil(b, err)
	defer orm.Close()

	require.Nil(b, orm.Create().Exec())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = orm.Insert().InsertPerson(&example.Person{Name: "xxx", Age: i}).Exec()
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
	orm, err := example.OpenPersonORM("sqlite3", ":memory:")
	require.Nil(b, err)
	defer orm.Close()

	require.Nil(b, orm.Create().Exec())

	for i := 0; i < datasetSize; i++ {
		_, err = orm.Insert().InsertPerson(&example.Person{Name: "xxx", Age: i}).Exec()
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
	orm, err := example.OpenAllORM("sqlite3", ":memory:")
	require.Nil(b, err)
	defer orm.Close()

	require.Nil(b, orm.Create().Exec())

	tm := time.Now().Round(time.Millisecond).UTC()

	for i := 0; i < datasetSize; i++ {
		_, err = orm.Insert().InsertAll(&example.All{String: "xxx", Select: i, Int: i, Time: tm, Bool: true, NotNil: "notnil"}).Exec()
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
		err = db.Create(&example.All{String: "xxx", Select: i, Int: i, Time: tm, Bool: true, NotNil: "notnil"}).Error
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

	_, err = db.Exec("CREATE TABLE  'all' ( 'auto' INTEGER PRIMARY KEY AUTOINCREMENT, 'notnil' TEXT NOT NULL, 'int' INTEGER, 'int8' INTEGER, 'int16' INTEGER, 'int32' INTEGER, 'int64' INTEGER, 'uint' INTEGER, 'uint8' INTEGER, 'uint16' INTEGER, 'uint32' INTEGER, 'uint64' INTEGER, 'time' TIMESTAMP, 'varcharstring' VARCHAR(100), 'varcharbyte' VARCHAR(100), 'string' TEXT, 'bytes' BLOB, 'bool' BOOLEAN, 'pint' INTEGER, 'pint8' INTEGER, 'pint16' INTEGER, 'pint32' INTEGER, 'pint64' INTEGER, 'puint' INTEGER, 'puint8' INTEGER, 'puint16' INTEGER, 'puint32' INTEGER, 'puint64' INTEGER, 'ptime' TIMESTAMP, 'pvarcharstring' VARCHAR(100), 'pvarcharbyte' VARCHAR(100), 'pstring' TEXT, 'pbytes' BLOB, 'pbool' BOOLEAN, 'select' INTEGER )")
	require.Nil(b, err)

	tm := time.Now().Round(time.Millisecond).UTC()

	for i := 0; i < datasetSize; i++ {
		_, err = db.Exec(`INSERT INTO 'all' ('int', 'string', 'bool', 'time', 'select', 'notnil') VALUES (?, ?, ?, ?, ?, ?)`, i, "xxx", true, tm, i, "notnil")
		require.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := db.Query(`SELECT * FROM 'all'`)
		require.Nil(b, err)
		var ps []example.All
		for rows.Next() {
			var p1, p2 example.All
			err := rows.Scan(
				&p1.Auto, &p1.NotNil, &p2.PInt, &p2.PInt8, &p2.PInt16,
				&p2.PInt32, &p2.PInt64, &p2.PUInt, &p2.PUInt8, &p2.PUInt16,
				&p2.PUInt32, &p2.PUInt64, &p2.Time, &p2.PVarCharString, &p1.VarCharByte,
				&p1.String, &p1.Bytes, &p1.Bool,
				&p1.PInt, &p1.PInt8, &p1.PInt16, &p1.PInt32, &p1.PInt64,
				&p1.PUInt, &p1.PUInt8, &p1.PUInt16, &p1.PUInt32, &p1.PUInt64,
				&p1.PTime, &p1.PVarCharString, &p1.PVarCharByte,
				&p1.PString, &p1.PBytes, &p1.PBool, &p1.Select,
			)
			require.Nil(b, err)
			ps = append(ps, p1)
		}
		assert.Equal(b, datasetSize, len(ps))
		rows.Close()
	}
}
