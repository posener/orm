package builder

import (
	"bytes"
	"fmt"
	"sync"
)

const initialBufferSize = 1024

var pool = &sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, initialBufferSize))
	},
}

// Dialect is dialect interface for a builder
type Dialect interface {
	Quote(string) string
	Var(int) string
}

// New returns a new SQL statement builder
func New(d Dialect, word string) *Builder {
	b := &Builder{
		dialect: d,
		buf:     pool.Get().(*bytes.Buffer),
	}
	b.buf.Reset()
	b.buf.WriteString(word)
	return b
}

// Builder is a builder for SQL statements
type Builder struct {
	dialect   Dialect
	buf       *bytes.Buffer
	args      []interface{}
	skipSpace bool
}

// Append a phrase to the statement
func (b *Builder) Append(word string) {
	if word != "" {
		if !b.skipSpace {
			b.buf.WriteByte(' ')
		}
		b.buf.WriteString(word)
		b.skipSpace = false
	}
}

// Quote appends the word quoted
func (b *Builder) Quote(word string) {
	b.Append(b.dialect.Quote(word))
}

// QuoteFullName quotes two names with a dot between them
func (b *Builder) QuoteFullName(surename, firstname string) {
	b.Append(fmt.Sprintf("%s.%s", b.dialect.Quote(surename), b.dialect.Quote(firstname)))
}

// Var adds a placeholder for a variable in the statement, and adds
// the variable to the list of arguments
func (b *Builder) Var(arg interface{}) {
	b.args = append(b.args, arg)
	b.Append(b.dialect.Var(len(b.args)))
}

// Open opens parentheses
func (b *Builder) Open() {
	b.buf.Write([]byte(" ("))
	b.skipSpace = true
}

// Close closes parentheses
func (b *Builder) Close() {
	b.buf.WriteByte(')')
	b.skipSpace = false
}

// Comma adds a comma
func (b *Builder) Comma() {
	b.buf.WriteByte(',')
	b.skipSpace = false
}

// Statement returns the statement, and return the used buffer
// so no longer calls to builder can be done afterwards
func (b *Builder) Statement() string {
	stmt := b.buf.String()
	pool.Put(b.buf)
	b.buf = nil
	return stmt
}

// Args returns the collected argument list
func (b *Builder) Args() []interface{} {
	return b.args
}
