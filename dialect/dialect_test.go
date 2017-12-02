package dialect

import (
	"regexp"
	"strings"
	"testing"

	"github.com/posener/orm"
	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/mysql"
	"github.com/posener/orm/dialect/sqlite3"
	"github.com/stretchr/testify/assert"
)

const table = "name"

func TestSelect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		sel             common.SelectParams
		mysqlWantStmt   string
		mysqlWantArgs   []interface{}
		sqlite3WantStmt string
		sqlite3WantArgs []interface{}
	}{
		{
			mysqlWantStmt:   "SELECT * FROM `name`",
			sqlite3WantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:             common.SelectParams{Columns: &columner{}},
			mysqlWantStmt:   "SELECT * FROM `name`",
			sqlite3WantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:             common.SelectParams{Columns: &columner{}, Page: common.Page{}},
			mysqlWantStmt:   "SELECT * FROM `name`",
			sqlite3WantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:             common.SelectParams{Columns: &columner{count: true}},
			mysqlWantStmt:   "SELECT COUNT(*) FROM `name`",
			sqlite3WantStmt: "SELECT COUNT(*) FROM 'name'",
		},
		{
			sel:             common.SelectParams{Columns: &columner{columns: []string{"a", "b", "c"}, count: true}},
			mysqlWantStmt:   "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM `name`",
			sqlite3WantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM 'name'",
		},
		{
			sel:             common.SelectParams{Columns: &columner{columns: []string{"a", "b", "c"}}},
			mysqlWantStmt:   "SELECT `name`.`a`, `name`.`b`, `name`.`c` FROM `name`",
			sqlite3WantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c` FROM 'name'",
		},
		{
			sel:             common.SelectParams{Page: common.Page{}},
			mysqlWantStmt:   "SELECT * FROM `name`",
			sqlite3WantStmt: "SELECT * FROM 'name'",
		},
		{
			sel:             common.SelectParams{Page: common.Page{Limit: 1}},
			mysqlWantStmt:   "SELECT * FROM `name` LIMIT 1",
			sqlite3WantStmt: "SELECT * FROM 'name' LIMIT 1",
		},
		{
			sel:             common.SelectParams{Page: common.Page{Limit: 1, Offset: 2}},
			mysqlWantStmt:   "SELECT * FROM `name` LIMIT 1 OFFSET 2",
			sqlite3WantStmt: "SELECT * FROM 'name' LIMIT 1 OFFSET 2",
		},
		{
			sel:             common.SelectParams{Page: common.Page{Offset: 1}},
			mysqlWantStmt:   "SELECT * FROM `name`",
			sqlite3WantStmt: "SELECT * FROM 'name'",
		},
		{
			sel: common.SelectParams{
				Columns: &columner{columns: []string{"a", "b", "c"}, count: true},
				Page:    common.Page{Limit: 1, Offset: 2},
			},
			mysqlWantStmt:   "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM `name` LIMIT 1 OFFSET 2",
			sqlite3WantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM 'name' LIMIT 1 OFFSET 2",
		},
		{
			sel: common.SelectParams{
				Groups: common.Groups{{Column: "a"}, {Column: "b"}},
			},
			mysqlWantStmt:   "SELECT * FROM `name` GROUP BY `name`.`a`, `name`.`b`",
			sqlite3WantStmt: "SELECT * FROM 'name' GROUP BY `name`.`a`, `name`.`b`",
		},
		{
			sel: common.SelectParams{
				Orders: common.Orders{
					{Column: "c", Dir: "ASC"},
					{Column: "d", Dir: "DESC"},
				},
			},
			mysqlWantStmt:   "SELECT * FROM `name` ORDER BY `name`.`c` ASC, `name`.`d` DESC",
			sqlite3WantStmt: "SELECT * FROM 'name' ORDER BY `name`.`c` ASC, `name`.`d` DESC",
		},
		{
			sel:             common.SelectParams{Where: common.NewWhere(orm.OpEq, "name", "k", 3)},
			mysqlWantStmt:   "SELECT * FROM `name` WHERE `name`.`k` = ?",
			mysqlWantArgs:   []interface{}{3},
			sqlite3WantStmt: "SELECT * FROM 'name' WHERE `name`.`k` = ?",
			sqlite3WantArgs: []interface{}{3},
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
			mysqlWantStmt:   "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM `name` WHERE `name`.`k` > ? GROUP BY `name`.`a`, `name`.`b` ORDER BY `name`.`c` ASC, `name`.`d` DESC LIMIT 1 OFFSET 2",
			mysqlWantArgs:   []interface{}{3},
			sqlite3WantStmt: "SELECT `name`.`a`, `name`.`b`, `name`.`c`, COUNT(*) FROM 'name' WHERE `name`.`k` > ? GROUP BY `name`.`a`, `name`.`b` ORDER BY `name`.`c` ASC, `name`.`d` DESC LIMIT 1 OFFSET 2",
			sqlite3WantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.mysqlWantStmt, func(t *testing.T) {
			tt.sel.Table = table
			stmt, args := new(mysql.Dialect).Select(&tt.sel)
			assert.Equal(t, tt.mysqlWantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.mysqlWantArgs, args)

			stmt, args = new(sqlite3.Dialect).Select(&tt.sel)
			assert.Equal(t, tt.sqlite3WantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.sqlite3WantArgs, args)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign          common.Assignments
		mysqlWantStmt   string
		mysqlWantArgs   []interface{}
		sqlite3WantStmt string
		sqlite3WantArgs []interface{}
	}{
		{
			mysqlWantStmt:   "INSERT INTO `name` () VALUES ()",
			mysqlWantArgs:   []interface{}(nil),
			sqlite3WantStmt: "INSERT INTO 'name' () VALUES ()",
			sqlite3WantArgs: []interface{}(nil),
		},
		{
			assign:          common.Assignments{{Column: "c1", Value: 1}},
			mysqlWantStmt:   "INSERT INTO `name` (`c1`) VALUES (?)",
			mysqlWantArgs:   []interface{}{1},
			sqlite3WantStmt: "INSERT INTO 'name' (`c1`) VALUES (?)",
			sqlite3WantArgs: []interface{}{1},
		},
		{
			assign:          common.Assignments{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			mysqlWantStmt:   "INSERT INTO `name` (`c1`, `c2`) VALUES (?, ?)",
			mysqlWantArgs:   []interface{}{1, ""},
			sqlite3WantStmt: "INSERT INTO 'name' (`c1`, `c2`) VALUES (?, ?)",
			sqlite3WantArgs: []interface{}{1, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.mysqlWantStmt, func(t *testing.T) {
			params := &common.InsertParams{Table: table, Assignments: tt.assign}
			stmt, args := new(mysql.Dialect).Insert(params)
			assert.Equal(t, tt.mysqlWantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.mysqlWantArgs, args, " ")

			stmt, args = new(sqlite3.Dialect).Insert(params)
			assert.Equal(t, tt.sqlite3WantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.sqlite3WantArgs, args, " ")
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign          common.Assignments
		where           common.Where
		mysqlWantStmt   string
		mysqlWantArgs   []interface{}
		sqlite3WantStmt string
		sqlite3WantArgs []interface{}
	}{
		{
			mysqlWantStmt:   "UPDATE `name` SET",
			mysqlWantArgs:   []interface{}(nil),
			sqlite3WantStmt: "UPDATE 'name' SET",
			sqlite3WantArgs: []interface{}(nil),
		},
		{
			assign:          common.Assignments{{Column: "c1", Value: 1}},
			mysqlWantStmt:   "UPDATE `name` SET `c1` = ?",
			mysqlWantArgs:   []interface{}{1},
			sqlite3WantStmt: "UPDATE 'name' SET `c1` = ?",
			sqlite3WantArgs: []interface{}{1},
		},
		{
			assign:          common.Assignments{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			mysqlWantStmt:   "UPDATE `name` SET `c1` = ?, `c2` = ?",
			mysqlWantArgs:   []interface{}{1, ""},
			sqlite3WantStmt: "UPDATE 'name' SET `c1` = ?, `c2` = ?",
			sqlite3WantArgs: []interface{}{1, ""},
		},
		{
			assign:          common.Assignments{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			where:           common.NewWhere(orm.OpGt, "name", "k", 3),
			mysqlWantStmt:   "UPDATE `name` SET `c1` = ?, `c2` = ? WHERE `name`.`k` > ?",
			mysqlWantArgs:   []interface{}{1, "", 3},
			sqlite3WantStmt: "UPDATE 'name' SET `c1` = ?, `c2` = ? WHERE `name`.`k` > ?",
			sqlite3WantArgs: []interface{}{1, "", 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.mysqlWantStmt, func(t *testing.T) {
			params := &common.UpdateParams{Table: table, Assignments: tt.assign, Where: tt.where}
			stmt, args := new(mysql.Dialect).Update(params)
			assert.Equal(t, tt.mysqlWantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.mysqlWantArgs, args)

			stmt, args = new(sqlite3.Dialect).Update(params)
			assert.Equal(t, tt.sqlite3WantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.sqlite3WantArgs, args, " ")
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		where           common.Where
		mysqlWantStmt   string
		mysqlWantArgs   []interface{}
		sqlite3WantStmt string
		sqlite3WantArgs []interface{}
	}{
		{
			mysqlWantStmt:   "DELETE FROM `name`",
			mysqlWantArgs:   []interface{}(nil),
			sqlite3WantStmt: "DELETE FROM 'name'",
			sqlite3WantArgs: []interface{}(nil),
		},
		{
			where:           common.NewWhere(orm.OpGt, "name", "k", 3),
			mysqlWantStmt:   "DELETE FROM `name` WHERE `name`.`k` > ?",
			mysqlWantArgs:   []interface{}{3},
			sqlite3WantStmt: "DELETE FROM 'name' WHERE `name`.`k` > ?",
			sqlite3WantArgs: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.mysqlWantStmt, func(t *testing.T) {
			params := &common.DeleteParams{Table: table, Where: tt.where}
			stmt, args := new(mysql.Dialect).Delete(params)
			assert.Equal(t, tt.mysqlWantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.mysqlWantArgs, args)

			stmt, args = new(sqlite3.Dialect).Delete(params)
			assert.Equal(t, tt.sqlite3WantStmt, reduceSpaces(stmt), " ")
			assert.Equal(t, tt.sqlite3WantArgs, args, " ")
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
