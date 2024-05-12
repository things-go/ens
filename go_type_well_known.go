package ens

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var sqlNullValueGoType = map[Type]*GoType{
	TypeBool:    SQLNullBoolType(),
	TypeUint8:   SQLNullByteType(),
	TypeInt16:   SQLNullInt16Type(),
	TypeInt32:   SQLNullInt32Type(),
	TypeInt64:   SQLNullInt64Type(),
	TypeFloat64: SQLNullFloat64Type(),
	TypeString:  SQLNullStringType(),
	TypeTime:    SQLNullTimeType(),
}

func BoolType() *GoType           { return NewGoType(TypeBool, false) }
func IntType() *GoType            { return NewGoType(TypeInt, int(0)) }
func Int8Type() *GoType           { return NewGoType(TypeInt8, int8(0)) }
func Int16Type() *GoType          { return NewGoType(TypeInt16, int16(0)) }
func Int32Type() *GoType          { return NewGoType(TypeInt32, int32(0)) }
func Int64Type() *GoType          { return NewGoType(TypeInt64, int64(0)) }
func UintType() *GoType           { return NewGoType(TypeUint, uint(0)) }
func Uint8Type() *GoType          { return NewGoType(TypeUint8, uint8(0)) }
func Uint16Type() *GoType         { return NewGoType(TypeUint16, uint16(0)) }
func Uint32Type() *GoType         { return NewGoType(TypeUint32, uint32(0)) }
func Uint64Type() *GoType         { return NewGoType(TypeUint64, uint64(0)) }
func Float32Type() *GoType        { return NewGoType(TypeFloat32, float32(0)) }
func Float64Type() *GoType        { return NewGoType(TypeFloat64, float64(0)) }
func StringType() *GoType         { return NewGoType(TypeString, "") }
func DecimalType() *GoType        { return NewGoType(TypeDecimal, "") }
func EnumType() *GoType           { return NewGoType(TypeEnum, "") }
func TimeType() *GoType           { return NewGoType(TypeTime, time.Time{}) }
func BytesType() *GoType          { return NewGoType(TypeBytes, []byte(nil)) }
func SQLNullBoolType() *GoType    { return NewGoType(TypeBool, sql.NullBool{}) }
func SQLNullByteType() *GoType    { return NewGoType(TypeUint8, sql.NullByte{}) }
func SQLNullFloat64Type() *GoType { return NewGoType(TypeFloat64, sql.NullFloat64{}) }
func SQLNullInt16Type() *GoType   { return NewGoType(TypeInt16, sql.NullInt16{}) }
func SQLNullInt32Type() *GoType   { return NewGoType(TypeInt32, sql.NullInt32{}) }
func SQLNullInt64Type() *GoType   { return NewGoType(TypeInt64, sql.NullInt64{}) }
func SQLNullStringType() *GoType  { return NewGoType(TypeString, sql.NullString{}) }
func SQLNullTimeType() *GoType    { return NewGoType(TypeTime, sql.NullTime{}) }
func JSONRawMessageType() *GoType { return NewGoType(TypeJSON, json.RawMessage{}) }
func SoftDeleteType() *GoType     { return NewGoType(TypeInt64, soft_delete.DeletedAt(0)) }
func GormDeletedAtType() *GoType  { return NewGoType(TypeInt64, gorm.DeletedAt{}) }
func DatatypesDateType() *GoType  { return NewGoType(TypeTime, datatypes.Date{}) }
func DatatypesJSONType() *GoType  { return NewGoType(TypeJSON, datatypes.JSON{}) }
