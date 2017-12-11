package example

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

		cItem := &C{Name: "The Hitchhiker's Guide to the Galaxy", Year: 1979}

		res, err := cORM.Insert().SetName(cItem.Name).SetYear(cItem.Year).Exec()
		require.Nil(t, err)
		cItem.ID, err = res.LastInsertId()
		require.Nil(t, err)

		aItem := &A{Name: "James", Age: 42, CPointer: cItem}

		res, err = aORM.Insert().InsertA(aItem).Exec()
		require.Nil(t, err)
		aItem.ID, err = res.LastInsertId()
		require.Nil(t, err)

		// query without join, A.CPointer should be nil
		aItems, err := aORM.Select().Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, aItem.Name, aItems[0].Name)
			assert.Nil(t, aItems[0].CPointer)
		}

		// query with join, A.CPointer should be filled with aORM's properties
		aItems, err = aORM.Select().JoinCPointer(cORM.Select().Scanner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, aItem, &aItems[0])
		}

		// query with join, A.CPointer should be filled with aORM's properties
		aItems, err = aORM.Select().JoinCPointer(cORM.Select().SelectName().Scanner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, aItem.Name, aItems[0].Name)
			assert.Equal(t, aItem.Age, aItems[0].Age)
			assert.Equal(t, aItem.CPointer.Name, aItems[0].CPointer.Name)
			assert.Equal(t, 0, aItems[0].CPointer.Year) // was not selected in query, thus expecting zero value
		}

		// query with join, A.CPointer should be filled with aORM's properties
		aItems, err = aORM.Select().SelectName().JoinCPointer(cORM.Select().SelectYear().Scanner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, aItem.Name, aItems[0].Name)
			assert.Equal(t, 0, aItems[0].Age)            // was not selected in query, thus expecting zero value
			assert.Equal(t, "", aItems[0].CPointer.Name) // was not selected in query, thus expecting zero value
			assert.Equal(t, aItem.CPointer.Year, aItems[0].CPointer.Year)
		}
	})
}

func TestRelationOneToMany(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		_, bORM, cORM := orms(t, conn)

		bItem := &B{Name: "Marks", Hobbies: "drones"}
		res, err := bORM.Insert().InsertB(bItem).Exec()
		require.Nil(t, err)
		bItem.ID, err = res.LastInsertId()
		require.Nil(t, err)

		generateCs(t, cORM, bItem, 10)

		bItems, err := bORM.Select().JoinCsPointer(cORM.Select().Scanner()).Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(bItems))
		assert.Equal(t, bItem, &bItems[0])

		bItem2 := &B{Name: "Yoko", Hobbies: "music"}
		res, err = bORM.Insert().InsertB(bItem2).Exec()
		require.Nil(t, err)
		bItem2.ID, err = res.LastInsertId()
		require.Nil(t, err)

		// expect to get only one b, since the second b doesn't have any c related
		bItems, err = bORM.Select().JoinCsPointer(cORM.Select().Scanner()).Query()
		require.Nil(t, err)
		require.Equal(t, 1, len(bItems))
		assert.Equal(t, bItem, &bItems[0])

		generateCs(t, cORM, bItem2, 5)

		bItems, err = bORM.Select().JoinCsPointer(cORM.Select().Scanner()).Query()
		require.Nil(t, err)
		require.Equal(t, 2, len(bItems))
		assert.Equal(t, bItem, &bItems[0])
		assert.Equal(t, bItem2, &bItems[1])

		bItems, err = bORM.Select().
			Where(bORM.Where().Name(orm.OpEq, "Yoko")).
			JoinCsPointer(
				cORM.Select().
					Where(cORM.Where().Year(orm.OpGt, 1996).And(cORM.Where().Year(orm.OpLt, 1999))).
					Scanner(),
			).Query()

		require.Nil(t, err)
		require.Equal(t, 1, len(bItems))
		bItem2.CsPointer = bItem2.CsPointer[2:4]
		assert.Equal(t, bItem2, &bItems[0])
	})
}

