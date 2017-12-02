package mysql

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
				{Name: "Int", Type: load.Type{Name: "int"}},
				{Name: "String", Type: load.Type{Name: "string"}},
				{Name: "Bool", Type: load.Type{Name: "bool"}},
				{Name: "Time", Type: load.Type{Name: "time.Time"}},
			},
			want: "`int` INTEGER, `string` TEXT, `bool` BOOLEAN, `time` DATETIME(3)",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}, SQL: load.SQL{PrimaryKey: true}},
				{Name: "String", Type: load.Type{Name: "string"}},
			},
			want: "`int` INTEGER PRIMARY KEY, `string` TEXT",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}, SQL: load.SQL{PrimaryKey: true, AutoIncrement: true}},
				{Name: "String", Type: load.Type{Name: "string"}},
			},
			want: "`int` INTEGER PRIMARY KEY AUTO_INCREMENT, `string` TEXT",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}},
				{Name: "String", Type: load.Type{Name: "string"}, SQL: load.SQL{NotNull: true, Default: "xxx"}},
			},
			want: "`int` INTEGER, `string` TEXT NOT NULL DEFAULT xxx",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}},
				{Name: "String", Type: load.Type{Name: "string"}, SQL: load.SQL{CustomType: "VARCHAR(10)"}},
			},
			want: "`int` INTEGER, `string` VARCHAR(10)",
		},
		{
			fields: []load.Field{
				{Name: "Int", Type: load.Type{Name: "int"}},
				{Name: "Time", Type: load.Type{Name: "time.Time"}, SQL: load.SQL{CustomType: "DATETIME"}},
			},
			want: "`int` INTEGER, `time` DATETIME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			s := Gen{Tp: &load.Type{Name: "name", Fields: tt.fields}}
			got := s.ColumnsStatement()
			assert.Equal(t, tt.want, got)
		})
	}

}
