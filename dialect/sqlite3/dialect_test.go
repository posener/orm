package sqlite3

import (
	"regexp"
	"strings"
	"testing"

	"github.com/posener/orm"
	"github.com/posener/orm/common"
	"github.com/stretchr/testify/assert"
)

const table = "name"

func TestSelect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		sel      common.SelectParams
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:      common.SelectParams{Columns: &columner{}},
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:      common.SelectParams{Columns: &columner{}, Page: common.Page{}},
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:      common.SelectParams{Columns: &columner{count: true}},
			wantStmt: "SELECT COUNT(*) FROM 'name'",
		},
		{
			sel:      common.SelectParams{Columns: &columner{columns: []string{"a", "b", "c"}, count: true}},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM 'name'",
		},
		{
			sel:      common.SelectParams{Columns: &columner{columns: []string{"a", "b", "c"}}},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c` FROM 'name'",
		},
		{
			sel:      common.SelectParams{Page: common.Page{}},
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:      common.SelectParams{Page: common.Page{Limit: 1}},
			wantStmt: "SELECT * FROM 'name' LIMIT 1",
		},
		{
			sel:      common.SelectParams{Page: common.Page{Limit: 1, Offset: 2}},
			wantStmt: "SELECT * FROM 'name' LIMIT 1 OFFSET 2",
		},
		{
			sel:      common.SelectParams{Page: common.Page{Offset: 1}},
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel: common.SelectParams{
				Columns: &columner{columns: []string{"a", "b", "c"}, count: true},
				Page:    common.Page{Limit: 1, Offset: 2},
			},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM 'name' LIMIT 1 OFFSET 2",
		},
		{
			sel: common.SelectParams{
				Groups: common.Groups{{Column: "a"}, {Column: "b"}},
			},
			wantStmt: "SELECT * FROM 'name' GROUP BY `name`.`a`, `name`.`b`",
		},
		{
			sel: common.SelectParams{
				Orders: common.Orders{
					{Column: "c", Dir: "ASC"},
					{Column: "d", Dir: "DESC"},
				},
			},
			wantStmt: "SELECT * FROM 'name' ORDER BY `name`.`c` ASC, `name`.`d` DESC",
		},
		{
			sel:      common.SelectParams{Where: common.NewWhere(orm.OpEq, "name", "k", 3)},
			wantStmt: "SELECT * FROM 'name' WHERE `name`.`k` = ?",
			wantArgs: []interface{}{3},
		},
		{
			sel: common.SelectParams{
				Columns: &columner{columns: []string{"a", "b", "c"}, count: true},
				Where:   common.NewWhere(orm.OpGt, "name", "k", 3),
				Groups:  common.Groups{{Column: "a"}, {Column: "b"}},
				Orders: common.Orders{
					{Column: "c", Dir: "ASC"},
					{Column: "d", Dir: "DESC"},
				},
				Page: common.Page{Limit: 1, Offset: 2},
			},
			wantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM 'name' WHERE `name`.`k` > ? GROUP BY `name`.`a`, `name`.`b` ORDER BY `name`.`c` ASC, `name`.`d` DESC LIMIT 1 OFFSET 2",
			wantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			tt.sel.Table = table
			d := &Dialect{}
			stmt, args := d.Select(&tt.sel)
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
			wantStmt: "INSERT INTO 'name' () VALUES ()",
			wantArgs: []interface{}(nil),
		},
		{
			assign:   common.Assignments{{Column: "c1", Value: 1}},
			wantStmt: "INSERT INTO 'name' (`c1`) VALUES (?)",
			wantArgs: []interface{}{1},
		},
		{
			assign:   common.Assignments{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			wantStmt: "INSERT INTO 'name' (`c1`, `c2`) VALUES (?, ?)",
			wantArgs: []interface{}{1, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			d := &Dialect{}
			stmt, args := d.Insert(&common.InsertParams{Table: table, Assignments: tt.assign})
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
			wantStmt: "UPDATE 'name' SET",
			wantArgs: []interface{}(nil),
		},
		{
			assign:   common.Assignments{{Column: "c1", Value: 1}},
			wantStmt: "UPDATE 'name' SET `c1` = ?",
			wantArgs: []interface{}{1},
		},
		{
			assign:   common.Assignments{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			wantStmt: "UPDATE 'name' SET `c1` = ?, `c2` = ?",
			wantArgs: []interface{}{1, ""},
		},
		{
			assign:   common.Assignments{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			where:    common.NewWhere(orm.OpGt, "name", "k", 3),
			wantStmt: "UPDATE 'name' SET `c1` = ?, `c2` = ? WHERE `name`.`k` > ?",
			wantArgs: []interface{}{1, "", 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			d := &Dialect{}
			stmt, args := d.Update(&common.UpdateParams{Table: table, Assignments: tt.assign, Where: tt.where})
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
			wantStmt: "DELETE FROM 'name'",
			wantArgs: []interface{}(nil),
		},
		{
			where:    common.NewWhere(orm.OpGt, "name", "k", 3),
			wantStmt: "DELETE FROM 'name' WHERE `name`.`k` > ?",
			wantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			d := &Dialect{}
			stmt, args := d.Delete(&common.DeleteParams{Table: table, Where: tt.where})
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args)
		})
	}
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

func (c *columner) Joins() []common.Join {
	return nil
}

func (c *columner) Count() bool {
	return c.count
}
