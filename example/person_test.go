// an example for a wanted autogenerated file
package example_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm/example"
	porm "github.com/posener/orm/example/personorm"
	"github.com/stretchr/testify/assert"
)

var (
	p1 = example.Person{Name: "moshe", Age: 1}
	p2 = example.Person{Name: "haim", Age: 2}
	p3 = example.Person{Name: "zvika", Age: 3}
)

func TestPersonSelect(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = porm.Create().Exec(db)
	if err != nil {
		t.Fatalf("Failed creating table: %s", err)
	}

	err = porm.Insert().SetName(p1.Name).SetAge(p1.Age).Exec(db)
	if err != nil {
		t.Fatalf("Failed inserting: %s", err)
	}
	err = porm.Insert().SetName(p2.Name).SetAge(p2.Age).Exec(db)
	if err != nil {
		t.Fatalf("Failed inserting: %s", err)
	}
	err = porm.InsertPerson(&p3).Exec(db)
	if err != nil {
		t.Fatalf("Failed inserting: %s", err)
	}

	tests := []struct {
		q    *porm.Select
		want []example.Person
	}{
		{
			q:    porm.Query(),
			want: []example.Person{p1, p2, p3},
		},
		{
			q:    porm.Query().SelectName(),
			want: []example.Person{{Name: "moshe"}, {Name: "haim"}, {Name: "zvika"}},
		},
		{
			q:    porm.Query().SelectAge(),
			want: []example.Person{{Age: 1}, {Age: 2}, {Age: 3}},
		},
		{
			q:    porm.Query().SelectAge().SelectName(),
			want: []example.Person{p1, p2, p3},
		},
		{
			q:    porm.Query().SelectName().SelectAge(),
			want: []example.Person{p1, p2, p3},
		},
		{
			q:    porm.Query().Where(porm.WhereName(porm.OpEq, "moshe")),
			want: []example.Person{p1},
		},
		{
			q:    porm.Query().Where(porm.WhereName(porm.OpEq, "moshe").Or(porm.WhereAge(porm.OpEq, 2))),
			want: []example.Person{p1, p2},
		},
		{
			q:    porm.Query().Where(porm.WhereName(porm.OpEq, "moshe").And(porm.WhereAge(porm.OpEq, 1))),
			want: []example.Person{p1},
		},
		{
			q: porm.Query().Where(porm.WhereName(porm.OpEq, "moshe").And(porm.WhereAge(porm.OpEq, 2))),
		},
		{
			q:    porm.Query().Where(porm.WhereAge(porm.OpGE, 2)),
			want: []example.Person{p2, p3},
		},
		{
			q:    porm.Query().Where(porm.WhereAge(porm.OpGt, 2)),
			want: []example.Person{p3},
		},
		{
			q:    porm.Query().Where(porm.WhereAge(porm.OpLE, 2)),
			want: []example.Person{p1, p2},
		},
		{
			q:    porm.Query().Where(porm.WhereAge(porm.OpLt, 2)),
			want: []example.Person{p1},
		},
		{
			q:    porm.Query().Where(porm.WhereName(porm.OpNe, "moshe")),
			want: []example.Person{p2, p3},
		},
		{
			q:    porm.Query().Where(porm.WhereNameIn("moshe", "haim")),
			want: []example.Person{p1, p2},
		},
		{
			q:    porm.Query().Where(porm.WhereAgeBetween(0, 2)),
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

func TestPersonCRUD(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	prepare(t, db)
	ps, err := porm.Query().Exec(db)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []example.Person{p1, p2, p3}, ps)
	porm.Delete().Where(porm.WhereName(porm.OpEq, "moshe")).Exec(db)
	ps, err = porm.Query().Exec(db)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []example.Person{p2, p3}, ps)
	ps, err = porm.Query().Where(porm.WhereName(porm.OpEq, "moshe")).Exec(db)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []example.Person(nil), ps)
}

func prepare(t *testing.T, db porm.SQLExecer) {
	err := porm.Create().Exec(db)
	if err != nil {
		t.Fatalf("Failed creating table: %s", err)
	}
	for _, p := range []example.Person{p1, p2, p3} {
		err = porm.InsertPerson(&p).Exec(db)
		if err != nil {
			t.Fatalf("Failed inserting: %s", err)
		}
	}
}
