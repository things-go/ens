package ens

import (
	"database/sql"

	"gorm.io/datatypes"
	"gorm.io/plugin/soft_delete"
)

var (
	zeroSoftDelete = soft_delete.DeletedAt(0)
	zeroUint       = uint(0)
	zeroInt        = int(0)
	zeroDate       = datatypes.Date{}
	zeroJson       = datatypes.JSON{}
)

var zeroSqlNullValue = map[Type]any{
	TypeBool:    sql.NullBool{},
	TypeUint:    sql.NullByte{},
	TypeFloat64: sql.NullFloat64{},
	TypeInt16:   sql.NullInt16{},
	TypeInt32:   sql.NullInt32{},
	TypeInt64:   sql.NullInt64{},
	TypeString:  sql.NullString{},
	TypeTime:    sql.NullTime{},
}
