package {{.Package}}

func (c *TCreate) String() string {
    // Create statement has a line for each variable with it's name and it's type.
	return `{{.Type.CreateString}}`
}
