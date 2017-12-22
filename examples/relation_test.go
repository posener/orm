package examples

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/posener/orm"
)

// ExampleRelationship demonstrates relationships between structs with ORM
//
// In this example we use the following structs:
//
//type One struct {
//	ID   int64 `sql:"primary key;auto increment"`
//	Name string
//	// ToOne is one-to-one relationship with a type called "OtherOne"
//	ToOne *OtherOne
//	// ToMany is one-to-many relationship with a type called "OtherMany"
//	// In order for this relationship to exists, OtherMany must have a
//	// field that references a "One" type
//	ToMany []OtherMany
//}
//
//type OtherOne struct {
//	ID   int64 `sql:"primary key;auto increment"`
//	Name string
//}
//
//type OtherMany struct {
//	ID   int64 `sql:"primary key;auto increment"`
//	Name string
//	// MyOneIs is a filed that must exists for allowing a one to many relationship
//	// between a One type and OtherMany type.
//	MyOneIs *One
//}
func ExampleRelationship() {
	db := conn()

	oneORM, err := NewOneORM("mysql", db)
	panicOnErr(err)
	otherOneORM, err := NewOtherOneORM("mysql", db)
	panicOnErr(err)
	otherManyORM, err := NewOtherManyORM("mysql", db)
	panicOnErr(err)

	oneORM.Logger(log.Printf)
	otherOneORM.Logger(log.Printf)
	otherManyORM.Logger(log.Printf)

	// Create the tables: order is important!
	// When have a forward reference (one-to-one relationship) to another table, the other
	// table should be created before.
	// - One have a reference to OtherOne, thus OtherOne should be created before One.
	// - OtherMany have a reference to One, thus One should be created before OtherMany.
	// Notice that we are not talking about a field with a slice of other type, which
	// is a one-to-many relationship.
	// Creating the tables in another order will return an error: "a foreign key constraint fails"
	panicOnErr(otherOneORM.Create().Exec())
	panicOnErr(oneORM.Create().Exec())
	panicOnErr(otherManyORM.Create().Exec())

	// create other-one rows, rows that have one-to-one relationship with One
	oo1, err := otherOneORM.Insert().SetName("oo1").Exec()
	panicOnErr(err)
	oo2, err := otherOneORM.Insert().SetName("oo2").Exec()
	panicOnErr(err)
	oo3, err := otherOneORM.Insert().SetName("oo3").Exec()
	panicOnErr(err)

	// create one rows, relate to other-one rows
	// to add the relationship, just pass an object with the right primary keys filled in,
	// in our case, we can just set oo1, oo2 and oo3
	o1, err := oneORM.Insert().SetName("o1").SetOtherOne(oo1).Exec()
	panicOnErr(err)
	o2, err := oneORM.Insert().SetName("o2").SetOtherOne(oo2).Exec()
	panicOnErr(err)
	o3, err := oneORM.Insert().SetName("o3").SetOtherOne(oo3).Exec()
	panicOnErr(err)

	// Now, lets add other-many objects, 5 for each one object
	// To set the relationship, we should just use the `SetMyOneIs` function with the right
	// reference to One object.
	for i := 0; i < 5; i++ {
		_, err = otherManyORM.Insert().SetName(fmt.Sprintf("om1%d", i)).SetMyOne(o1).Exec()
		panicOnErr(err)
		_, err = otherManyORM.Insert().SetName(fmt.Sprintf("om2%d", i)).SetMyOne(o2).Exec()
		panicOnErr(err)
		_, err = otherManyORM.Insert().SetName(fmt.Sprintf("om3%d", i)).SetMyOne(o3).Exec()
		panicOnErr(err)
	}

	// Now that we have the data-set in our table, let's see how we query it
	// We can query One object:
	ones, err := oneORM.Select().Query()
	panicOnErr(err)

	// this simple query, will not return any relationships of the ones object
	// we got 3 objects:
	fmt.Println("1. ones len:", len(ones))
	// and all of their references are empty:
	fmt.Println("2. one's references:", ones[0].OtherOne, ones[0].OtherMany)

	// In order to get also referenced items, we must use an SQL join query, which,
	// for reasons of efficiency, must be given explicitly.
	// The select statement builder have a a join function, for each of the struct
	// relationships:

	ones, err = oneORM.Select().
		JoinOtherOne(otherOneORM.Select().Joiner()).
		JoinOtherMany(otherManyORM.Select().Joiner()).
		Query()

	var otherMany []string
	for _, om := range ones[0].OtherMany {
		otherMany = append(otherMany, om.Name)
	}
	fmt.Println("3. one with join's references:", ones[0].OtherOne.Name, otherMany)

	// All the select operations: Where, GroupBy, OrderBy, Page, and so, can be give to the
	// joined selector as well, and joins can be also applied recursively

	ones, err = oneORM.Select().
		// choose only o3 from all the "ones"
		Where(oneORM.Where().Name(orm.OpEq, "o3")).
		// Filter One.OtherMany to have only OtherMany.ID < 10
		// Nested join, join One from OtherMany
		JoinOtherMany(otherManyORM.Select().
			Where(otherManyORM.Where().ID(orm.OpLt, 10)).
			JoinMyOne(oneORM.Select().Joiner()).
			Joiner()).
		Query()

	// we expect to have only one entry:
	fmt.Println("4. complex join len:", len(ones))
	fmt.Println("5. complex join element:", ones[0].Name)

	// it shouldn't have OtherOne, since we didn't joined it
	fmt.Println("6. complex join other one:", ones[0].OtherOne)

	// it's OtherMany should be of length 3, since only 3 elements have their id less than 10
	fmt.Println("7. complex join other many len:", len(ones[0].OtherMany))

	// the OtherMany has a nested join, so it's MyOne reference shouldn't be nil
	fmt.Println("8. complex join other many's my one:", ones[0].OtherMany[0].MyOne.Name)

	// Output:
	// 1. ones len: 3
	// 2. one's references: <nil> []
	// 3. one with join's references: oo1 [om10 om11 om12 om13 om14]
	// 4. complex join len: 1
	// 5. complex join element: o3
	// 6. complex join other one: <nil>
	// 7. complex join other many len: 3
	// 8. complex join other many's my one: o3
}
