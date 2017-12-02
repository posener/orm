package dialect

import (
	"testing"

	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/sqlite3"
	"github.com/posener/orm/load"
	"github.com/stretchr/testify/assert"
)

func TestSqlite3_ColumnsStatement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fields      []load.Field
		sqlite3Want string
		mysqlWant   string
	}{
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}},
				{Name: "String", Type: load.Type{Name: "string"}},
				{Name: "Bool", Type: load.Type{Name: "bool"}},
				{Name: "Time", Type: load.Type{Name: "time.Time"}},
			},
			sqlite3Want: "'int' INTEGER, 'string' TEXT, 'bool' BOOLEAN, 'time' TIMESTAMP",
			mysqlWant:   "`int` INTEGER, `string` TEXT, `bool` BOOLEAN, `time` DATETIME(3)",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}, SQL: load.SQL{PrimaryKey: true}},
				{Name: "String", Type: load.Type{Name: "string"}},
			},
			sqlite3Want: "'int' INTEGER PRIMARY KEY, 'string' TEXT",
			mysqlWant:   "`int` INTEGER PRIMARY KEY, `string` TEXT",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}, SQL: load.SQL{PrimaryKey: true, AutoIncrement: true}},
				{Name: "String", Type: load.Type{Name: "string"}},
			},
			sqlite3Want: "'int' INTEGER PRIMARY KEY AUTOINCREMENT, 'string' TEXT",
			mysqlWant:   "`int` INTEGER PRIMARY KEY AUTO_INCREMENT, `string` TEXT",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}},
				{Name: "String", Type: load.Type{Name: "string"}, SQL: load.SQL{NotNull: true, Default: "xxx"}},
			},
			sqlite3Want: "'int' INTEGER, 'string' TEXT NOT NULL DEFAULT xxx",
			mysqlWant:   "`int` INTEGER, `string` TEXT NOT NULL DEFAULT xxx",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}},
				{Name: "String", Type: load.Type{Name: "string"}, SQL: load.SQL{CustomType: "VARCHAR(10)"}},
			},
			sqlite3Want: "'int' INTEGER, 'string' VARCHAR(10)",
			mysqlWant:   "`int` INTEGER, `string` VARCHAR(10)",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}},
				{Name: "Time", Type: load.Type{Name: "time.Time"}, SQL: load.SQL{CustomType: "DATETIME"}},
			},
			sqlite3Want: "'int' INTEGER, 'time' DATETIME",
			mysqlWant:   "`int` INTEGER, `time` DATETIME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.sqlite3Want, func(t *testing.T) {
			tp := &load.Type{Name: "name", Fields: tt.fields}
			got := new(sqlite3.Gen).ColumnsStatement(tp)
			assert.Equal(t, tt.sqlite3Want, got)
			got = new(mysql.Gen).ColumnsStatement(tp)
			assert.Equal(t, tt.mysqlWant, got)
		})
	}
}
