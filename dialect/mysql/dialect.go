package mysql

import (
	"fmt"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/format"
)

type Dialect struct{}

func (*Dialect) Name() string {
	return "mysql"
}

// Insert returns an SQL CREATE statement and arguments
func (d *Dialect) Create(p *common.Create) (string, []interface{}) {
	stmt := fmt.Sprintf("CREATE TABLE %s `%s` ( %s )",
		format.IfNotExists(p.IfNotExists),
		p.Table,
		p.ColumnsStatement,
	)

	return stmt, nil
}

// Insert returns an SQL INSERT statement and arguments
func (d *Dialect) Insert(p *common.Insert) (string, []interface{}) {
	stmt := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
		p.Table,
		format.AssignColumns(p.Assignments),
		common.QMarks(len(p.Assignments)),
	)

	var args []interface{}
	if p.Assignments != nil {
		args = append(args, p.Assignments.Args()...)
	}

	return stmt, args
}

// Select returns an SQL SELECT statement and arguments
func (d *Dialect) Select(p *common.Select) (string, []interface{}) {
	stmt := fmt.Sprintf("SELECT %s FROM `%s` %s %s %s %s",
		format.Columns(p.Columns),
		p.Table,
		format.Where(p.Where),
		format.GroupBy(p.Groups),
		format.OrderBy(p.Orders),
		format.Page(p.Page),
	)

	var args []interface{}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}

// Delete returns an SQL DELETE statement and arguments
func (d *Dialect) Delete(p *common.Delete) (string, []interface{}) {
	stmt := fmt.Sprintf("DELETE FROM `%s` %s",
		p.Table,
		format.Where(p.Where),
	)

	var args []interface{}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}

// Update returns an SQL UPDATE statement and arguments
func (d *Dialect) Update(p *common.Update) (string, []interface{}) {
	stmt := fmt.Sprintf("UPDATE `%s` SET %s %s",
		p.Table,
		format.AssignSets(p.Assignments),
		format.Where(p.Where),
	)

	var args []interface{}
	if p.Assignments != nil {
		args = append(args, p.Assignments.Args()...)
	}
	if p.Where != nil {
		args = append(args, p.Where.Args()...)
	}

	return stmt, args
}
