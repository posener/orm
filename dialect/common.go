package dialect

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const errConvert = "converting %s: column %d with value %v (type %T) to %s"

// ErrConvert returns a conversion error when processing sql values
func ErrConvert(field string, index int, value interface{}, to string) error {
	return fmt.Errorf(errConvert, field, index, value, value, to)
}

// Values is a hack to the sql.Rows struct
// since the rows struct does not expose it's lastcols values, or a way to give
// a custom scanner to the Scan method.
// See issue https://github.com/golang/go/issues/22544
func Values(r sql.Rows) []driver.Value {
	// some ugly hack to access lastcols field
	rs := reflect.ValueOf(&r).Elem()
	rf := rs.FieldByName("lastcols")

	// overcome panic reflect.ColumnValue.Interface: cannot return value obtained from unexported field or method
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	return rf.Interface().([]driver.Value)
}

// ParseInt parses a byte slice to int64
func ParseInt(s []byte) int64 {
	i, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		log.Printf("Failed parsing %s to int", string(s))
	}
	return i
}

// ParseUInt parses a byte slice to uint64
func ParseUInt(s []byte) uint64 {
	i, err := strconv.ParseUint(string(s), 10, 64)
	if err != nil {
		log.Printf("Failed parsing %s to uint", string(s))
	}
	return i
}

// ParseFloat parses a byte slice to float64
func ParseFloat(s []byte) float64 {
	i, err := strconv.ParseFloat(string(s), 64)
	if err != nil {
		log.Printf("Failed parsing %s to float", string(s))
	}
	return i
}

// ParseTime parses time from mysql database
// zero time is a special case, since it's month is 00, which causes month out of range error.
func ParseTime(b []byte, precision int) time.Time {
	var (
		format = "2006-01-02 15:04:05"
		zero   = "0000-00-00 00:00:00"
		s      = string(b)
	)
	if precision > 0 {
		ext := "." + strings.Repeat("0", precision)
		format += ext
		zero += ext
	}
	if s == zero {
		return time.Time{}
	}
	t, err := time.Parse(format, s)
	if err != nil {
		log.Printf("Failed parsing '%s' to time.Time with format '%s'", string(s), format)
	}
	return t
}

// ParseBool parses a byte slice to bool
func ParseBool(s []byte) bool {
	return s[0] != 0
}

// ContextOrBackground returns background context if the given context is nil,
// otherwise it returns the context itself.
func ContextOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
