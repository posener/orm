package sqlite3

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
		sel      orm.Select
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:      orm.Select{Columns: &columner{}},
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:      orm.Select{Columns: &columner{}, Page: orm.Page{}},
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:      orm.Select{Columns: &columner{count: true}},
			wantStmt: "SELECT COUNT(*) FROM 'name'",
		},
		{
			sel:      orm.Select{Columns: &columner{columns: []string{"a", "b", "c"}, count: true}},
			wantStmt: "SELECT `a`, `b`, `c`, COUNT(*) FROM 'name'",
		},
		{
			sel:      orm.Select{Columns: &columner{columns: []string{"a", "b", "c"}}},
			wantStmt: "SELECT `a`, `b`, `c` FROM 'name'",
		},
		{
			sel:      orm.Select{Page: orm.Page{}},
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:      orm.Select{Page: orm.Page{Limit: 1}},
			wantStmt: "SELECT * FROM 'name' LIMIT 1",
		},
		{
			sel:      orm.Select{Page: orm.Page{Limit: 1, Offset: 2}},
			wantStmt: "SELECT * FROM 'name' LIMIT 1 OFFSET 2",
		},
		{
			sel:      orm.Select{Page: orm.Page{Offset: 1}},
			wantStmt: "SELECT * FROM 'name'",
		},
		{
			sel: orm.Select{
				Columns: &columner{columns: []string{"a", "b", "c"}, count: true},
				Page:    orm.Page{Limit: 1, Offset: 2},
			},
			wantStmt: "SELECT `a`, `b`, `c`, COUNT(*) FROM 'name' LIMIT 1 OFFSET 2",
		},
		{
			sel: orm.Select{
				Groups: orm.Groups{{Column: "a"}, {Column: "b"}},
			},
			wantStmt: "SELECT * FROM 'name' GROUP BY `a`, `b`",
		},
		{
			sel: orm.Select{
				Orders: orm.Orders{
					{Column: "c", Dir: "ASC"},
					{Column: "d", Dir: "DESC"},
				},
			},
			wantStmt: "SELECT * FROM 'name' ORDER BY `c` ASC, `d` DESC",
		},
		{
			sel:      orm.Select{Where: orm.NewWhere(orm.OpEq, "k", 3)},
			wantStmt: "SELECT * FROM 'name' WHERE `k` = ?",
			wantArgs: []interface{}{3},
		},
		{
			sel: orm.Select{
				Columns: &columner{columns: []string{"a", "b", "c"}, count: true},
				Where:   orm.NewWhere(orm.OpGt, "k", 3),
				Groups:  orm.Groups{{Column: "a"}, {Column: "b"}},
				Orders: orm.Orders{
					{Column: "c", Dir: "ASC"},
					{Column: "d", Dir: "DESC"},
				},
				Page: orm.Page{Limit: 1, Offset: 2},
			},
			wantStmt: "SELECT `a`, `b`, `c`, COUNT(*) FROM 'name' WHERE `k` > ? GROUP BY `a`, `b` ORDER BY `c` ASC, `d` DESC LIMIT 1 OFFSET 2",
			wantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			tt.sel.Table = table
			stmt, args := Select(&tt.sel)
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign   orm.Assignments
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "INSERT INTO 'name' () VALUES ()",
			wantArgs: []interface{}(nil),
		},
		{
			assign:   orm.Assignments{{Column: "c1", Value: 1}},
			wantStmt: "INSERT INTO 'name' (`c1`) VALUES (?)",
			wantArgs: []interface{}{1},
		},
		{
			assign:   orm.Assignments{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			wantStmt: "INSERT INTO 'name' (`c1`, `c2`) VALUES (?, ?)",
			wantArgs: []interface{}{1, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			stmt, args := Insert(&orm.Insert{Table: table, Assignments: tt.assign})
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args, " ")
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign   orm.Assignments
		where    orm.Where
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "UPDATE 'name' SET",
			wantArgs: []interface{}(nil),
		},
		{
			assign:   orm.Assignments{{Column: "c1", Value: 1}},
			wantStmt: "UPDATE 'name' SET `c1` = ?",
			wantArgs: []interface{}{1},
		},
		{
			assign:   orm.Assignments{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			wantStmt: "UPDATE 'name' SET `c1` = ?, `c2` = ?",
			wantArgs: []interface{}{1, ""},
		},
		{
			assign:   orm.Assignments{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			where:    orm.NewWhere(orm.OpGt, "k", 3),
			wantStmt: "UPDATE 'name' SET `c1` = ?, `c2` = ? WHERE `k` > ?",
			wantArgs: []interface{}{1, "", 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			stmt, args := Update(&orm.Update{Table: table, Assignments: tt.assign, Where: tt.where})
			assert.Equal(t, tt.wantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		where    orm.Where
		wantStmt string
		wantArgs []interface{}
	}{
		{
			wantStmt: "DELETE FROM 'name'",
			wantArgs: []interface{}(nil),
		},
		{
			where:    orm.NewWhere(orm.OpGt, "k", 3),
			wantStmt: "DELETE FROM 'name' WHERE `k` > ?",
			wantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantStmt, func(t *testing.T) {
			stmt, args := Delete(&orm.Delete{Table: table, Where: tt.where})
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

func (c *columner) Count() bool {
	return c.count
}
