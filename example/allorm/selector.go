// Package allorm was auto-generated by github.com/posener/orm; DO NOT EDIT
package allorm

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/posener/orm/example"

	"github.com/posener/orm/common"
)

const errMsg = "converting %s: column %d with value %v (type %T) to %s"

// selector selects columns for SQL queries and for parsing SQL rows
type selector struct {
	SelectAuto           bool
	SelectNotNil         bool
	SelectInt            bool
	SelectInt8           bool
	SelectInt16          bool
	SelectInt32          bool
	SelectInt64          bool
	SelectUInt           bool
	SelectUInt8          bool
	SelectUInt16         bool
	SelectUInt32         bool
	SelectUInt64         bool
	SelectTime           bool
	SelectVarCharString  bool
	SelectVarCharByte    bool
	SelectString         bool
	SelectBytes          bool
	SelectBool           bool
	SelectPInt           bool
	SelectPInt8          bool
	SelectPInt16         bool
	SelectPInt32         bool
	SelectPInt64         bool
	SelectPUInt          bool
	SelectPUInt8         bool
	SelectPUInt16        bool
	SelectPUInt32        bool
	SelectPUInt64        bool
	SelectPTime          bool
	SelectPVarCharString bool
	SelectPVarCharByte   bool
	SelectPString        bool
	SelectPBytes         bool
	SelectPBool          bool
	SelectSelect         bool
	count                bool // used for sql COUNT(*) column
}

// Columns are the names of selected columns
func (s *selector) Columns() []string {
	var cols []string
	if s.SelectAuto {
		cols = append(cols, "auto")
	}
	if s.SelectNotNil {
		cols = append(cols, "notnil")
	}
	if s.SelectInt {
		cols = append(cols, "int")
	}
	if s.SelectInt8 {
		cols = append(cols, "int8")
	}
	if s.SelectInt16 {
		cols = append(cols, "int16")
	}
	if s.SelectInt32 {
		cols = append(cols, "int32")
	}
	if s.SelectInt64 {
		cols = append(cols, "int64")
	}
	if s.SelectUInt {
		cols = append(cols, "uint")
	}
	if s.SelectUInt8 {
		cols = append(cols, "uint8")
	}
	if s.SelectUInt16 {
		cols = append(cols, "uint16")
	}
	if s.SelectUInt32 {
		cols = append(cols, "uint32")
	}
	if s.SelectUInt64 {
		cols = append(cols, "uint64")
	}
	if s.SelectTime {
		cols = append(cols, "time")
	}
	if s.SelectVarCharString {
		cols = append(cols, "varcharstring")
	}
	if s.SelectVarCharByte {
		cols = append(cols, "varcharbyte")
	}
	if s.SelectString {
		cols = append(cols, "string")
	}
	if s.SelectBytes {
		cols = append(cols, "bytes")
	}
	if s.SelectBool {
		cols = append(cols, "bool")
	}
	if s.SelectPInt {
		cols = append(cols, "pint")
	}
	if s.SelectPInt8 {
		cols = append(cols, "pint8")
	}
	if s.SelectPInt16 {
		cols = append(cols, "pint16")
	}
	if s.SelectPInt32 {
		cols = append(cols, "pint32")
	}
	if s.SelectPInt64 {
		cols = append(cols, "pint64")
	}
	if s.SelectPUInt {
		cols = append(cols, "puint")
	}
	if s.SelectPUInt8 {
		cols = append(cols, "puint8")
	}
	if s.SelectPUInt16 {
		cols = append(cols, "puint16")
	}
	if s.SelectPUInt32 {
		cols = append(cols, "puint32")
	}
	if s.SelectPUInt64 {
		cols = append(cols, "puint64")
	}
	if s.SelectPTime {
		cols = append(cols, "ptime")
	}
	if s.SelectPVarCharString {
		cols = append(cols, "pvarcharstring")
	}
	if s.SelectPVarCharByte {
		cols = append(cols, "pvarcharbyte")
	}
	if s.SelectPString {
		cols = append(cols, "pstring")
	}
	if s.SelectPBytes {
		cols = append(cols, "pbytes")
	}
	if s.SelectPBool {
		cols = append(cols, "pbool")
	}
	if s.SelectSelect {
		cols = append(cols, "select")
	}
	return cols
}

