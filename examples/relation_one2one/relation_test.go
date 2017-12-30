package examples

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/posener/orm"
	"github.com/posener/orm/examples"
)

// ExampleRelationOne2One demonstrates one-to-one relationship between structs with ORM
//
// In this example we use the following structs:
//
//type One struct {
//	ID   int64 `sql:"primary key;auto increment"`
//	Name string
//	// ToOne is one-to-one relationship with a type called "OtherOne"
//	ToOne *OtherOne
//}
//
//type OtherOne struct {
//	ID   int64 `sql:"primary key;auto increment"`
//	Name string
//}
func ExampleRelationOne2One() {
	conn := examples.Conn("one2one")
	if conn == nil {
		return // mysql address was not defined
	}
	defer conn.Close()

	conn.Logger(log.Printf)

	oneORM, err := NewOneORM(conn)
	examples.PanicOnErr(err)
	otherOneORM, err := NewOtherOneORM(conn)
	examples.PanicOnErr(err)

	// Create the tables: order is important!
	// When have a forward reference (one-to-one relationship) to another table, the other
	// table should be created before.
	// One have a reference to OtherOne, thus OtherOne should be created before One.
	// Creating the tables in another order will return an error: "a foreign key constraint fails"
	examples.PanicOnErr(otherOneORM.Create().Exec())
	examples.PanicOnErr(oneORM.Create().Exec())

	// create other-one rows, rows that have one-to-one relationship with One
	oo1, err := otherOneORM.Insert().SetName("oo1").Exec()
	examples.PanicOnErr(err)
	oo2, err := otherOneORM.Insert().SetName("oo2").Exec()
	examples.PanicOnErr(err)
	oo3, err := otherOneORM.Insert().SetName("oo3").Exec()
	examples.PanicOnErr(err)

	// create one rows, relate to other-one rows
	// to add the relationship, just pass an object with the right primary keys filled in,
	// in our case, we can just set oo1, oo2 and oo3
	_, err = oneORM.Insert().SetName("o1").SetOtherOne(oo1).Exec()
	examples.PanicOnErr(err)
	_, err = oneORM.Insert().SetName("o2").SetOtherOne(oo2).Exec()
	examples.PanicOnErr(err)
	_, err = oneORM.Insert().SetName("o3").SetOtherOne(oo3).Exec()
	examples.PanicOnErr(err)

	// Now that we have the data-set in our table, let's see how we query it
	// We can query One object:
	ones, err := oneORM.Select().Query()
	examples.PanicOnErr(err)

	// this simple query, will not return any relationships of the ones object
	// we got 3 objects:
	fmt.Println("1. ones len:", len(ones))
	// and all of their references are empty:
	fmt.Println("2. one's references:", ones[0].OtherOne)

	// In order to get also referenced items, we must use an SQL join query, which,
	// for reasons of efficiency, must be given explicitly.
	// The select statement builder have a a join function, for each of the struct
	// relationships:

	ones, err = oneORM.Select().
		JoinOtherOne(otherOneORM.Select().Joiner()).
		Query()

	fmt.Println("3. one with join's references:", ones[0].OtherOne.Name)

	// All the select operations: Where, GroupBy, OrderBy, Page, and so, can be give to the
	// joined selector as well, and joins can be also applied recursively
	ones, err = oneORM.Select().
		JoinOtherOne(otherOneORM.Select().
			Where(otherOneORM.Where().ID(orm.OpLt, 3)).Joiner()).
		Query()

	// we expect to have only 2 entries:
	fmt.Println("4. complex join len:", len(ones))
	fmt.Println("5. complex join element:", ones[0].Name)
	fmt.Println("6. complex join other one:", ones[0].OtherOne)

	// Output:
	// 1. ones len: 3
	// 2. one's references: <nil>
	// 3. one with join's references: oo1
	// 4. complex join len: 2
	// 5. complex join element: o1
	// 6. complex join other one: &{1 oo1}
}
