package sqlite3

import (
	"testing"

	"github.com/posener/orm/common"
	"github.com/stretchr/testify/assert"
)

func TestSqlite3_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fields []common.Field
		want   string
	}{
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int"}},
				{Type: "string", SQL: common.SQL{Column: "string"}},
				{Type: "bool", SQL: common.SQL{Column: "bool"}},
				{Type: "time.Time", SQL: common.SQL{Column: "time"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER, 'string' TEXT, 'bool' BOOLEAN, 'time' TIMESTAMP )",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int", PrimaryKey: true}},
				{Type: "string", SQL: common.SQL{Column: "string"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER PRIMARY KEY, 'string' TEXT )",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int", PrimaryKey: true, AutoIncrement: true}},
				{Type: "string", SQL: common.SQL{Column: "string"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER PRIMARY KEY AUTOINCREMENT, 'string' TEXT )",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int"}},
				{Type: "string", SQL: common.SQL{Column: "string", NotNull: true, Default: "xxx"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER, 'string' TEXT NOT NULL DEFAULT xxx )",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int"}},
				{Type: "string", SQL: common.SQL{Column: "string", CustomType: "VARCHAR(10)"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER, 'string' VARCHAR(10) )",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int"}},
				{Type: "time.Time", SQL: common.SQL{Column: "time", CustomType: "DATETIME"}},
			},
			want: "CREATE TABLE 'name' ( 'int' INTEGER, 'time' DATETIME )",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			s := sqlite3{tp: common.Type{Name: "name", Fields: tt.fields}}
			got := s.Create()
			assert.Equal(t, tt.want, got)
		})
	}

}
