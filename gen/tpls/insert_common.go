package tpls

// NewInsert returns a new INSERT statement
func NewInsert() Insert {
	return Insert{}
}

type Insert struct {
	cols   []string
	values []interface{}
}

func (i Insert) add(name string, value interface{}) Insert {
	i.cols = append(i.cols, name)
	i.values = append(i.values, value)
	return i
}
