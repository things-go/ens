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
	Decimal
	Bytes
	Time
	endType

	JSON = Field
	UUID = Field
	Enum = String
)

type StructField struct {
	Type       Type   // 类型
	GoName     string // go名称, camel case
	Nullable   bool   // 是否Nullable
	ColumnName string // 列名, snake case
	Comment    string // 注释
}

type Struct struct {
	GoName    string         // go名称, camel case
	TableName string         // 表名, snake name
	Comment   string         // 注释
	Fields    []*StructField // 字段
}

type Schema struct {
	Name     string
	Entities []*Struct
}
