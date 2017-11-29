package sqlite3

import (
	"testing"

	"github.com/posener/orm/load"
	"github.com/stretchr/testify/assert"
)

func TestSqlite3_ColumnsStatement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fields []load.Field
		want   string
	}{
		{
			fields: []load.Field{
				{GoType: load.GoType{Type: "int"}, SQL: load.SQL{Column: "int"}},
				{GoType: load.GoType{Type: "string"}, SQL: load.SQL{Column: "string"}},
				{GoType: load.GoType{Type: "bool"}, SQL: load.SQL{Column: "bool"}},
				{GoType: load.GoType{Type: "time.Time"}, SQL: load.SQL{Column: "time"}},
			},
			want: "'int' INTEGER, 'string' TEXT, 'bool' BOOLEAN, 'time' TIMESTAMP",
		},
		{
			fields: []load.Field{
				{GoType: load.GoType{Type: "int"}, SQL: load.SQL{Column: "int", PrimaryKey: true}},
				{GoType: load.GoType{Type: "string"}, SQL: load.SQL{Column: "string"}},
			},
			want: "'int' INTEGER PRIMARY KEY, 'string' TEXT",
		},
		{
			fields: []load.Field{
				{GoType: load.GoType{Type: "int"}, SQL: load.SQL{Column: "int", PrimaryKey: true, AutoIncrement: true}},
				{GoType: load.GoType{Type: "string"}, SQL: load.SQL{Column: "string"}},
			},
			want: "'int' INTEGER PRIMARY KEY AUTOINCREMENT, 'string' TEXT",
		},
		{
			fields: []load.Field{
				{GoType: load.GoType{Type: "int"}, SQL: load.SQL{Column: "int"}},
				{GoType: load.GoType{Type: "string"}, SQL: load.SQL{Column: "string", NotNull: true, Default: "xxx"}},
			},
			want: "'int' INTEGER, 'string' TEXT NOT NULL DEFAULT xxx",
		},
		{
			fields: []load.Field{
				{GoType: load.GoType{Type: "int"}, SQL: load.SQL{Column: "int"}},
				{GoType: load.GoType{Type: "string"}, SQL: load.SQL{Column: "string", CustomType: "VARCHAR(10)"}},
			},
			want: "'int' INTEGER, 'string' VARCHAR(10)",
		},
		{
			fields: []load.Field{
				{GoType: load.GoType{Type: "int"}, SQL: load.SQL{Column: "int"}},
				{GoType: load.GoType{Type: "time.Time"}, SQL: load.SQL{Column: "time", CustomType: "DATETIME"}},
			},
			want: "'int' INTEGER, 'time' DATETIME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			s := Gen{Tp: &load.Type{GoType: load.GoType{Type: "name"}, Fields: tt.fields}}
			got := s.ColumnsStatement()
			assert.Equal(t, tt.want, got)
		})
	}

}
