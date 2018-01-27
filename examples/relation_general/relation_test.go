package relationone2many

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
	"github.com/posener/orm/examples"
)

// ExampleRelationGeneral demonstrates general relationship between structs with ORM
//
// In this example we use the following structs:
//
//type One struct {
//	ID   int64 `sql:"primary key;auto increment"`
//	Name string
//	// A relation to Two, since Two doesn't have a field which is a pointer to One,
//	// This is considered to be a general relation, which will be stored in a table
//	// called `rel_one_two` if created.
//	Two []Two
//	// A relation to Two with a "relation name" modifier, the relation table will be
//	// called "mutual" since this is the relation name.
//	// Since Two also has a field with the same relation name, this relation is mutual,
//	// and adding a relation with One's orm will be shown in Two's queries and vice versa.
//	Mutual []Two `sql:"relation name:mutual"`
//}
//type Two struct {
//	ID   int64 `sql:"primary key;auto increment"`
//	Name string
//	// A relation to One, since One doesn't have a field which is a pointer to Two,
//	// This is considered to be a general relation, which will be stored in a table
//	// called `rel_two_one` if created.
//	One []One
//	// A relation to One with a "relation name" modifier, the relation table will be
//	// called "mutual" since this is the relation name.
//	// Since One also has a field with the same relation name, this relation is mutual,
//	// and adding a relation with Two's orm will be shown in One's queries and vice versa.
//	Mutual []One `sql:"relation name:mutual"`
//}

func ExampleRelationGeneral() {
	conn, err := orm.Open("sqlite3", ":memory:", orm.OptLogger(log.Printf))
	examples.PanicOnErr(err)
	defer conn.Close()

	one, err := NewOneORM(conn)
	examples.PanicOnErr(err)
	two, err := NewTwoORM(conn)
	examples.PanicOnErr(err)

	// Create SQL tables
	examples.PanicOnErr(one.Create().Exec())
	examples.PanicOnErr(two.Create().Exec())

	// Since One and Two have relation tables, we need to create them also.
	// Since they have a common table named "mutual", we need to create the second one
	// with the IfNotExists modifier.
	examples.PanicOnErr(one.Create().Relations().Exec())
	examples.PanicOnErr(two.Create().Relations().IfNotExists().Exec())

	// Let's create rows in both types so we can start making relations
	const n = 5
	var (
		ones = make([]*One, n)
		twos = make([]*Two, n)
	)
	for i := 0; i < n; i++ {
		ones[i], err = one.Insert().SetName(fmt.Sprintf("one %d", i)).Exec()
		examples.PanicOnErr(err)
		twos[i], err = two.Insert().SetName(fmt.Sprintf("two %d", i)).Exec()
		examples.PanicOnErr(err)
	}

	// Add a relations

	// Add a relation to `One` to it's `One.Two` field is done by the auto-generated
	// function in One's orm `RemoveTwo`.
	err = one.RelationTwo().Add(ones[0].ID, twos[0].ID)
	examples.PanicOnErr(err)

	// Test this by a query:
	o, err := one.Select().
		JoinTwo(two.Select().Joiner()).
		Where(one.Where().ID(orm.OpEq, ones[0].ID)).
		First()
	examples.PanicOnErr(err)
	fmt.Printf("1. %+v\n", *o)

	// On the same way we could add Two's relation to One, in it's field Two.One:
	err = two.RelationOne().Add(twos[1].ID, ones[1].ID)
	examples.PanicOnErr(err)

	// Test this by a query:
	t, err := two.Select().
		JoinOne(one.Select().Joiner()).
		Where(two.Where().ID(orm.OpEq, twos[1].ID)).
		First()
	examples.PanicOnErr(err)
	fmt.Printf("2. %+v\n", *t)

	// If we will add to the mutual relation, we will see the relation added in both types.
	// This is because One and Two are sharing the same relation table, because they both
	// named it "mutual"
	err = one.Relationmutual().Add(ones[0].ID, twos[1].ID)
	examples.PanicOnErr(err)

	o, err = one.Select().
		JoinTwo(two.Select().Joiner()).
		JoinMutual(two.Select().Joiner()).
		Where(one.Where().ID(orm.OpEq, ones[0].ID)).
		First()
	examples.PanicOnErr(err)
	t, err = two.Select().
		JoinOne(one.Select().Joiner()).
		JoinMutual(one.Select().Joiner()).
		Where(two.Where().ID(orm.OpEq, twos[1].ID)).
		First()
	examples.PanicOnErr(err)
	fmt.Printf("3. %+v\n   %+v\n", *o, *t)

	// Remove a relation

	// The relation builder also has a `Remove` method, that removes a relation between two objects.
	err = one.RelationTwo().Remove(ones[0].ID, twos[0].ID)
	examples.PanicOnErr(err)

	o, err = one.Select().
		JoinTwo(two.Select().Joiner()).
		JoinMutual(two.Select().Joiner()).
		Where(one.Where().ID(orm.OpEq, ones[0].ID)).
		First()
	examples.PanicOnErr(err)
	fmt.Printf("4. %+v\n", *o)

	// Output:
	// 1. {ID:1 Name:one 0 Two:[{ID:1 Name:two 0 One:[] Mutual:[]}] Mutual:[]}
	// 2. {ID:2 Name:two 1 One:[{ID:2 Name:one 1 Two:[] Mutual:[]}] Mutual:[]}
	// 3. {ID:1 Name:one 0 Two:[{ID:1 Name:two 0 One:[] Mutual:[]}] Mutual:[{ID:2 Name:two 1 One:[] Mutual:[]}]}
	//    {ID:2 Name:two 1 One:[{ID:2 Name:one 1 Two:[] Mutual:[]}] Mutual:[{ID:1 Name:one 0 Two:[] Mutual:[]}]}
	// 4. {ID:1 Name:one 0 Two:[] Mutual:[{ID:2 Name:two 1 One:[] Mutual:[]}]}
}
