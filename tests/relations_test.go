package tests

import (
	"fmt"
	"testing"

	"github.com/posener/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelationOneToOne(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		aORM, _, cORM := orms(t, conn)

		if conn.name != "sqlite3" { // this is not enforced in sqlite3
			_, err := aORM.Insert().InsertA(&A{Name: "James", CPointer: &C{ID: -10}}).Exec()
			require.NotNil(t, err)
		}

		c1, err := cORM.Insert().
			SetName("The Hitchhiker's Guide to the Galaxy").
			SetYear(1979).
			Exec()
		require.Nil(t, err)

		a1, err := aORM.Insert().InsertA(&A{Name: "James", Age: 42, CPointer: c1}).Exec()
		require.Nil(t, err)

		// query without join, A.CPointer should be nil
		aItems, err := aORM.Select().Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, a1.Name, aItems[0].Name)
			assert.Nil(t, aItems[0].CPointer)
		}

		// query with join, A.CPointer should be filled with aORM's properties
		aItems, err = aORM.Select().
			JoinCPointer(cORM.Select().Joiner()).
			Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, a1, &aItems[0])
		}

		// query with join, A.CPointer should be filled with aORM's properties
		aItems, err = aORM.Select().
			JoinCPointer(cORM.Select().SelectName().Joiner()).
			Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, a1.Name, aItems[0].Name)
			assert.Equal(t, a1.Age, aItems[0].Age)
			assert.Equal(t, a1.CPointer.Name, aItems[0].CPointer.Name)
			assert.Equal(t, 0, aItems[0].CPointer.Year) // was not selected in query, thus expecting zero value
		}

		// query with join, A.CPointer should be filled with aORM's properties
		aItems, err = aORM.Select().
			SelectName().
			JoinCPointer(cORM.Select().SelectYear().Joiner()).
			Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, a1.Name, aItems[0].Name)
			assert.Equal(t, 0, aItems[0].Age)            // was not selected in query, thus expecting zero value
			assert.Equal(t, "", aItems[0].CPointer.Name) // was not selected in query, thus expecting zero value
			assert.Equal(t, a1.CPointer.Year, aItems[0].CPointer.Year)
		}

		// test update of reference
		c2, err := cORM.Insert().
			SetName("1984").
			SetYear(1949).
			Exec()

		_, err = aORM.Update().
			Where(aORM.Where().ID(orm.OpEq, a1.ID)).
			SetCPointer(c2).
			Exec()
		require.Nil(t, err)

		gotA1, err := aORM.Select().
			JoinCPointer(cORM.Select().Joiner()).
			First()
		require.Nil(t, err)
		assert.Equal(t, a1.Name, gotA1.Name)
		assert.Equal(t, c2.Name, gotA1.CPointer.Name)
	})
}

func TestRelationOneToMany(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		_, bORM, cORM := orms(t, conn)

		b1, err := bORM.Insert().InsertB(&B{Name: "Marks", Hobbies: "drones"}).Exec()
		require.Nil(t, err)

		generateCs(t, cORM, b1, 10)

		bItems, err := bORM.Select().JoinCsPointer(cORM.Select().Joiner()).Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(bItems))
		assert.Equal(t, b1, &bItems[0])

		b2, err := bORM.Insert().InsertB(&B{Name: "Yoko", Hobbies: "music"}).Exec()
		require.Nil(t, err)

		// expect to get only one b, since the second b doesn't have any c related
		bItems, err = bORM.Select().JoinCsPointer(cORM.Select().Joiner()).Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(bItems))
		assert.Equal(t, b1, &bItems[0])

		generateCs(t, cORM, b2, 5)

		bItems, err = bORM.Select().JoinCsPointer(cORM.Select().Joiner()).Query()
		require.Nil(t, err)
		require.Equal(t, 2, len(bItems))
		assert.Equal(t, b1, &bItems[0])
		assert.Equal(t, b2, &bItems[1])

		bItems, err = bORM.Select().
			Where(bORM.Where().Name(orm.OpEq, "Yoko")).
			JoinCsPointer(
				cORM.Select().
					Where(cORM.Where().Year(orm.OpGt, 1996).And(cORM.Where().Year(orm.OpLt, 1999))).
					Joiner(),
			).Query()

		require.Nil(t, err)
		require.Equal(t, 1, len(bItems))
		b2.CsPointer = b2.CsPointer[2:4]
		assert.Equal(t, b2, &bItems[0])

		// test update of reference
		cORM.Update().
			Where(cORM.Where().Year(orm.OpEq, 2000)).
			SetB(b2).
			Exec()

		bItems, err = bORM.Select().
			Where(bORM.Where().ID(orm.OpEq, b1.ID)).
			JoinCsPointer(cORM.Select().Joiner()).
			Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(bItems))
		assert.Equal(t, 9, len(bItems[0].CsPointer))
	})
}

