package {{.Package}}

// Select returns an object to create a SELECT statement
func (o *ORM) Select() *TSelect {
	return &TSelect{orm: o}
}
