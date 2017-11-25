// Autogenerated by github.com/posener/orm; DO NOT EDIT
package allsqlite3

type columns struct {
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

	count bool // used for sql COUNT(*) column
}

// String is the SQL representation of columns
func (c *columns) Columns() []string {
	var cols []string
	if c.SelectAuto {
		cols = append(cols, "auto")
	}
	if c.SelectNotNil {
		cols = append(cols, "notnil")
	}
	if c.SelectInt {
		cols = append(cols, "int")
	}
	if c.SelectInt8 {
		cols = append(cols, "int8")
	}
	if c.SelectInt16 {
		cols = append(cols, "int16")
	}
	if c.SelectInt32 {
		cols = append(cols, "int32")
	}
	if c.SelectInt64 {
		cols = append(cols, "int64")
	}
	if c.SelectUInt {
		cols = append(cols, "uint")
	}
	if c.SelectUInt8 {
		cols = append(cols, "uint8")
	}
	if c.SelectUInt16 {
		cols = append(cols, "uint16")
	}
	if c.SelectUInt32 {
		cols = append(cols, "uint32")
	}
	if c.SelectUInt64 {
		cols = append(cols, "uint64")
	}
	if c.SelectTime {
		cols = append(cols, "time")
	}
	if c.SelectVarCharString {
		cols = append(cols, "varcharstring")
	}
	if c.SelectVarCharByte {
		cols = append(cols, "varcharbyte")
	}
	if c.SelectString {
		cols = append(cols, "string")
	}
	if c.SelectBytes {
		cols = append(cols, "bytes")
	}
	if c.SelectBool {
		cols = append(cols, "bool")
	}
	if c.SelectPInt {
		cols = append(cols, "pint")
	}
	if c.SelectPInt8 {
		cols = append(cols, "pint8")
	}
	if c.SelectPInt16 {
		cols = append(cols, "pint16")
	}
	if c.SelectPInt32 {
		cols = append(cols, "pint32")
	}
	if c.SelectPInt64 {
		cols = append(cols, "pint64")
	}
	if c.SelectPUInt {
		cols = append(cols, "puint")
	}
	if c.SelectPUInt8 {
		cols = append(cols, "puint8")
	}
	if c.SelectPUInt16 {
		cols = append(cols, "puint16")
	}
	if c.SelectPUInt32 {
		cols = append(cols, "puint32")
	}
	if c.SelectPUInt64 {
		cols = append(cols, "puint64")
	}
	if c.SelectPTime {
		cols = append(cols, "ptime")
	}
	if c.SelectPVarCharString {
		cols = append(cols, "pvarcharstring")
	}
	if c.SelectPVarCharByte {
		cols = append(cols, "pvarcharbyte")
	}
	if c.SelectPString {
		cols = append(cols, "pstring")
	}
	if c.SelectPBytes {
		cols = append(cols, "pbytes")
	}
	if c.SelectPBool {
		cols = append(cols, "pbool")
	}
	if c.SelectSelect {
		cols = append(cols, "select")
	}

	return cols
}

func (c *columns) Count() bool {
	return c.count
}

// selectAll returns true if no column was specifically selected
func (c *columns) selectAll() bool {
	return !c.SelectAuto && !c.SelectNotNil && !c.SelectInt && !c.SelectInt8 && !c.SelectInt16 && !c.SelectInt32 && !c.SelectInt64 && !c.SelectUInt && !c.SelectUInt8 && !c.SelectUInt16 && !c.SelectUInt32 && !c.SelectUInt64 && !c.SelectTime && !c.SelectVarCharString && !c.SelectVarCharByte && !c.SelectString && !c.SelectBytes && !c.SelectBool && !c.SelectPInt && !c.SelectPInt8 && !c.SelectPInt16 && !c.SelectPInt32 && !c.SelectPInt64 && !c.SelectPUInt && !c.SelectPUInt8 && !c.SelectPUInt16 && !c.SelectPUInt32 && !c.SelectPUInt64 && !c.SelectPTime && !c.SelectPVarCharString && !c.SelectPVarCharByte && !c.SelectPString && !c.SelectPBytes && !c.SelectPBool && !c.SelectSelect && !c.count
}
