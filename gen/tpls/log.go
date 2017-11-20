package tpls

type Logger func(string, ...interface{})

func (o *ORM) log(s string, args ...interface{}) {
	if o.logger == nil {
		return
	}
	o.logger(s, args...)
}