func generateCs(t *testing.T, cORM CORM, bItem *B, count int) {
	t.Helper()
	var cItems []C
	for i := 0; i < count; i++ {
		cItem, err := cORM.Insert().
			InsertC(&C{Name: fmt.Sprintf("Book %d", i), Year: 2000 - i, B: bItem}).
			Exec()
		require.Nil(t, err)
		cItem.B = nil
		cItems = append(cItems, *cItem)
		bItem.CsPointer = append(bItem.CsPointer, cItem)
	}
}

func TestRelationOneToOneNonPointerNested(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		a, err := NewA2ORM(conn.name, conn)
		require.Nil(t, err)
		b, err := NewB2ORM(conn.name, conn)
		require.Nil(t, err)
		c, err := NewC2ORM(conn.name, conn)
		require.Nil(t, err)
		d, err := NewD2ORM(conn.name, conn)
		require.Nil(t, err)
		if testing.Verbose() {
			a.Logger(t.Logf)
			b.Logger(t.Logf)
			c.Logger(t.Logf)
			d.Logger(t.Logf)
		}

		require.Nil(t, d.Create().Exec())
		require.Nil(t, c.Create().Exec())
		require.Nil(t, b.Create().Exec())
		require.Nil(t, a.Create().Exec())

		dItem, err := d.Insert().SetName("D").Exec()
		require.Nil(t, err)

		cItem, err := c.Insert().SetName("C").Exec()
		require.Nil(t, err)

		bItem, err := b.Insert().InsertB2(&B2{Name: "B", C: cItem, D: dItem}).Exec()
		require.Nil(t, err)

		aItem, err := a.Insert().InsertA2(&A2{B: *bItem}).Exec()
		require.Nil(t, err)

		// query without join, A.CPointer should be nil
		aItems, err := a.Select().Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, aItem.ID, aItems[0].ID)
			assert.Equal(t, int64(0), aItems[0].B.ID)
		}

		// test nested join
		aItems, err = a.Select().
			JoinB(b.Select().
				JoinC(c.Select().Joiner()).
				JoinD(d.Select().Joiner()).
				Joiner()).
			Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, aItem, &aItems[0])
		}

		// test nested join only one side
		aItems, err = a.Select().
			JoinB(b.Select().
				JoinC(c.Select().Joiner()).
				Joiner()).
			Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			aItem.B.D = nil
			assert.Equal(t, aItem, &aItems[0])
		}

		// test one level join
		aItems, err = a.Select().JoinB(b.Select().Joiner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			aItem.B.C = nil
			assert.Equal(t, aItem, &aItems[0])
		}
	})
}

func TestBidirectionalOneToManyRelationship(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		a, err := NewA3ORM(conn.name, conn)
		require.Nil(t, err)
		b, err := NewB3ORM(conn.name, conn)
		require.Nil(t, err)
		if testing.Verbose() {
			a.Logger(t.Logf)
			b.Logger(t.Logf)
		}

		require.Nil(t, a.Create().Exec())
		require.Nil(t, b.Create().Exec())

		aItem, err := a.Insert().SetName("A").Exec()
		require.Nil(t, err)

		for i := 0; i < 10; i++ {
			bItem, err := b.Insert().InsertB3(&B3{Name: fmt.Sprintf("B%d", i), A: aItem}).Exec()
			require.Nil(t, err)
			aItem.B = append(aItem.B, bItem)
			bItem.A = nil // for later query comparison
		}

		aList, err := a.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(aList))
		assert.Equal(t, "A", aList[0].Name)
		assert.Equal(t, 0, len(aList[0].B))

		aList, err = a.Select().JoinB(b.Select().Joiner()).Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(aList))
		assert.Equal(t, aItem, &aList[0])

		bList, err := b.Select().Query()
		require.Nil(t, err)
		if assert.Equal(t, 10, len(bList)) {
			assert.Equal(t, "B0", bList[0].Name)
			assert.Nil(t, bList[0].A)
		}

		bList, err = b.Select().JoinA(a.Select().Joiner()).Query()
		require.Nil(t, err)
		assert.Equal(t, 10, len(bList))
		for i, bItem := range bList {
			assert.Equal(t, fmt.Sprintf("B%d", i), bItem.Name)
			if assert.NotNil(t, bItem.A) {
				assert.Equal(t, "A", bItem.A.Name)
			}
		}
	})
}

