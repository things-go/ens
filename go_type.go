package ens

import (
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

var typeNames = [...]string{
	TypeInvalid: "invalid",
	TypeBool:    "bool",
	TypeInt8:    "int8",
	TypeInt16:   "int16",
	TypeInt32:   "int32",
	TypeInt64:   "int64",
	TypeInt:     "int",
	TypeUint8:   "uint8",
	TypeUint16:  "uint16",
	TypeUint32:  "uint32",
	TypeUint64:  "uint64",
	TypeUint:    "uint",
	TypeFloat32: "float32",
	TypeFloat64: "float64",
	TypeDecimal: "string",
	TypeString:  "string",
	TypeEnum:    "string",
	TypeBytes:   "[]byte",
	TypeTime:    "time.Time",
	TypeJSON:    "json.RawMessage",
	TypeUUID:    "[16]byte",
	TypeOther:   "other",
}

// A Type represents a field type.
type Type uint8

// List of field types.
const (
	TypeInvalid Type = iota
	TypeBool
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeInt
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64
	TypeUint
	TypeFloat32
	TypeFloat64
	TypeDecimal
	TypeString
	TypeEnum
	TypeBytes
	TypeTime
	TypeJSON
	TypeUUID
	TypeOther
	endTypes
)

// String returns the string representation of a type.
func (t Type) String() string {
	if t < endTypes {
		return typeNames[t]
	}
	return typeNames[TypeInvalid]
}

// IsNumeric reports if the given type is a numeric type.
func (t Type) IsNumeric() bool {
	return t >= TypeInt8 && t <= TypeFloat64
}

// IsFloat reports if the given type is a float type.
func (t Type) IsFloat() bool {
	return t == TypeFloat32 || t == TypeFloat64
}

// IsInteger reports if the given type is an integral type.
func (t Type) IsInteger() bool {
	return t.IsNumeric() && !t.IsFloat()
}

// IsBool reports if the given type is an bool type.
func (t Type) IsBool() bool {
	return t == TypeBool
}

// IsTime reports if the given type is an time.Time type.
func (t Type) IsTime() bool {
	return t == TypeTime
}

// IsValid reports if the given type if known type.
func (t Type) IsValid() bool {
	return t > TypeInvalid && t < endTypes
}

func (t Type) IntoProtoDataType() string {
	dataType := ""
	switch t {
	case TypeBool:
		dataType = "bool"
	case TypeInt8, TypeInt16, TypeInt32, TypeInt:
		dataType = "int32"
	case TypeInt64:
		dataType = "int64"
	case TypeUint8, TypeUint16, TypeUint32, TypeUint:
		dataType = "uint32"
	case TypeUint64:
		dataType = "uint64"
	case TypeFloat32:
		dataType = "float"
	case TypeFloat64:
		dataType = "double"
	case TypeDecimal, TypeString, TypeEnum, TypeJSON, TypeUUID, TypeOther:
		dataType = "string"
	case TypeBytes:
		dataType = "bytes"
	case TypeTime:
		dataType = "google.protobuf.Timestamp"
	default:
		dataType = "string"
	}
	return dataType
}

func (t Type) IntoAssistDataType() string {
	switch t {
	case TypeBool,
		TypeInt8, TypeInt16, TypeInt32, TypeInt64, TypeInt,
		TypeUint8, TypeUint16, TypeUint32, TypeUint64, TypeUint,
		TypeFloat32, TypeFloat64, TypeString, TypeEnum:
		return strings.Title(t.String())
	case TypeDecimal:
		return "Decimal"
	case TypeBytes:
		return "Bytes"
	case TypeTime:
		return "Time"
	case TypeJSON, TypeUUID, TypeOther:
		fallthrough
	default:
		return "Field"
	}
}

type GoType struct {
	Type         Type   // Type enum.
	Ident        string // Type identifier,  e.g. int, time.Time, sql.NullInt64.
	PkgPath      string // import path. e.g. "", time, database/sql.
	PkgQualifier string // a package qualifier. e.g. "", time, sql.
	Nullable     bool   // pointers or slices, means not need point.
}

// String returns the string representation of a type.
func (t GoType) String() string {
	switch {
	case t.Ident != "":
		return t.Ident
	case t.Type < endTypes:
		return typeNames[t.Type]
	default:
		return typeNames[TypeInvalid]
	}
}

// IsNumeric reports if the given type is a numeric type.
func (t GoType) IsNumeric() bool {
	return t.Type.IsNumeric()
}

// IsFloat reports if the given type is a float type.
func (t GoType) IsFloat() bool {
	return t.Type.IsFloat()
}

// IsInteger reports if the given type is an integral type.
func (t GoType) IsInteger() bool {
	return t.Type.IsInteger()
}

// IsBool reports if the given type is an bool type.
func (t GoType) IsBool() bool {
	return t.Type.IsBool()
}

// IsTime reports if the given type is an time.Time type.
func (t GoType) IsTime() bool {
	return t.Type.IsTime()
}

// IsValid reports if the given type if known type.
func (t GoType) IsValid() bool {
	return t.Type.IsValid()
}

// Comparable reports whether values of this type are comparable.
func (t GoType) Comparable() bool {
	switch t.Type {
	case TypeBool, TypeTime, TypeUUID, TypeEnum, TypeString:
		return true
	case TypeOther:
		// Always accept custom types as comparable on the database side.
		// In the future, we should consider adding an interface to let
		// custom types tell if they are comparable or not (see #1304).
		return true
	default:
		return t.Type.IsNumeric()
	}
}

func NewGoType(t Type, v any) *GoType {
	tt := reflect.TypeOf(v)
	tv := indirect(tt)
	return &GoType{
		Type:         t,
		Ident:        tt.String(),
		PkgPath:      tv.PkgPath(),
		PkgQualifier: PkgQualifier(tv.String()),
		Nullable:     slices.Contains([]reflect.Kind{reflect.Slice, reflect.Ptr, reflect.Map}, tt.Kind()),
	}
}
