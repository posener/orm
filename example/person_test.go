// an example for a wanted autogenerated file
package example_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm/example"
	porm "github.com/posener/orm/example/personorm"
	"github.com/posener/orm/where"
	"github.com/stretchr/testify/assert"
)

func TestPersonSelect(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	mustExec(db.Exec("CREATE TABLE person (name VARCHAR(255), age INT)"))
	mustExec(db.Exec("INSERT INTO person (name, age) VALUES ('moshe', 1)"))
	mustExec(db.Exec("INSERT INTO person (name, age) VALUES ('haim', 2)"))
	mustExec(db.Exec("INSERT INTO person (name, age) VALUES ('zvika', 3)"))

	p1 := example.Person{Name: "moshe", Age: 1}
	p2 := example.Person{Name: "haim", Age: 2}
	p3 := example.Person{Name: "zvika", Age: 3}

	tests := []struct {
		q    *porm.Query
		want []example.Person
	}{
		{
			q:    &porm.Query{},
			want: []example.Person{p1, p2, p3},
		},
		{
			q:    &porm.Query{Select: &porm.Select{Name: true}},
			want: []example.Person{{Name: "moshe"}, {Name: "haim"}, {Name: "zvika"}},
		},
		{
			q:    &porm.Query{Select: &porm.Select{Age: true}},
			want: []example.Person{{Age: 1}, {Age: 2}, {Age: 3}},
		},
		{
			q:    &porm.Query{Where: porm.WhereName(where.OpEq, "moshe")},
			want: []example.Person{p1},
		},
		{
			q:    &porm.Query{Where: porm.WhereName(where.OpEq, "moshe").Or(porm.WhereAge(where.OpEq, 2))},
			want: []example.Person{p1, p2},
		},
		{
			q:    &porm.Query{Where: porm.WhereName(where.OpEq, "moshe").And(porm.WhereAge(where.OpEq, 1))},
			want: []example.Person{p1},
		},
		{
			q: &porm.Query{Where: porm.WhereName(where.OpEq, "moshe").And(porm.WhereAge(where.OpEq, 2))},
		},
		{
			q:    &porm.Query{Where: porm.WhereName(where.OpNe, "moshe")},
			want: []example.Person{p2, p3},
		},
		{
			q:    &porm.Query{Where: porm.WhereNameIn("moshe", "haim")},
			want: []example.Person{p1, p2},
		},
		{
			q:    &porm.Query{Where: porm.WhereAgeBetween(0, 2)},
			want: []example.Person{p1, p2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.q.String(), func(t *testing.T) {
			p, err := tt.q.Exec(db)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.want, p)
		})
	}
}

func mustExec(result sql.Result, err error) sql.Result {
	if err != nil {
		panic(err)
	}
	return result
}

// CREATE {{.Table}} {{.ColumnWithTypes}}
// SELECT {{.Columns}} {{.Where}} FROM {{Table}}, where args
// DELETE {{.Where}} FROM {{.Table}}, where args

//type Statement struct {
//	Template *template.Template
//	Table string
//}
//
//func (s *Statement) Write(w io.Writer) {
//	s.Template.Execute(w, s)
//}
//
//func (s *Statement) String() string {
//	b := bytes.NewBuffer(nil)
//	s.Write(b)
//	return b.String()
//}
//
//type CreateStatement struct {
//	Statement
//	ColumnTypes
//}
//
//var (
//	Create = Statement{
//		Template: template.Must(template.New("create").Parse("CREATE moshe ")),
//	}
//	Select = Statement{
//		Template: template.Must(template.New("select").Parse("SELECT {{.Columns}} {{.Where}} FROM {{.Table}}")),
//	}
//	Delete = Statement{
//		Template: template.Must(template.New("delete").Parse("DELETE {{.Where}} FROM {{.Table}}")),
//	}
//)
