package example

import (
	"fmt"
	"testing"

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

func TestRelationOneToOneNonPointer(t *testing.T) {
	testDBs(t, func(t *testing.T, conn conn) {
		a, err := NewA2ORM(conn.name, conn)
		require.Nil(t, err)
		b, err := NewB2ORM(conn.name, conn)
		require.Nil(t, err)
		if testing.Verbose() {
			a.Logger(t.Logf)
			b.Logger(t.Logf)
		}

		_, err = b.Create().Exec()
		require.Nil(t, err)
		_, err = a.Create().Exec()
		require.Nil(t, err)

		bItem := &B2{Name: "The Hitchhiker's Guide to the Galaxy"}

		res, err := b.Insert().InsertB2(bItem).Exec()
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

		aItems, err = a.Select().JoinB(b.Select().Scanner()).Query()
		require.Nil(t, err)
		if assert.Equal(t, 1, len(aItems)) {
			assert.Equal(t, aItem, &aItems[0])
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
