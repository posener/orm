package tpls

// Query returns a struct to perform query operations
func Query() TQuery {
	return TQuery{}
}

// TQuery is a struct to hold query information
type TQuery struct {
	sel   TSelect
	where Where
}

// Select applies columns select on the query
func (q TQuery) Select(s TSelect) TQuery {
	q.sel = s
	return q
}

// Where applies where conditions on the query
func (q TQuery) Where(w Where) TQuery {
	q.where = w
	return q
}
