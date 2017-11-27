package tpls

// Logger is a fmt.Printf - like function
type Logger func(string, ...interface{})

// log if a logger was set
func (o *orm) log(s string, args ...interface{}) {
	if o.logger == nil {
		return
	}
	o.logger(s, args...)
}
