package dialect

import (
	"regexp"
	"strings"
	"testing"

	"github.com/posener/orm"
	"github.com/stretchr/testify/assert"
)

const table = "name"

func TestSelect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		sel      SelectParams
		wantStmt string
		wantArgs []interface{}
	}{
		{
			sel:      SelectParams{Columns: map[string]bool{"col": true}, OrderedColumns: []string{"col"}},
			wantStmt: "SELECT `name`.`col` FROM `name`",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"col": true},
				OrderedColumns: []string{"col"},
				Page:           Page{},
			},
			wantStmt: "SELECT `name`.`col` FROM `name`",
		},
		{
			sel:      SelectParams{Count: true},
			wantStmt: "SELECT COUNT(*) FROM `name`",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"a": true, "b": true, "c": true},
				OrderedColumns: []string{"a", "b", "c"},
				Count:          true,
			},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM `name`",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"a": true, "b": true, "c": true},
				OrderedColumns: []string{"a", "c", "b"},
				Count:          true,
			},
			wantStmt: "SELECT `name`.`a`, `name`.`c`, `name`.`b`, COUNT(*) FROM `name`",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"a": true, "b": true, "c": true},
				OrderedColumns: []string{"a", "c", "b"},
			},
			wantStmt: "SELECT `name`.`a`, `name`.`c`, `name`.`b` FROM `name`",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"col": true},
				OrderedColumns: []string{"col"},
				Page:           Page{},
			},
			wantStmt: "SELECT `name`.`col` FROM `name`",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"col": true},
				OrderedColumns: []string{"col"},
				Page:           Page{Limit: 1},
			},
			wantStmt: "SELECT `name`.`col` FROM `name` LIMIT 1",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"col": true},
				OrderedColumns: []string{"col"},
				Page:           Page{Limit: 1, Offset: 2},
			},
			wantStmt: "SELECT `name`.`col` FROM `name` LIMIT 1 OFFSET 2",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"col": true},
				OrderedColumns: []string{"col"},
				Page:           Page{Offset: 1},
			},
			wantStmt: "SELECT `name`.`col` FROM `name`",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"a": true, "b": true, "c": true},
				OrderedColumns: []string{"a", "b", "c"},
				Count:          true,
				Page:           Page{Limit: 1, Offset: 2},
			},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM `name` LIMIT 1 OFFSET 2",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"a": true},
				OrderedColumns: []string{"a"},
				Groups:         Groups{{Column: "a"}, {Column: "b"}},
			},
			wantStmt: "SELECT `name`.`a` FROM `name` GROUP BY `name`.`a`, `name`.`b`",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"c": true},
				OrderedColumns: []string{"c"},
				Orders: Orders{
					{Column: "c", Dir: "ASC"},
					{Column: "d", Dir: "DESC"},
				},
			},
			wantStmt: "SELECT `name`.`c` FROM `name` ORDER BY `name`.`c` ASC, `name`.`d` DESC",
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"k": true},
				OrderedColumns: []string{"k"},
				Where:          NewWhere(orm.OpEq, "k", 3),
			},
			wantStmt: "SELECT `name`.`k` FROM `name` WHERE `name`.`k` = ?",
			wantArgs: []interface{}{3},
		},
		{
			sel: SelectParams{
				Columns:        map[string]bool{"a": true, "b": true, "c": true},
				OrderedColumns: []string{"a", "b", "c"},
				Count:          true,
				Where:          NewWhere(orm.OpGt, "k", 3),
				Groups:         Groups{{Column: "a"}, {Column: "b"}},
				Orders: Orders{
					{Column: "c", Dir: "ASC"},
					{Column: "d", Dir: "DESC"},
				},
				Page: Page{Limit: 1, Offset: 2},
			},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM `name` WHERE `name`.`k` > ? GROUP BY `name`.`a`, `name`.`b` ORDER BY `name`.`c` ASC, `name`.`d` DESC LIMIT 1 OFFSET 2",
			wantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			tt.sel.Table = table
			stmt, args := Get("mysql").Select(&tt.sel)
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign   Assignments
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "INSERT INTO `name` () VALUES ()",
			wantArgs: []interface{}(nil),
		},
		{
			assign:   Assignments{{Column: "c1", ColumnValue: 1}},
			wantStmt: "INSERT INTO `name` (`c1`) VALUES (?)",
			wantArgs: []interface{}{1},
		},
		{
			assign:   Assignments{{Column: "c1", ColumnValue: 1}, {Column: "c2", ColumnValue: ""}},
			wantStmt: "INSERT INTO `name` (`c1`, `c2`) VALUES (?, ?)",
			wantArgs: []interface{}{1, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			params := &InsertParams{Table: table, Assignments: tt.assign}
			stmt, args := Get("mysql").Insert(params)
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args, " ")
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign   Assignments
		where    Where
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "UPDATE `name` SET",
			wantArgs: []interface{}(nil),
		},
		{
			assign:   Assignments{{Column: "c1", ColumnValue: 1}},
			wantStmt: "UPDATE `name` SET `c1` = ?",
			wantArgs: []interface{}{1},
		},
		{
			assign:   Assignments{{Column: "c1", ColumnValue: 1}, {Column: "c2", ColumnValue: ""}},
			wantStmt: "UPDATE `name` SET `c1` = ?, `c2` = ?",
			wantArgs: []interface{}{1, ""},
		},
		{
			assign:   Assignments{{Column: "c1", ColumnValue: 1}, {Column: "c2", ColumnValue: ""}},
			where:    NewWhere(orm.OpGt, "k", 3),
			wantStmt: "UPDATE `name` SET `c1` = ?, `c2` = ? WHERE `name`.`k` > ?",
			wantArgs: []interface{}{1, "", 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			params := &UpdateParams{Table: table, Assignments: tt.assign, Where: tt.where}
			stmt, args := Get("mysql").Update(params)
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		where    Where
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "DELETE FROM `name`",
			wantArgs: []interface{}(nil),
		},
		{
			where:    NewWhere(orm.OpGt, "k", 3),
			wantStmt: "DELETE FROM `name` WHERE `name`.`k` > ?",
			wantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			params := &DeleteParams{Table: table, Where: tt.where}
			stmt, args := Get("mysql").Delete(params)
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestColumnsJoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		p        SelectParams
		wantCols string
		wantJoin string
	}{
		{
			p: SelectParams{Table: "table"},
		},
		{
			p:        SelectParams{Table: "table", Count: true},
			wantCols: "COUNT(*)",
		},
		{
			p: SelectParams{
				Table:          "table",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
			},
			wantCols: "`table`.`a`, `table`.`b`",
		},
		{
			p: SelectParams{
				Table:          "table",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Count:          true,
			},
			wantCols: "`table`.`a`, `table`.`b`, COUNT(*)",
		},
		{
			p: SelectParams{
				Table:          "table",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"b", "a"},
				Count:          true,
			},
			wantCols: "`table`.`b`, `table`.`a`, COUNT(*)",
		},
		{
			p: SelectParams{
				Table:          "table",
				Columns:        map[string]bool{"a": true, "b": true},
				Count:          true,
				OrderedColumns: []string{"a", "b"},
			},
			wantCols: "`table`.`a`, `table`.`b`, COUNT(*)",
		},
		{
			p: SelectParams{
				Table: "A",
				Joins: []JoinParams{
					{
						SelectParams: SelectParams{
							Table: "B",
						},
						Pairings: []Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "",
			wantJoin: "LEFT OUTER JOIN `B` AS `A_B_id` ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: SelectParams{
				Table: "A",
				Count: true,
				Joins: []JoinParams{
					{
						SelectParams: SelectParams{Table: "B"},
						Pairings:     []Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "COUNT(*)",
			wantJoin: "LEFT OUTER JOIN `B` AS `A_B_id` ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: SelectParams{
				Table:          "A",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Joins: []JoinParams{
					{
						SelectParams: SelectParams{Table: "B"},
						Pairings:     []Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`",
			wantJoin: "LEFT OUTER JOIN `B` AS `A_B_id` ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: SelectParams{
				Table: "A",
				Joins: []JoinParams{
					{
						SelectParams: SelectParams{
							Table:          "B",
							Columns:        map[string]bool{"c": true, "d": true},
							OrderedColumns: []string{"c", "d"},
						},
						Pairings: []Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A_B_id`.`c`, `A_B_id`.`d`",
			wantJoin: "LEFT OUTER JOIN `B` AS `A_B_id` ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: SelectParams{
				Table:          "A",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Joins: []JoinParams{
					{
						SelectParams: SelectParams{
							Table:          "B",
							Columns:        map[string]bool{"c": true, "d": true},
							OrderedColumns: []string{"c", "d"},
						},
						Pairings: []Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`",
			wantJoin: "LEFT OUTER JOIN `B` AS `A_B_id` ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: SelectParams{
				Table:          "A",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Joins: []JoinParams{
					{
						SelectParams: SelectParams{
							Table:          "B",
							Columns:        map[string]bool{"c": true, "d": true},
							OrderedColumns: []string{"c", "d"},
						},
						Pairings: []Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
					{
						SelectParams: SelectParams{
							Table:          "C",
							Columns:        map[string]bool{"e": true, "f": true},
							OrderedColumns: []string{"e", "f"},
						},
						Pairings: []Pairing{{Column: "C_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`, `A_C_id`.`e`, `A_C_id`.`f`",
			wantJoin: "LEFT OUTER JOIN `B` AS `A_B_id` ON (`A`.`B_id` = `A_B_id`.`id`) LEFT OUTER JOIN `C` AS `A_C_id` ON (`A`.`C_id` = `A_C_id`.`id`)",
		},
		{
			p: SelectParams{
				Table:          "A",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Joins: []JoinParams{
					{
						SelectParams: SelectParams{
							Table:          "B",
							Columns:        map[string]bool{"c": true, "d": true},
							OrderedColumns: []string{"c", "d"},
							Joins: []JoinParams{
								{
									SelectParams: SelectParams{
										Table:          "D",
										Columns:        map[string]bool{"g": true, "h": true},
										OrderedColumns: []string{"g", "h"},
									},
									Pairings: []Pairing{{Column: "D_id", JoinedColumn: "id"}},
								},
							},
						},
						Pairings: []Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
					{
						SelectParams: SelectParams{
							Table:          "C",
							Columns:        map[string]bool{"e": true, "f": true},
							OrderedColumns: []string{"e", "f"},
						},
						Pairings: []Pairing{{Column: "C_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`, `A_B_id_D_id`.`g`, `A_B_id_D_id`.`h`, `A_C_id`.`e`, `A_C_id`.`f`",
			wantJoin: "LEFT OUTER JOIN `B` AS `A_B_id` ON (`A`.`B_id` = `A_B_id`.`id`) LEFT OUTER JOIN `D` AS `A_B_id_D_id` ON (`A_B_id`.`D_id` = `A_B_id_D_id`.`id`) LEFT OUTER JOIN `C` AS `A_C_id` ON (`A`.`C_id` = `A_C_id`.`id`)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantCols, func(t *testing.T) {
			d := Get("mysql").(*dialect)

			cols := newBuilder(d, "")
			cols.selectColumns(&tt.p)
			assert.Equal(t, tt.wantCols, strings.Trim(cols.Statement(), " "))

			j := newBuilder(d, "")
			j.join(&tt.p)
			assert.Equal(t, tt.wantJoin, strings.Trim(j.Statement(), " "))
		})
	}
}

func reduceSpaces(s string) string {
	re := regexp.MustCompile("([ ]+)")
	return strings.Trim(re.ReplaceAllString(s, " "), " ")
}
