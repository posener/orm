package dialect

import "bytes"

// newBuilder returns a new SQL statement builder
func newBuilder(d Dialect, word string) *builder {
	return &builder{Dialect: d, buf: bytes.NewBufferString(word)}
}

type builder struct {
	Dialect
	buf       *bytes.Buffer
	args      []interface{}
	skipSpace bool
}

// Append a phrase to the statement
func (b *builder) Append(word string) {
	if word != "" {
		if !b.skipSpace {
			b.buf.WriteByte(' ')
		}
		b.buf.WriteString(word)
		b.skipSpace = false
	}
}

// Var adds a placeholder for a variable in the statement, and adds
// the variable to the list of arguments
func (b *builder) Var(arg interface{}) {
	b.args = append(b.args, arg)
	b.Append(b.Dialect.Var(len(b.args)))
}

// Open opens parentheses
func (b *builder) Open() {
	b.buf.Write([]byte(" ("))
	b.skipSpace = true
}

// Close closes parentheses
func (b *builder) Close() {
	b.buf.WriteByte(')')
	b.skipSpace = false
}

// Comma adds a comma
func (b *builder) Comma() {
	b.buf.WriteByte(',')
	b.skipSpace = false
}

// Statement returns the built statement
func (b *builder) Statement() string {
	return b.buf.String()
}

// Args returns the list of arguments
func (b *builder) Args() []interface{} {
	return b.args
}
