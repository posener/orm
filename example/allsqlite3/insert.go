// Autogenerated by github.com/posener/orm; DO NOT EDIT
package allsqlite3

import (
	"time"

	"github.com/posener/orm/example"
)

// InsertAll creates an INSERT statement according to the given object
func (o *ORM) InsertAll(p *example.All) *Insert {
	i := o.Insert()
	i.internal.Assignments.Add("auto", p.Auto)
	i.internal.Assignments.Add("notnil", p.NotNil)
	i.internal.Assignments.Add("int", p.Int)
	i.internal.Assignments.Add("int8", p.Int8)
	i.internal.Assignments.Add("int16", p.Int16)
	i.internal.Assignments.Add("int32", p.Int32)
	i.internal.Assignments.Add("int64", p.Int64)
	i.internal.Assignments.Add("uint", p.UInt)
	i.internal.Assignments.Add("uint8", p.UInt8)
	i.internal.Assignments.Add("uint16", p.UInt16)
	i.internal.Assignments.Add("uint32", p.UInt32)
	i.internal.Assignments.Add("uint64", p.UInt64)
	i.internal.Assignments.Add("time", p.Time)
	i.internal.Assignments.Add("varcharstring", p.VarCharString)
	i.internal.Assignments.Add("varcharbyte", p.VarCharByte)
	i.internal.Assignments.Add("string", p.String)
	i.internal.Assignments.Add("bytes", p.Bytes)
	i.internal.Assignments.Add("bool", p.Bool)
	i.internal.Assignments.Add("pint", p.PInt)
	i.internal.Assignments.Add("pint8", p.PInt8)
	i.internal.Assignments.Add("pint16", p.PInt16)
	i.internal.Assignments.Add("pint32", p.PInt32)
	i.internal.Assignments.Add("pint64", p.PInt64)
	i.internal.Assignments.Add("puint", p.PUInt)
	i.internal.Assignments.Add("puint8", p.PUInt8)
	i.internal.Assignments.Add("puint16", p.PUInt16)
	i.internal.Assignments.Add("puint32", p.PUInt32)
	i.internal.Assignments.Add("puint64", p.PUInt64)
	i.internal.Assignments.Add("ptime", p.PTime)
	i.internal.Assignments.Add("pvarcharstring", p.PVarCharString)
	i.internal.Assignments.Add("pvarcharbyte", p.PVarCharByte)
	i.internal.Assignments.Add("pstring", p.PString)
	i.internal.Assignments.Add("pbytes", p.PBytes)
	i.internal.Assignments.Add("pbool", p.PBool)
	i.internal.Assignments.Add("select", p.Select)
	return i
}

// SetAuto sets value for column auto in the INSERT statement
func (i *Insert) SetAuto(value int) *Insert {
	i.internal.Assignments.Add("auto", value)
	return i
}

// SetNotNil sets value for column notnil in the INSERT statement
func (i *Insert) SetNotNil(value string) *Insert {
	i.internal.Assignments.Add("notnil", value)
	return i
}

// SetInt sets value for column int in the INSERT statement
func (i *Insert) SetInt(value int) *Insert {
	i.internal.Assignments.Add("int", value)
	return i
}

// SetInt8 sets value for column int8 in the INSERT statement
func (i *Insert) SetInt8(value int8) *Insert {
	i.internal.Assignments.Add("int8", value)
	return i
}

// SetInt16 sets value for column int16 in the INSERT statement
func (i *Insert) SetInt16(value int16) *Insert {
	i.internal.Assignments.Add("int16", value)
	return i
}

// SetInt32 sets value for column int32 in the INSERT statement
func (i *Insert) SetInt32(value int32) *Insert {
	i.internal.Assignments.Add("int32", value)
	return i
}

// SetInt64 sets value for column int64 in the INSERT statement
func (i *Insert) SetInt64(value int64) *Insert {
	i.internal.Assignments.Add("int64", value)
	return i
}

// SetUInt sets value for column uint in the INSERT statement
func (i *Insert) SetUInt(value uint) *Insert {
	i.internal.Assignments.Add("uint", value)
	return i
}

// SetUInt8 sets value for column uint8 in the INSERT statement
func (i *Insert) SetUInt8(value uint8) *Insert {
	i.internal.Assignments.Add("uint8", value)
	return i
}

// SetUInt16 sets value for column uint16 in the INSERT statement
func (i *Insert) SetUInt16(value uint16) *Insert {
	i.internal.Assignments.Add("uint16", value)
	return i
}

// SetUInt32 sets value for column uint32 in the INSERT statement
func (i *Insert) SetUInt32(value uint32) *Insert {
	i.internal.Assignments.Add("uint32", value)
	return i
}

