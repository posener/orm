// Autogenerated by github.com/posener/orm
package personorm

import (
	"fmt"

	"github.com/posener/orm/example"
)

// API is the interface of the ORM object
type API interface {
	Create() *TCreate
	Select() *TSelect
	Insert() *TInsert
	Update() *TUpdate
	Delete() *TDelete
	InsertPerson(*example.Person) *TInsert
	UpdatePerson(*example.Person) *TUpdate
}

// Querier is the interface for a SELECT SQL statement
type Querier interface {
	fmt.Stringer
	Query() ([]example.Person, error)
}