func generateCs(t *testing.T, cORM CORM, bItem *B, count int) {
	t.Helper()
	var cItems []C
	for i := 0; i < count; i++ {
		cItem := &C{Name: fmt.Sprintf("Book %d", i), Year: 2000 - i, BID: bItem.ID}
		res, err := cORM.Insert().InsertC(cItem).Exec()
		require.Nil(t, err)
		cItem.ID, err = res.LastInsertId()
		require.Nil(t, err)
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

		_, err = d.Create().Exec()
		require.Nil(t, err)
		_, err = c.Create().Exec()
		require.Nil(t, err)
		_, err = b.Create().Exec()
		require.Nil(t, err)
		_, err = a.Create().Exec()
		require.Nil(t, err)

		dItem := &D2{Name: "D"}
		res, err := d.Insert().InsertD2(dItem).Exec()
		require.Nil(t, err)
		dItem.ID, err = res.LastInsertId()
		require.Nil(t, err)

		cItem := &C2{Name: "C"}
		res, err = c.Insert().InsertC2(cItem).Exec()
		require.Nil(t, err)
		cItem.ID, err = res.LastInsertId()
		require.Nil(t, err)

		bItem := &B2{Name: "B", C: cItem, D: dItem}
		res, err = b.Insert().InsertB2(bItem).Exec()
		require.Nil(t, err)
		bItem.ID, err = res.LastInsertId()
		require.Nil(t, err)

		aItem := &A2{B: *bItem}

		res, err = a.Insert().InsertA2(aItem).Exec()
		require.Nil(t, err)
		aItem.ID, err = res.LastInsertId()
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
				JoinC(c.Select().Scanner()).
				JoinD(d.Select().Scanner()).
				Scanner()).
			Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, aItem, &aItems[0])
		}

		// test nested join only one side
		aItems, err = a.Select().
			JoinB(b.Select().
				JoinC(c.Select().Scanner()).
				Scanner()).
			Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			aItem.B.D = nil
			assert.Equal(t, aItem, &aItems[0])
		}

		// test one level join
		aItems, err = a.Select().JoinB(b.Select().Scanner()).Query()
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

		_, err = a.Create().Exec()
		require.Nil(t, err)
		_, err = b.Create().Exec()
		require.Nil(t, err)

		aItem := &A3{Name: "A"}
		res, err := a.Insert().InsertA3(aItem).Exec()
		require.Nil(t, err)
		aItem.ID, err = res.LastInsertId()
		require.Nil(t, err)

		for i := 0; i < 10; i++ {
			bItem := &B3{Name: fmt.Sprintf("B%d", i), A: aItem}
			res, err := b.Insert().InsertB3(bItem).Exec()
			require.Nil(t, err)
			bItem.ID, err = res.LastInsertId()
			require.Nil(t, err)
			aItem.B = append(aItem.B, bItem)
			bItem.A = nil // for later query comparison
		}

		aList, err := a.Select().Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(aList))
		assert.Equal(t, "A", aList[0].Name)
		assert.Equal(t, 0, len(aList[0].B))

		aList, err = a.Select().JoinB(b.Select().Scanner()).Query()
		require.Nil(t, err)
		assert.Equal(t, 1, len(aList))
		assert.Equal(t, aItem, &aList[0])

		bList, err := b.Select().Query()
		require.Nil(t, err)
		if assert.Equal(t, 10, len(bList)) {
			assert.Equal(t, "B0", bList[0].Name)
			assert.Nil(t, bList[0].A)
		}

		bList, err = b.Select().JoinA(a.Select().Scanner()).Query()
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

		_, err = b.Create().Exec()
		require.Nil(t, err)
		_, err = a.Create().Exec()
		require.Nil(t, err)

		b1 := &B4{Name: "B1"}
		res, err := b.Insert().InsertB4(b1).Exec()
		require.Nil(t, err)
		b1.ID, err = res.LastInsertId()
		require.Nil(t, err)

		b2 := &B4{Name: "B2"}
		res, err = b.Insert().InsertB4(b2).Exec()
		require.Nil(t, err)
		b2.ID, err = res.LastInsertId()
		require.Nil(t, err)

		aItem := &A4{Name: "A", B1: b1, B2: b2}
		res, err = a.Insert().InsertA4(aItem).Exec()
		require.Nil(t, err)
		aItem.ID, err = res.LastInsertId()
		require.Nil(t, err)

		aList, err := a.Select().Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aList)) {
			assert.Equal(t, "A", aList[0].Name)
			assert.Nil(t, aList[0].B1)
			assert.Nil(t, aList[0].B2)
		}

		aList, err = a.Select().JoinB1(b.Select().Scanner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aList)) {
			assert.Equal(t, "A", aList[0].Name)
			assert.Equal(t, b1, aList[0].B1)
			assert.Nil(t, aList[0].B2)
		}

		aList, err = a.Select().JoinB1(b.Select().Scanner()).JoinB2(b.Select().Scanner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aList)) {
			assert.Equal(t, "A", aList[0].Name)
			assert.Equal(t, b1, aList[0].B1)
			assert.Equal(t, b2, aList[0].B2)
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

	_, err = b.Create().Exec()
	require.Nil(t, err)

	_, err = c.Create().Exec()
	require.Nil(t, err)

	_, err = a.Create().Exec()
	require.Nil(t, err)

	return a, b, c
}
