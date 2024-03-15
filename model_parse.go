package ens

import (
	"fmt"
	"reflect"

	"github.com/things-go/ens/utils"
)

func ParseModel(v any) (MixinEntity, error) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return nil, nil
	}
	vt := indirect(value.Type())
	for vt.Kind() == reflect.Pointer {
		vt = vt.Elem()
	}
	if vt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s is not a struct", vt.String())
	}
	entityBuilder := &EntityBuilder{}
	fields := structToFielders(vt)
	return entityBuilder.
		SetMetadata(utils.SnakeCase(vt.Name()), "").
		SetFields(fields...), nil
}

func structToFielders(vt reflect.Type) []Fielder {
	fields := make([]Fielder, 0, vt.NumField())
	for i := 0; i < vt.NumField(); i++ {
		fv := vt.Field(i)
		if !fv.IsExported() {
			continue
		}
		// an embedded field
		if fv.Anonymous {
			fvt := fv.Type
			for fvt.Kind() == reflect.Ptr {
				fvt = fv.Type.Elem()
			}
			if fvt.Kind() != reflect.Struct {
				continue
			}
			fields = append(fields, structToFielders(fvt)...)
		} else {
			fields = append(fields, structFieldToFielder(fv))
		}
	}
	return fields
}

func structFieldToFielder(fv reflect.StructField) Fielder {
	fvt := fv.Type
	nullable := false
	for fvt.Kind() == reflect.Ptr {
		fvt = fv.Type.Elem()
		nullable = true
	}

	fieldName := utils.SnakeCase(fv.Name)
	// setting := schema.ParseTagSetting(fv.Tag.Get("gorm"), ";")
	// if v, ok := setting["COLUMN"]; ok && v != "" {
	// 	fieldName = v
	// }
	ident := fvt.String()
	return Field(
		&GoType{
			Type:         intoGoTypeType(fvt.Kind(), ident),
			Ident:        fvt.String(),
			PkgPath:      fvt.PkgPath(),
			PkgQualifier: PkgQualifier(fvt.String()),
			Nullable:     nullable,
		},
		fieldName,
	)
}

func intoGoTypeType(kind reflect.Kind, ident string) Type {
	switch kind {
	case reflect.Bool:
		return TypeBool
	case reflect.Int:
		return TypeInt
	case reflect.Int8:
		return TypeInt8
	case reflect.Int16:
		return TypeInt16
	case reflect.Int32:
		return TypeInt32
	case reflect.Int64:
		return TypeInt64
	case reflect.Uint:
		return TypeUint
	case reflect.Uint8:
		return TypeUint8
	case reflect.Uint16:
		return TypeUint16
	case reflect.Uint32:
		return TypeUint32
	case reflect.Uint64:
		return TypeUint64
	case reflect.Float32:
		return TypeFloat32
	case reflect.Float64:
		return TypeFloat64
	case reflect.String:
		return TypeString
	case reflect.Struct:
		switch ident {
		case "time.Time", "sql.NullTime":
			return TypeTime
		case "sql.NullBool":
			return TypeBool
		case "sql.NullByte":
			return TypeBytes
		case "sql.NullString":
			return TypeString
		case "sql.NullFloat64":
			return TypeFloat64
		case "sql.NullInt16":
			return TypeInt16
		case "sql.NullInt32":
			return TypeInt32
		case "sql.NullInt64":
			return TypeInt64
		default:
			return TypeOther
		}
	// case reflect.Array:
	// case reflect.Slice:
	default:
		return TypeOther
	}
}
