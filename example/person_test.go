// an example for a wanted autogenerated file
package example_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm/example"
	porm "github.com/posener/orm/example/personorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	p1 = example.Person{Name: "moshe", Age: 1}
	p2 = example.Person{Name: "haim", Age: 2}
	p3 = example.Person{Name: "zvika", Age: 3}
)

func TestPersonSelect(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err)
	defer db.Close()

	orm := porm.New(db)

	if testing.Verbose() {
		orm.Logger(t.Logf)
	}

	_, err = orm.Create().Exec()
	require.Nil(t, err, "Failed creating table")

	res, err := orm.Insert().SetName(p1.Name).SetAge(p1.Age).Exec()
	require.Nil(t, err, "Failed inserting")
	assertRowsAffected(t, 1, res)
	res, err = orm.Insert().SetName(p2.Name).SetAge(p2.Age).Exec()
	require.Nil(t, err, "Failed inserting")
	assertRowsAffected(t, 1, res)
	res, err = orm.InsertPerson(&p3).Exec()
	require.Nil(t, err, "Failed inserting")
	assertRowsAffected(t, 1, res)

	tests := []struct {
		q    porm.Querier
		want []example.Person
	}{
		{
			q:    orm.Select(),
			want: []example.Person{p1, p2, p3},
		},
		{
			q:    orm.Select().SelectName(),
			want: []example.Person{{Name: "moshe"}, {Name: "haim"}, {Name: "zvika"}},
		},
		{
			q:    orm.Select().SelectAge(),
			want: []example.Person{{Age: 1}, {Age: 2}, {Age: 3}},
		},
		{
			q:    orm.Select().SelectAge().SelectName(),
			want: []example.Person{p1, p2, p3},
		},
		{
			q:    orm.Select().SelectName().SelectAge(),
			want: []example.Person{p1, p2, p3},
		},
		{
			q:    orm.Select().Where(porm.WhereName(porm.OpEq, "moshe")),
			want: []example.Person{p1},
		},
		{
			q:    orm.Select().Where(porm.WhereName(porm.OpEq, "moshe").Or(porm.WhereAge(porm.OpEq, 2))),
			want: []example.Person{p1, p2},
		},
		{
			q:    orm.Select().Where(porm.WhereName(porm.OpEq, "moshe").And(porm.WhereAge(porm.OpEq, 1))),
			want: []example.Person{p1},
		},
		{
			q: orm.Select().Where(porm.WhereName(porm.OpEq, "moshe").And(porm.WhereAge(porm.OpEq, 2))),
		},
		{
			q:    orm.Select().Where(porm.WhereAge(porm.OpGE, 2)),
			want: []example.Person{p2, p3},
		},
		{
			q:    orm.Select().Where(porm.WhereAge(porm.OpGt, 2)),
			want: []example.Person{p3},
		},
		{
			q:    orm.Select().Where(porm.WhereAge(porm.OpLE, 2)),
			want: []example.Person{p1, p2},
		},
		{
			q:    orm.Select().Where(porm.WhereAge(porm.OpLt, 2)),
			want: []example.Person{p1},
		},
		{
			q:    orm.Select().Where(porm.WhereName(porm.OpNe, "moshe")),
			want: []example.Person{p2, p3},
		},
		{
			q:    orm.Select().Where(porm.WhereNameIn("moshe", "haim")),
			want: []example.Person{p1, p2},
		},
		{
			q:    orm.Select().Where(porm.WhereAgeBetween(0, 2)),
			want: []example.Person{p1, p2},
		},
		{
			q:    orm.Select().Limit(2),
			want: []example.Person{p1, p2},
		},
		{
			q:    orm.Select().Page(1, 1),
			want: []example.Person{p2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.q.String(), func(t *testing.T) {
			p, err := tt.q.Query()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.want, p)
		})
	}
}

func TestPersonCRUD(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err)
	defer db.Close()

	orm := porm.New(db)

	prepare(t, orm)
	ps, err := orm.Select().Query()
	require.Nil(t, err)
	assert.Equal(t, []example.Person{p1, p2, p3}, ps)

	// Test delete
	delete := orm.Delete().Where(porm.WhereName(porm.OpEq, "moshe"))
	res, err := delete.Exec()
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)
	ps, err = orm.Select().Query()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []example.Person{p2, p3}, ps)
	ps, err = orm.Select().Where(porm.WhereName(porm.OpEq, "moshe")).Query()
	require.Nil(t, err)
	assert.Equal(t, []example.Person(nil), ps)

	// Test Update
	update := orm.Update().SetName("Jonney").Where(porm.WhereName(porm.OpEq, "zvika"))
	res, err = update.Exec()
	require.Nil(t, err)
	assertRowsAffected(t, 1, res)

	ps, err = orm.Select().Where(porm.WhereName(porm.OpEq, "Jonney")).Query()
	require.Nil(t, err)
	assert.Equal(t, []example.Person{{Name: "Jonney", Age: 3}}, ps)
}

func prepare(t *testing.T, orm porm.API) {
	_, err := orm.Create().Exec()
	require.Nil(t, err)
	for _, p := range []example.Person{p1, p2, p3} {
		res, err := orm.InsertPerson(&p).Exec()
		require.Nil(t, err, "Failed inserting")
		assertRowsAffected(t, 1, res)
	}
}

func assertRowsAffected(t *testing.T, wantRows int64, result sql.Result) {
	gotRows, err := result.RowsAffected()
	require.Nil(t, err)
	assert.Equal(t, wantRows, gotRows)
}
