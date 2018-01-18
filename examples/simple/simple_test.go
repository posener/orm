package simple

import (
	"context"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
	"github.com/posener/orm/examples"
)

// ExampleSimpleUsage
// This example shows simple usage with a simple type "Simple":
//type Simple struct {
//  ID     int64  `sql:"primary key;auto increment"`
//	Field1 string
//	Field2 int
//	Field3 bool
//}
func ExampleSimpleUsage() {
	conn, err := orm.Open("sqlite3", ":memory:", orm.OptLogger(log.Printf))
	examples.PanicOnErr(err)
	defer conn.Close()

	// Create an ORM controller for the Simple struct
	sorm, err := NewSimpleORM(conn)
	examples.PanicOnErr(err)

	// Create will create a new table
	// IfNotExists will add an SQL "IF NOT EXISTS" statement, so that the
	// CREATE statement won't fail in case that the table already exists.
	// Invoking Exec() causes the create statement to be executed.
	err = sorm.Create().IfNotExists().Exec()
	examples.PanicOnErr(err)

	// Insert an item with a Simple object.
	// Invoking Exec() causes the insert statement to be executed.
	// We get the type in the response, with it's ID filled in
	simple1, err := sorm.Insert().
		InsertSimple(&Simple{Field1: "simple 1", Field2: 12, Field3: true}).
		Exec()
	examples.PanicOnErr(err)

	// The first object in the table gets id 1, since we chosen auto increment for
	// this column.
	fmt.Println("simple1.ID:", simple1.ID)

	// Insert an object to the table using fields instead of an instance of the struct
	// We can also use context in all the ORM executions, if we want, using the
	// `Context` function.
	// Notice that in this execution, all the arguments for "SetField1" and "SetField2"
	// are typed according to the field type, and not an interface{}
	simple2, err := sorm.Insert().
		SetField1("simple 2").
		SetField2(12).
		Context(context.Background()).
		Exec()
	examples.PanicOnErr(err)

	// This time we expect the ID to be 2
	fmt.Println("simple2.ID:", simple2.ID)

	// Lets get an object from the database.
	// If we want a specific object, we can use the Get function, and give as argument
	// the object's primary key or keys (in case of multiple primary keys)
	getSimple2, err := sorm.Get(2).Exec()
	examples.PanicOnErr(err)

	// This time we expect the ID to be 2
	fmt.Println("getSimple2.ID:", getSimple2.ID)

	// If we try and get an object that the key is not in the table, we will get a not found error
	_, err = sorm.Get(3).Exec()
	fmt.Println("Error:", err.Error())

	// Listing objects is done with the Select function
	simples, err := sorm.Select().Query()
	examples.PanicOnErr(err)

	// We expect the select length to be 2 since we inserted 2 rows
	fmt.Println("Select len:", len(simples))

	// We can select with a where condition:
	// lets filter by field3 == true, notice that the sorm.Where().Field3() function
	// gets two typed arguments, no strings, no interface{}
	// The where statement can have OR and AND operator for applying several conditions
	simples, err = sorm.Select().
		Where(sorm.Where().Field3(orm.OpEq, true)).
		Query()
	examples.PanicOnErr(err)

	// We expect the select length to be 1 since only the 1st row agrees with the
	// where conditions
	fmt.Println("Select where len:", len(simples))

	// We can group by and order by. Let's try:
	simples, err = sorm.Select().
		GroupBy(SimpleColField2).
		Query()
	examples.PanicOnErr(err)

	// We expect the select length to be 1 since both rows have the same value in field2
	fmt.Println("Select group len:", len(simples))

	// Update row 1 and set field1 and field2
	_, err = sorm.Update().
		Where(sorm.Where().ID(orm.OpEq, 1)).
		SetField1("updated").
		SetField2(1).
		Exec()
	examples.PanicOnErr(err)

	// Check the new values:
	simple1Updated, err := sorm.Get(1).Exec()
	examples.PanicOnErr(err)

	fmt.Println("Updated:", simple1Updated.Field1, simple1Updated.Field2, simple1Updated.Field3)

	// Delete a row:
	_, err = sorm.Delete().Where(sorm.Where().Field1(orm.OpEq, "updated")).Exec()
	examples.PanicOnErr(err)

	_, err = sorm.Get(1).Exec()
	fmt.Println("Should not be found:", err.Error())

	// Output: simple1.ID: 1
	// simple2.ID: 2
	// getSimple2.ID: 2
	// Error: not found
	// Select len: 2
	// Select where len: 1
	// Select group len: 1
	// Updated: updated 1 true
	// Should not be found: not found
}