func TestFieldsWithTheSameType(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		a, err := NewA4ORM(conn.name, conn)
		require.Nil(t, err)
		b, err := NewB4ORM(conn.name, conn)
		require.Nil(t, err)
		if testing.Verbose() {
			a.Logger(t.Logf)
			b.Logger(t.Logf)
		}

		require.Nil(t, b.Create().Exec())
		require.Nil(t, a.Create().Exec())

		b1, err := b.Insert().SetName("B1").Exec()
		require.Nil(t, err)

		b2, err := b.Insert().SetName("B2").Exec()
		require.Nil(t, err)

		_, err = a.Insert().InsertA4(&A4{Name: "A", B1: b1, B2: b2}).Exec()
		require.Nil(t, err)

		aList, err := a.Select().Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aList)) {
			assert.Equal(t, "A", aList[0].Name)
			assert.Nil(t, aList[0].B1)
			assert.Nil(t, aList[0].B2)
		}

		aList, err = a.Select().JoinB1(b.Select().Joiner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aList)) {
			assert.Equal(t, "A", aList[0].Name)
			assert.Equal(t, b1, aList[0].B1)
			assert.Nil(t, aList[0].B2)
		}

		aList, err = a.Select().JoinB1(b.Select().Joiner()).JoinB2(b.Select().Joiner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aList)) {
			assert.Equal(t, "A", aList[0].Name)
			assert.Equal(t, b1, aList[0].B1)
			assert.Equal(t, b2, aList[0].B2)
		}
	})
}

func TestSelfReferencing(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		a, err := NewA5ORM(conn.name, conn)
		require.Nil(t, err)
		if testing.Verbose() {
			a.Logger(t.Logf)
		}

		require.Nil(t, a.Create().Exec())

		// lets create this graph:
		//     a1
		//    /  \
		//   a2  a3
		//  /  \ / \
		// a4  a5  a6
		a6 := insertA5(t, a, "6", nil, nil)
		a5 := insertA5(t, a, "5", nil, nil)
		a4 := insertA5(t, a, "4", nil, nil)
		a3 := insertA5(t, a, "3", a5, a6)
		a2 := insertA5(t, a, "2", a4, a5)
		a1 := insertA5(t, a, "1", a2, a3)

		as, err := a.Select().Query()
		if assert.Equal(t, 6, len(as)) {
			for _, ai := range as {
				assert.Nil(t, ai.Left)
				assert.Nil(t, ai.Right)
			}
		}
		a1Got, err := a.Select().Where(a.Where().Name(orm.OpEq, "1")).First()
		require.Nil(t, err)
		assert.Equal(t, "1", a1Got.Name)
		assert.Nil(t, a1Got.Left)
		assert.Nil(t, a1Got.Right)

		a1Got, err = a.Select().
			Where(a.Where().Name(orm.OpEq, "1")).
			JoinLeft(a.Select().Joiner()).
			First()
		require.Nil(t, err)
		assert.Equal(t, "1", a1Got.Name)
		if assert.NotNil(t, a1Got.Left) {
			assert.Equal(t, "2", a1Got.Left.Name)
		}
		assert.Nil(t, a1Got.Right)

		joinLeftRight := a.Select().
			JoinLeft(a.Select().Joiner()).
			JoinRight(a.Select().Joiner()).
			Joiner()

		a1Got, err = a.Select().
			Where(a.Where().Name(orm.OpEq, "1")).
			JoinLeft(joinLeftRight).
			JoinRight(joinLeftRight).
			First()
		require.Nil(t, err)
		assert.Equal(t, a1, a1Got)
	})
}

