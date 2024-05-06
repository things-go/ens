//go:generate stringer -type Type
package rapier

type Type int

const (
	Field Type = iota
	Bool
	Int8
	Int16
	Int32
	Int64
	Int
	Uint8
	Uint16
	Uint32
	Uint64
	Uint
	Float32
	Float64
	String
	Enum
	Decimal
	Bytes
	Time
	endType

	JSON = Field
	UUID = Field
)

type StructField struct {
	Type       Type
	GoName     string
	Nullable   bool
	ColumnName string
}

type Struct struct {
	GoName    string
	TableName string
	Fields    []*StructField
}
