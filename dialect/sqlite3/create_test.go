package sqlite3

import (
	"testing"

	"github.com/posener/orm/def"
	"github.com/stretchr/testify/assert"
)

func TestSqlite3_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fields []def.Field
		want   string
	}{
		{
			fields: []def.Field{
				{Type: "int", SQL: def.SQL{Column: "int"}},
				{Type: "string", SQL: def.SQL{Column: "string"}},
				{Type: "bool", SQL: def.SQL{Column: "bool"}},
				{Type: "time.Time", SQL: def.SQL{Column: "time"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER, 'string' TEXT, 'bool' BOOLEAN, 'time' TIMESTAMP )",
		},
		{
			fields: []def.Field{
				{Type: "int", SQL: def.SQL{Column: "int", PrimaryKey: true}},
				{Type: "string", SQL: def.SQL{Column: "string"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER PRIMARY KEY, 'string' TEXT )",
		},
		{
			fields: []def.Field{
				{Type: "int", SQL: def.SQL{Column: "int", PrimaryKey: true, AutoIncrement: true}},
				{Type: "string", SQL: def.SQL{Column: "string"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER PRIMARY KEY AUTOINCREMENT, 'string' TEXT )",
		},
		{
			fields: []def.Field{
				{Type: "int", SQL: def.SQL{Column: "int"}},
				{Type: "string", SQL: def.SQL{Column: "string", NotNull: true, Default: "xxx"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER, 'string' TEXT NOT NULL DEFAULT xxx )",
		},
		{
			fields: []def.Field{
				{Type: "int", SQL: def.SQL{Column: "int"}},
				{Type: "string", SQL: def.SQL{Column: "string", CustomType: "VARCHAR(10)"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER, 'string' VARCHAR(10) )",
		},
		{
			fields: []def.Field{
				{Type: "int", SQL: def.SQL{Column: "int"}},
				{Type: "time.Time", SQL: def.SQL{Column: "time", CustomType: "DATETIME"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER, 'time' DATETIME )",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			s := sqlite3{tp: def.Type{Name: "name", Fields: tt.fields}}
			got := s.Create()
			assert.Equal(t, tt.want, got)
		})
	}

}