func insertA5(t *testing.T, a A5ORM, name string, left, right *A5) *A5 {
	t.Helper()
	item, err := a.Insert().InsertA5(&A5{Name: name, Left: left, Right: right}).Exec()
	require.Nil(t, err)
	return item
}

func TestMultiplePrimaryKeys(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		if conn.name == "sqlite3" {
			t.Skip("sqlite3 does not support string type primary keys")
		}
		a, err := NewA6ORM(conn.name, conn)
		require.Nil(t, err)
		b, err := NewB6ORM(conn.name, conn)
		require.Nil(t, err)
		if testing.Verbose() {
			a.Logger(t.Logf)
			b.Logger(t.Logf)
		}

		require.Nil(t, b.Create().Exec())
		require.Nil(t, a.Create().Exec())

		b1, err := b.Insert().SetSureName("Jackson").SetFirstName("Michael").Exec()
		require.Nil(t, err)

		b2, err := b.Insert().SetSureName("Obama").SetFirstName("Barak").Exec()
		require.Nil(t, err)

		_, err = a.Insert().SetName("A").SetB(b1).Exec()
		require.Nil(t, err)

		_, err = a.Insert().SetName("B").SetB(b2).Exec()
		require.Nil(t, err)

		as, err := a.Select().Query()
		require.Nil(t, err)
		if assert.Equal(t, 2, len(as)) {
			assert.Equal(t, "A", as[0].Name)
			assert.Nil(t, as[0].B)
			assert.Equal(t, "B", as[1].Name)
			assert.Nil(t, as[1].B)
		}

		as, err = a.Select().JoinB(b.Select().Joiner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 2, len(as)) {
			assert.Equal(t, "A", as[0].Name)
			assert.Equal(t, b1, as[0].B)
			assert.Equal(t, "B", as[1].Name)
			assert.Equal(t, b2, as[1].B)
		}
	})
}

func TestReferencingField(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		a, err := NewA7ORM(conn.name, conn)
		require.Nil(t, err)
		b, err := NewB7ORM(conn.name, conn)
		require.Nil(t, err)
		if testing.Verbose() {
			a.Logger(t.Logf)
			b.Logger(t.Logf)
		}

		require.Nil(t, a.Create().Exec())
		require.Nil(t, b.Create().Exec())

		a1, err := a.Insert().SetName("A1").Exec()
		require.Nil(t, err)
		a2, err := a.Insert().SetName("A2").Exec()
		require.Nil(t, err)

		_, err = b.Insert().SetName("B1").SetA1(a1).SetA2(a2).Exec()
		require.Nil(t, err)
		_, err = b.Insert().SetName("B2").SetA1(a1).SetA2(a2).Exec()
		require.Nil(t, err)
		_, err = b.Insert().SetName("B3").SetA1(a2).SetA2(a1).Exec()
		require.Nil(t, err)
		_, err = b.Insert().SetName("B4").SetA1(a2).SetA2(a1).Exec()
		require.Nil(t, err)

		as, err := a.Select().JoinB(b.Select().Joiner()).Query()
		require.Nil(t, err)

		if assert.Equal(t, 2, len(as)) {
			if assert.Equal(t, "A1", as[0].Name) && assert.Equal(t, 2, len(as[0].B)) {
				assert.Equal(t, "B1", as[0].B[0].Name)
				assert.Equal(t, "B2", as[0].B[1].Name)
			}
			if assert.Equal(t, "A2", as[1].Name) && assert.Equal(t, 2, len(as[1].B)) {
				assert.Equal(t, "B3", as[1].B[0].Name)
				assert.Equal(t, "B4", as[1].B[1].Name)
			}
		}
	})
}

func orms(t *testing.T, conn conn) (AORM, BORM, CORM) {
	t.Helper()
	a, err := NewAORM(conn.name, conn)
	require.Nil(t, err)
	b, err := NewBORM(conn.name, conn)
	require.Nil(t, err)
	c, err := NewCORM(conn.name, conn)
	require.Nil(t, err)
	if testing.Verbose() {
		a.Logger(t.Logf)
		b.Logger(t.Logf)
		c.Logger(t.Logf)
	}

	require.Nil(t, b.Create().Exec())
	require.Nil(t, c.Create().Exec())
	require.Nil(t, a.Create().Exec())

	return a, b, c
}
