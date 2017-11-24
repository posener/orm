package api

type Dialect interface {
	Name() string
	Create() string
}

type Tabler interface {
	Table() string
}

type Columner interface {
	Columns() []string
	Count() bool
}

type Order struct {
	Column, Dir string
}

type Orderer interface {
	Orders() []Order
}

type Group struct {
	Column string
}

type Grouper interface {
	Groups() []Group
}

type Wherer interface {
	Where() string
}

type Pager interface {
	Page() (limit int64, offset int64)
}

type Assignment struct {
	Column string
	Value  interface{}
}

type Assigner interface {
	Assignments() []Assignment
}
