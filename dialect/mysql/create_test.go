package mysql

import (
	"testing"

	"github.com/posener/orm/common"
	"github.com/stretchr/testify/assert"
)

func TestSqlite3_ColumnsStatement(t *testing.T) {
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
			want: "`int` INTEGER, `string` TEXT, `bool` BOOLEAN, `time` DATETIME(3)",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int", PrimaryKey: true}},
				{Type: "string", SQL: common.SQL{Column: "string"}},
			},
			want: "`int` INTEGER PRIMARY KEY, `string` TEXT",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int", PrimaryKey: true, AutoIncrement: true}},
				{Type: "string", SQL: common.SQL{Column: "string"}},
			},
			want: "`int` INTEGER PRIMARY KEY AUTO_INCREMENT, `string` TEXT",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int"}},
				{Type: "string", SQL: common.SQL{Column: "string", NotNull: true, Default: "xxx"}},
			},
			want: "`int` INTEGER, `string` TEXT NOT NULL DEFAULT xxx",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int"}},
				{Type: "string", SQL: common.SQL{Column: "string", CustomType: "VARCHAR(10)"}},
			},
			want: "`int` INTEGER, `string` VARCHAR(10)",
		},
		{
			fields: []common.Field{
				{Type: "int", SQL: common.SQL{Column: "int"}},
				{Type: "time.Time", SQL: common.SQL{Column: "time", CustomType: "DATETIME"}},
			},
			want: "`int` INTEGER, `time` DATETIME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			s := Gen{Tp: common.Type{Name: "name", Fields: tt.fields}}
			got := s.ColumnsStatement()
			assert.Equal(t, tt.want, got)
		})
	}

}
