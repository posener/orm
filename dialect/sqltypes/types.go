package sqltypes

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var typeFormat = regexp.MustCompile(`([^(]+)(\((\d*)\))?`)

// Type represents an SQL column type
type Type struct {
	Name string
	Size int
}

// New parses a string into an sql type
func New(s string) (*Type, error) {
	t := new(Type)
	m := typeFormat.FindStringSubmatch(s)
	switch len(m) {
	case 0:
		return nil, fmt.Errorf("invalid SQL type: %s", s)
	case 1:
		t.Name = m[0]
	case 4:
		t.Name = m[1]
		t.Size, _ = strconv.Atoi(m[3])
	}
	t.Name = strings.ToLower(t.Name)
	return t, nil
}

// String returns the type as a string
func (t *Type) String() string {
	if t == nil {
		return ""
	}
	if t.Size == 0 {
		return t.Name
	}
	return fmt.Sprintf("%s(%d)", t.Name, t.Size)
}