// Joins are join options of the query
func (s *selector) Joins() []common.JoinParams {
	var joins []common.JoinParams

	return joins
}

// Count is true when a COUNT(*) column should be added to the query
func (s *selector) Count() bool {
	return s.count
}

// FirstCount scans an SQL row to a AllCount struct
func (s *selector) FirstCount(dialect string, vals []driver.Value) (*AllCount, error) {
	switch dialect {
	case "mysql":
		return s.scanmysql(vals)

	case "sqlite3":
		return s.scansqlite3(vals)
	default:
		return nil, fmt.Errorf("unsupported dialect %s", dialect)
	}
}

// First scans an SQL row to a All struct
func (s *selector) First(dialect string, vals []driver.Value) (*example.All, error) {
	item, err := s.FirstCount(dialect, vals)
	if err != nil {
		return nil, err
	}
	return &item.All, nil
}

// scanmysql scans mysql row to a All struct
func (s *selector) scanmysql(vals []driver.Value) (*AllCount, error) {
	var (
		row AllCount
		all = s.selectAll()
		i   int
	)

	if all || s.SelectAuto {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int(parseInt(val))
				row.Auto = tmp
			case int64:
				tmp := int(val)
				row.Auto = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Auto", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectNotNil {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := string(val)
				row.NotNil = tmp
			default:
				return nil, fmt.Errorf(errMsg, "NotNil", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectInt {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int(parseInt(val))
				row.Int = tmp
			case int64:
				tmp := int(val)
				row.Int = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Int", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectInt8 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int8(parseInt(val))
				row.Int8 = tmp
			case int64:
				tmp := int8(val)
				row.Int8 = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Int8", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectInt16 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int16(parseInt(val))
				row.Int16 = tmp
			case int64:
				tmp := int16(val)
				row.Int16 = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Int16", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectInt32 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int32(parseInt(val))
				row.Int32 = tmp
			case int64:
				tmp := int32(val)
				row.Int32 = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Int32", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectInt64 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int64(parseInt(val))
				row.Int64 = tmp
			case int64:
				tmp := int64(val)
				row.Int64 = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Int64", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectUInt {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint(parseFloat(val))
				row.UInt = tmp
			case int64:
				tmp := uint(val)
				row.UInt = tmp
			default:
				return nil, fmt.Errorf(errMsg, "UInt", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectUInt8 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint8(parseFloat(val))
				row.UInt8 = tmp
			case int64:
				tmp := uint8(val)
				row.UInt8 = tmp
			default:
				return nil, fmt.Errorf(errMsg, "UInt8", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectUInt16 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint16(parseFloat(val))
				row.UInt16 = tmp
			case int64:
				tmp := uint16(val)
				row.UInt16 = tmp
			default:
				return nil, fmt.Errorf(errMsg, "UInt16", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectUInt32 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint32(parseFloat(val))
				row.UInt32 = tmp
			case int64:
				tmp := uint32(val)
				row.UInt32 = tmp
			default:
				return nil, fmt.Errorf(errMsg, "UInt32", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectUInt64 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint64(parseFloat(val))
				row.UInt64 = tmp
			case int64:
				tmp := uint64(val)
				row.UInt64 = tmp
			default:
				return nil, fmt.Errorf(errMsg, "UInt64", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectTime {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := parseTime(val, 3)
				row.Time = tmp
			case time.Time:
				tmp := time.Time(val)
				row.Time = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Time", i, vals[i], vals[i], "[]byte, time.Time")
			}
		}
		i++
	}

	if all || s.SelectVarCharString {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := string(val)
				row.VarCharString = tmp
			default:
				return nil, fmt.Errorf(errMsg, "VarCharString", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectVarCharByte {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := []byte(val)
				row.VarCharByte = tmp
			default:
				return nil, fmt.Errorf(errMsg, "VarCharByte", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectString {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := string(val)
				row.String = tmp
			default:
				return nil, fmt.Errorf(errMsg, "String", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectBytes {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := []byte(val)
				row.Bytes = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Bytes", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectBool {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := parseBool(val)
				row.Bool = tmp
			case bool:
				tmp := bool(val)
				row.Bool = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Bool", i, vals[i], vals[i], "[]byte, bool")
			}
		}
		i++
	}

	if all || s.SelectPInt {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int(parseInt(val))
				row.PInt = &tmp
			case int64:
				tmp := int(val)
				row.PInt = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PInt", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPInt8 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int8(parseInt(val))
				row.PInt8 = &tmp
			case int64:
				tmp := int8(val)
				row.PInt8 = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PInt8", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPInt16 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int16(parseInt(val))
				row.PInt16 = &tmp
			case int64:
				tmp := int16(val)
				row.PInt16 = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PInt16", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPInt32 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int32(parseInt(val))
				row.PInt32 = &tmp
			case int64:
				tmp := int32(val)
				row.PInt32 = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PInt32", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPInt64 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int64(parseInt(val))
				row.PInt64 = &tmp
			case int64:
				tmp := int64(val)
				row.PInt64 = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PInt64", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPUInt {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint(parseFloat(val))
				row.PUInt = &tmp
			case int64:
				tmp := uint(val)
				row.PUInt = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PUInt", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPUInt8 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint8(parseFloat(val))
				row.PUInt8 = &tmp
			case int64:
				tmp := uint8(val)
				row.PUInt8 = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PUInt8", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPUInt16 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint16(parseFloat(val))
				row.PUInt16 = &tmp
			case int64:
				tmp := uint16(val)
				row.PUInt16 = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PUInt16", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPUInt32 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint32(parseFloat(val))
				row.PUInt32 = &tmp
			case int64:
				tmp := uint32(val)
				row.PUInt32 = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PUInt32", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPUInt64 {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := uint64(parseFloat(val))
				row.PUInt64 = &tmp
			case int64:
				tmp := uint64(val)
				row.PUInt64 = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PUInt64", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if all || s.SelectPTime {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := parseTime(val, 3)
				row.PTime = &tmp
			case time.Time:
				tmp := time.Time(val)
				row.PTime = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PTime", i, vals[i], vals[i], "[]byte, time.Time")
			}
		}
		i++
	}

	if all || s.SelectPVarCharString {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := string(val)
				row.PVarCharString = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PVarCharString", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectPVarCharByte {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := []byte(val)
				row.PVarCharByte = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PVarCharByte", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectPString {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := string(val)
				row.PString = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PString", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectPBytes {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := []byte(val)
				row.PBytes = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PBytes", i, vals[i], vals[i], "[]byte, []byte")
			}
		}
		i++
	}

	if all || s.SelectPBool {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := parseBool(val)
				row.PBool = &tmp
			case bool:
				tmp := bool(val)
				row.PBool = &tmp
			default:
				return nil, fmt.Errorf(errMsg, "PBool", i, vals[i], vals[i], "[]byte, bool")
			}
		}
		i++
	}

	if all || s.SelectSelect {
		if vals[i] != nil {
			switch val := vals[i].(type) {
			case []byte:
				tmp := int(parseInt(val))
				row.Select = tmp
			case int64:
				tmp := int(val)
				row.Select = tmp
			default:
				return nil, fmt.Errorf(errMsg, "Select", i, vals[i], vals[i], "[]byte, int64")
			}
		}
		i++
	}

	if s.count {
		switch val := vals[i].(type) {
		case int64:
			row.Count = val
		case []byte:
			row.Count = parseInt(val)
		default:
			return nil, fmt.Errorf(errMsg, "COUNT(*)", i, vals[i], vals[i], "int64, []byte")
		}
		i++
	}

	return &row, nil
}

// scansqlite3 scans sqlite3 row to a All struct
func (s *selector) scansqlite3(vals []driver.Value) (*AllCount, error) {
	var (
		row AllCount
		all = s.selectAll()
		i   int
	)

	if all || s.SelectAuto {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Auto", i, vals[i], vals[i], "int")
			}
			tmp := int(val)
			row.Auto = tmp
		}
		i++
	}

	if all || s.SelectNotNil {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "NotNil", i, vals[i], vals[i], "string")
			}
			tmp := string(val)
			row.NotNil = tmp
		}
		i++
	}

	if all || s.SelectInt {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Int", i, vals[i], vals[i], "int")
			}
			tmp := int(val)
			row.Int = tmp
		}
		i++
	}

	if all || s.SelectInt8 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Int8", i, vals[i], vals[i], "int8")
			}
			tmp := int8(val)
			row.Int8 = tmp
		}
		i++
	}

	if all || s.SelectInt16 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Int16", i, vals[i], vals[i], "int16")
			}
			tmp := int16(val)
			row.Int16 = tmp
		}
		i++
	}

	if all || s.SelectInt32 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Int32", i, vals[i], vals[i], "int32")
			}
			tmp := int32(val)
			row.Int32 = tmp
		}
		i++
	}

	if all || s.SelectInt64 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Int64", i, vals[i], vals[i], "int64")
			}
			tmp := int64(val)
			row.Int64 = tmp
		}
		i++
	}

	if all || s.SelectUInt {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "UInt", i, vals[i], vals[i], "uint")
			}
			tmp := uint(val)
			row.UInt = tmp
		}
		i++
	}

	if all || s.SelectUInt8 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "UInt8", i, vals[i], vals[i], "uint8")
			}
			tmp := uint8(val)
			row.UInt8 = tmp
		}
		i++
	}

	if all || s.SelectUInt16 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "UInt16", i, vals[i], vals[i], "uint16")
			}
			tmp := uint16(val)
			row.UInt16 = tmp
		}
		i++
	}

	if all || s.SelectUInt32 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "UInt32", i, vals[i], vals[i], "uint32")
			}
			tmp := uint32(val)
			row.UInt32 = tmp
		}
		i++
	}

	if all || s.SelectUInt64 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "UInt64", i, vals[i], vals[i], "uint64")
			}
			tmp := uint64(val)
			row.UInt64 = tmp
		}
		i++
	}

	if all || s.SelectTime {
		if vals[i] != nil {
			val, ok := vals[i].(time.Time)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Time", i, vals[i], vals[i], "time.Time")
			}
			tmp := time.Time(val)
			row.Time = tmp
		}
		i++
	}

	if all || s.SelectVarCharString {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "VarCharString", i, vals[i], vals[i], "string")
			}
			tmp := string(val)
			row.VarCharString = tmp
		}
		i++
	}

	if all || s.SelectVarCharByte {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "VarCharByte", i, vals[i], vals[i], "[]byte")
			}
			tmp := []byte(val)
			row.VarCharByte = tmp
		}
		i++
	}

	if all || s.SelectString {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "String", i, vals[i], vals[i], "string")
			}
			tmp := string(val)
			row.String = tmp
		}
		i++
	}

	if all || s.SelectBytes {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Bytes", i, vals[i], vals[i], "[]byte")
			}
			tmp := []byte(val)
			row.Bytes = tmp
		}
		i++
	}

	if all || s.SelectBool {
		if vals[i] != nil {
			val, ok := vals[i].(bool)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Bool", i, vals[i], vals[i], "bool")
			}
			tmp := bool(val)
			row.Bool = tmp
		}
		i++
	}

	if all || s.SelectPInt {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PInt", i, vals[i], vals[i], "*int")
			}
			tmp := int(val)
			row.PInt = &tmp
		}
		i++
	}

	if all || s.SelectPInt8 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PInt8", i, vals[i], vals[i], "*int8")
			}
			tmp := int8(val)
			row.PInt8 = &tmp
		}
		i++
	}

	if all || s.SelectPInt16 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PInt16", i, vals[i], vals[i], "*int16")
			}
			tmp := int16(val)
			row.PInt16 = &tmp
		}
		i++
	}

	if all || s.SelectPInt32 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PInt32", i, vals[i], vals[i], "*int32")
			}
			tmp := int32(val)
			row.PInt32 = &tmp
		}
		i++
	}

	if all || s.SelectPInt64 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PInt64", i, vals[i], vals[i], "*int64")
			}
			tmp := int64(val)
			row.PInt64 = &tmp
		}
		i++
	}

	if all || s.SelectPUInt {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PUInt", i, vals[i], vals[i], "*uint")
			}
			tmp := uint(val)
			row.PUInt = &tmp
		}
		i++
	}

	if all || s.SelectPUInt8 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PUInt8", i, vals[i], vals[i], "*uint8")
			}
			tmp := uint8(val)
			row.PUInt8 = &tmp
		}
		i++
	}

	if all || s.SelectPUInt16 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PUInt16", i, vals[i], vals[i], "*uint16")
			}
			tmp := uint16(val)
			row.PUInt16 = &tmp
		}
		i++
	}

	if all || s.SelectPUInt32 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PUInt32", i, vals[i], vals[i], "*uint32")
			}
			tmp := uint32(val)
			row.PUInt32 = &tmp
		}
		i++
	}

	if all || s.SelectPUInt64 {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PUInt64", i, vals[i], vals[i], "*uint64")
			}
			tmp := uint64(val)
			row.PUInt64 = &tmp
		}
		i++
	}

	if all || s.SelectPTime {
		if vals[i] != nil {
			val, ok := vals[i].(time.Time)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PTime", i, vals[i], vals[i], "*time.Time")
			}
			tmp := time.Time(val)
			row.PTime = &tmp
		}
		i++
	}

	if all || s.SelectPVarCharString {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PVarCharString", i, vals[i], vals[i], "*string")
			}
			tmp := string(val)
			row.PVarCharString = &tmp
		}
		i++
	}

	if all || s.SelectPVarCharByte {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PVarCharByte", i, vals[i], vals[i], "*[]byte")
			}
			tmp := []byte(val)
			row.PVarCharByte = &tmp
		}
		i++
	}

	if all || s.SelectPString {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PString", i, vals[i], vals[i], "*string")
			}
			tmp := string(val)
			row.PString = &tmp
		}
		i++
	}

	if all || s.SelectPBytes {
		if vals[i] != nil {
			val, ok := vals[i].([]byte)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PBytes", i, vals[i], vals[i], "*[]byte")
			}
			tmp := []byte(val)
			row.PBytes = &tmp
		}
		i++
	}

	if all || s.SelectPBool {
		if vals[i] != nil {
			val, ok := vals[i].(bool)
			if !ok {
				return nil, fmt.Errorf(errMsg, "PBool", i, vals[i], vals[i], "*bool")
			}
			tmp := bool(val)
			row.PBool = &tmp
		}
		i++
	}

	if all || s.SelectSelect {
		if vals[i] != nil {
			val, ok := vals[i].(int64)
			if !ok {
				return nil, fmt.Errorf(errMsg, "Select", i, vals[i], vals[i], "int")
			}
			tmp := int(val)
			row.Select = tmp
		}
		i++
	}

	if s.count {
		switch val := vals[i].(type) {
		case int64:
			row.Count = val
		case []byte:
			row.Count = parseInt(val)
		default:
			return nil, fmt.Errorf(errMsg, "COUNT(*)", i, vals[i], vals[i], "int64, []byte")
		}
		i++
	}

	return &row, nil
}

// selectAll returns true if no column was specifically selected
func (s *selector) selectAll() bool {
	return !s.SelectAuto && !s.SelectNotNil && !s.SelectInt && !s.SelectInt8 && !s.SelectInt16 && !s.SelectInt32 && !s.SelectInt64 && !s.SelectUInt && !s.SelectUInt8 && !s.SelectUInt16 && !s.SelectUInt32 && !s.SelectUInt64 && !s.SelectTime && !s.SelectVarCharString && !s.SelectVarCharByte && !s.SelectString && !s.SelectBytes && !s.SelectBool && !s.SelectPInt && !s.SelectPInt8 && !s.SelectPInt16 && !s.SelectPInt32 && !s.SelectPInt64 && !s.SelectPUInt && !s.SelectPUInt8 && !s.SelectPUInt16 && !s.SelectPUInt32 && !s.SelectPUInt64 && !s.SelectPTime && !s.SelectPVarCharString && !s.SelectPVarCharByte && !s.SelectPString && !s.SelectPBytes && !s.SelectPBool && !s.SelectSelect && !s.count
}