// SetUInt64 sets value for column uint64 in the INSERT statement
func (i *Insert) SetUInt64(value uint64) *Insert {
	i.internal.Assignments.Add("uint64", value)
	return i
}

// SetTime sets value for column time in the INSERT statement
func (i *Insert) SetTime(value time.Time) *Insert {
	i.internal.Assignments.Add("time", value)
	return i
}

// SetVarCharString sets value for column varcharstring in the INSERT statement
func (i *Insert) SetVarCharString(value string) *Insert {
	i.internal.Assignments.Add("varcharstring", value)
	return i
}

// SetVarCharByte sets value for column varcharbyte in the INSERT statement
func (i *Insert) SetVarCharByte(value []byte) *Insert {
	i.internal.Assignments.Add("varcharbyte", value)
	return i
}

// SetString sets value for column string in the INSERT statement
func (i *Insert) SetString(value string) *Insert {
	i.internal.Assignments.Add("string", value)
	return i
}

// SetBytes sets value for column bytes in the INSERT statement
func (i *Insert) SetBytes(value []byte) *Insert {
	i.internal.Assignments.Add("bytes", value)
	return i
}

// SetBool sets value for column bool in the INSERT statement
func (i *Insert) SetBool(value bool) *Insert {
	i.internal.Assignments.Add("bool", value)
	return i
}

// SetPInt sets value for column pint in the INSERT statement
func (i *Insert) SetPInt(value *int) *Insert {
	i.internal.Assignments.Add("pint", value)
	return i
}

// SetPInt8 sets value for column pint8 in the INSERT statement
func (i *Insert) SetPInt8(value *int8) *Insert {
	i.internal.Assignments.Add("pint8", value)
	return i
}

// SetPInt16 sets value for column pint16 in the INSERT statement
func (i *Insert) SetPInt16(value *int16) *Insert {
	i.internal.Assignments.Add("pint16", value)
	return i
}

// SetPInt32 sets value for column pint32 in the INSERT statement
func (i *Insert) SetPInt32(value *int32) *Insert {
	i.internal.Assignments.Add("pint32", value)
	return i
}

// SetPInt64 sets value for column pint64 in the INSERT statement
func (i *Insert) SetPInt64(value *int64) *Insert {
	i.internal.Assignments.Add("pint64", value)
	return i
}

// SetPUInt sets value for column puint in the INSERT statement
func (i *Insert) SetPUInt(value *uint) *Insert {
	i.internal.Assignments.Add("puint", value)
	return i
}

// SetPUInt8 sets value for column puint8 in the INSERT statement
func (i *Insert) SetPUInt8(value *uint8) *Insert {
	i.internal.Assignments.Add("puint8", value)
	return i
}

// SetPUInt16 sets value for column puint16 in the INSERT statement
func (i *Insert) SetPUInt16(value *uint16) *Insert {
	i.internal.Assignments.Add("puint16", value)
	return i
}

// SetPUInt32 sets value for column puint32 in the INSERT statement
func (i *Insert) SetPUInt32(value *uint32) *Insert {
	i.internal.Assignments.Add("puint32", value)
	return i
}

// SetPUInt64 sets value for column puint64 in the INSERT statement
func (i *Insert) SetPUInt64(value *uint64) *Insert {
	i.internal.Assignments.Add("puint64", value)
	return i
}

// SetPTime sets value for column ptime in the INSERT statement
func (i *Insert) SetPTime(value *time.Time) *Insert {
	i.internal.Assignments.Add("ptime", value)
	return i
}

// SetPVarCharString sets value for column pvarcharstring in the INSERT statement
func (i *Insert) SetPVarCharString(value *string) *Insert {
	i.internal.Assignments.Add("pvarcharstring", value)
	return i
}

// SetPVarCharByte sets value for column pvarcharbyte in the INSERT statement
func (i *Insert) SetPVarCharByte(value *[]byte) *Insert {
	i.internal.Assignments.Add("pvarcharbyte", value)
	return i
}

// SetPString sets value for column pstring in the INSERT statement
func (i *Insert) SetPString(value *string) *Insert {
	i.internal.Assignments.Add("pstring", value)
	return i
}

// SetPBytes sets value for column pbytes in the INSERT statement
func (i *Insert) SetPBytes(value *[]byte) *Insert {
	i.internal.Assignments.Add("pbytes", value)
	return i
}

// SetPBool sets value for column pbool in the INSERT statement
func (i *Insert) SetPBool(value *bool) *Insert {
	i.internal.Assignments.Add("pbool", value)
	return i
}

// SetSelect sets value for column select in the INSERT statement
func (i *Insert) SetSelect(value int) *Insert {
	i.internal.Assignments.Add("select", value)
	return i
}
