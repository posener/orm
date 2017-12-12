package dialect

import (
	"regexp"
	"strings"
	"testing"

	"github.com/posener/orm"
	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/load"
	"github.com/stretchr/testify/assert"
)

const table = "name"

func TestCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fields []*load.Field
		want   string
	}{
		{
			fields: []*load.Field{
				{Name: "Int", Type: load.Type{Naked: &load.Naked{Name: "int"}}},
				{Name: "String", Type: load.Type{Naked: &load.Naked{Name: "string"}}},
				{Name: "Bool", Type: load.Type{Naked: &load.Naked{Name: "bool"}}},
				{Name: "Time", Type: load.Type{Naked: &load.Naked{Name: "time.Time"}}},
			},
			want: "`int` INTEGER, `string` TEXT, `bool` BOOLEAN, `time` DATETIME(3)",
		},
		{
			fields: []*load.Field{
				{Name: "Int", Type: load.Type{Naked: &load.Naked{Name: "int"}}, PrimaryKey: true},
				{Name: "String", Type: load.Type{Naked: &load.Naked{Name: "string"}}},
			},
			want: "`int` INTEGER PRIMARY KEY, `string` TEXT",
		},
		{
			fields: []*load.Field{
				{Name: "Int", Type: load.Type{Naked: &load.Naked{Name: "int"}}, PrimaryKey: true, AutoIncrement: true},
				{Name: "String", Type: load.Type{Naked: &load.Naked{Name: "string"}}},
			},
			want: "`int` INTEGER PRIMARY KEY AUTO_INCREMENT, `string` TEXT",
		},
		{
			fields: []*load.Field{
				{Name: "Int", Type: load.Type{Naked: &load.Naked{Name: "int"}}},
				{Name: "String", Type: load.Type{Naked: &load.Naked{Name: "string"}}, NotNull: true, Default: "xxx"},
			},
			want: "`int` INTEGER, `string` TEXT NOT NULL DEFAULT xxx",
		},
		{
			fields: []*load.Field{
				{Name: "Int", Type: load.Type{Naked: &load.Naked{Name: "int"}}},
				{Name: "String", Type: load.Type{Naked: &load.Naked{Name: "string"}}, CustomType: "VARCHAR(10)"},
			},
			want: "`int` INTEGER, `string` VARCHAR(10)",
		},
		{
			fields: []*load.Field{
				{Name: "Int", Type: load.Type{Naked: &load.Naked{Name: "int"}}},
				{Name: "Time", Type: load.Type{Naked: &load.Naked{Name: "time.Time"}}, CustomType: "DATETIME"},
			},
			want: "`int` INTEGER, `time` DATETIME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			tp := &load.Type{Naked: &load.Naked{Name: "name", Fields: tt.fields}}
			genMysql := &gen{GenImplementer: new(mysql.Gen)}
			got := genMysql.ColumnsStatement(tp)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		sel      common.SelectParams
		wantStmt string
		wantArgs []interface{}
	}{
		{
			sel:      common.SelectParams{Columns: &columner{}},
			wantStmt: "SELECT `name`.* FROM `name`",
		},
		{
			sel:      common.SelectParams{Columns: &columner{}, Page: common.Page{}},
			wantStmt: "SELECT `name`.* FROM `name`",
		},
		{
			sel:      common.SelectParams{Columns: &columner{count: true}},
			wantStmt: "SELECT COUNT(*) FROM `name`",
		},
		{
			sel:      common.SelectParams{Columns: &columner{columns: []string{"a", "b", "c"}, count: true}},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM `name`",
		},
		{
			sel:      common.SelectParams{Columns: &columner{columns: []string{"a", "b", "c"}}},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c` FROM `name`",
		},
		{
			sel:      common.SelectParams{Columns: &columner{}, Page: common.Page{}},
			wantStmt: "SELECT `name`.* FROM `name`",
		},
		{
			sel:      common.SelectParams{Columns: &columner{}, Page: common.Page{Limit: 1}},
			wantStmt: "SELECT `name`.* FROM `name` LIMIT 1",
		},
		{
			sel:      common.SelectParams{Columns: &columner{}, Page: common.Page{Limit: 1, Offset: 2}},
			wantStmt: "SELECT `name`.* FROM `name` LIMIT 1 OFFSET 2",
		},
		{
			sel:      common.SelectParams{Columns: &columner{}, Page: common.Page{Offset: 1}},
			wantStmt: "SELECT `name`.* FROM `name`",
		},
		{
			sel: common.SelectParams{
				Columns: &columner{columns: []string{"a", "b", "c"}, count: true},
				Page:    common.Page{Limit: 1, Offset: 2},
			},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM `name` LIMIT 1 OFFSET 2",
		},
		{
			sel: common.SelectParams{
				Columns: &columner{},
				Groups:  common.Groups{{Column: "a"}, {Column: "b"}},
			},
			wantStmt: "SELECT `name`.* FROM `name` GROUP BY `name`.`a`, `name`.`b`",
		},
		{
			sel: common.SelectParams{
				Columns: &columner{},
				Orders: common.Orders{
					{Column: "c", Dir: "ASC"},
					{Column: "d", Dir: "DESC"},
				},
			},
			wantStmt: "SELECT `name`.* FROM `name` ORDER BY `name`.`c` ASC, `name`.`d` DESC",
		},
		{
			sel: common.SelectParams{
				Columns: &columner{},
				Where:   common.NewWhere(orm.OpEq, "k", 3),
			},
			wantStmt: "SELECT `name`.* FROM `name` WHERE `name`.`k` = ?",
			wantArgs: []interface{}{3},
		},
		{
			sel: common.SelectParams{
				Columns: &columner{columns: []string{"a", "b", "c"}, count: true},
				Where:   common.NewWhere(orm.OpGt, "k", 3),
				Groups:  common.Groups{{Column: "a"}, {Column: "b"}},
				Orders: common.Orders{
					{Column: "c", Dir: "ASC"},
					{Column: "d", Dir: "DESC"},
				},
				Page: common.Page{Limit: 1, Offset: 2},
			},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM `name` WHERE `name`.`k` > ? GROUP BY `name`.`a`, `name`.`b` ORDER BY `name`.`c` ASC, `name`.`d` DESC LIMIT 1 OFFSET 2",
			wantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			tt.sel.Table = table
			stmt, args := newMust("mysql").Select(&tt.sel)
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign   common.Assignments
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "INSERT INTO `name` () VALUES ()",
			wantArgs: []interface{}(nil),
		},
		{
			assign:   common.Assignments{{Column: "c1", ColumnValue: 1}},
			wantStmt: "INSERT INTO `name` (`c1`) VALUES (?)",
			wantArgs: []interface{}{1},
		},
		{
			assign:   common.Assignments{{Column: "c1", ColumnValue: 1}, {Column: "c2", ColumnValue: ""}},
			wantStmt: "INSERT INTO `name` (`c1`, `c2`) VALUES (?, ?)",
			wantArgs: []interface{}{1, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			params := &common.InsertParams{Table: table, Assignments: tt.assign}
			stmt, args := newMust("mysql").Insert(params)
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args, " ")
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign   common.Assignments
		where    common.Where
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "UPDATE `name` SET",
			wantArgs: []interface{}(nil),
		},
		{
			assign:   common.Assignments{{Column: "c1", ColumnValue: 1}},
			wantStmt: "UPDATE `name` SET `c1` = ?",
			wantArgs: []interface{}{1},
		},
		{
			assign:   common.Assignments{{Column: "c1", ColumnValue: 1}, {Column: "c2", ColumnValue: ""}},
			wantStmt: "UPDATE `name` SET `c1` = ?, `c2` = ?",
			wantArgs: []interface{}{1, ""},
		},
		{
			assign:   common.Assignments{{Column: "c1", ColumnValue: 1}, {Column: "c2", ColumnValue: ""}},
			where:    common.NewWhere(orm.OpGt, "k", 3),
			wantStmt: "UPDATE `name` SET `c1` = ?, `c2` = ? WHERE `name`.`k` > ?",
			wantArgs: []interface{}{1, "", 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			params := &common.UpdateParams{Table: table, Assignments: tt.assign, Where: tt.where}
			stmt, args := newMust("mysql").Update(params)
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		where    common.Where
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "DELETE FROM `name`",
			wantArgs: []interface{}(nil),
		},
		{
			where:    common.NewWhere(orm.OpGt, "k", 3),
			wantStmt: "DELETE FROM `name` WHERE `name`.`k` > ?",
			wantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			params := &common.DeleteParams{Table: table, Where: tt.where}
			stmt, args := newMust("mysql").Delete(params)
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func newMust(name string) Dialect {
	d, err := New(name)
	if err != nil {
		panic(err)
	}
	return d
}

func reduceSpaces(s string) string {
	re := regexp.MustCompile("([ ]+)")
	return strings.Trim(re.ReplaceAllString(s, " "), " ")
}

type columner struct {
	columns []string
	count   bool
}

func (c *columner) Columns() []string {
	return c.columns
}

func (c *columner) Joins() []common.JoinParams {
	return nil
}

func (c *columner) Count() bool {
	return c.count
}
