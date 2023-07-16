package sts

import (
	"fmt"
	"reflect"

	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

func IsNill(v reflect.Value) bool {
	// if vv is not ptr, return false(v is not nil)
	// if vv is ptr, return v.IsNil()
	return v.Kind() == reflect.Pointer && v.IsNil()
}

func indirect(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func Parse(v any) error {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return nil
	}
	vt := indirect(value.Type())
	for vt.Kind() == reflect.Pointer {
		vt = vt.Elem()
	}
	if vt.Kind() != reflect.Struct {
		return fmt.Errorf("%s is not a struct", vt.String())
	}
	fmt.Println("---> ", vt.Name())
	fmt.Println("---> ", vt.String())
	fmt.Println("---> ", vt.PkgPath())

	for i := 0; i < vt.NumField(); i++ {
		fv := vt.Field(i)
		if !fv.IsExported() {
			continue
		}
		fvt := fv.Type

		field := ens.Field(
			&ens.GoType{
				Type:         0,
				Ident:        fvt.String(),
				PkgPath:      fvt.PkgPath(),
				PkgQualifier: ens.PkgQualifier(fvt.String()),
				Nullable:     false,
			},
			utils.SnakeCase(fv.Name),
		)

		fmt.Println("-----")
		fmt.Println(fv.Name)
		fmt.Println(fvt.Kind())
		fmt.Printf("%#v\n", field.Build(nil))
	}

	return nil
}

func IntoType(kind reflect.Kind) ens.Type {
	t := ens.TypeInvalid
	switch kind {
	case reflect.Bool:
		return ens.TypeBool
	case reflect.Int:
		return ens.TypeInt
	case reflect.Int8:
		return ens.TypeInt8
	case reflect.Int16:
		return ens.TypeInt16
	case reflect.Int32:
		return ens.TypeInt32
	case reflect.Int64:
		return ens.TypeInt64
	case reflect.Uint:
		return ens.TypeUint
	case reflect.Uint8:
		return ens.TypeUint8
	case reflect.Uint16:
		return ens.TypeUint16
	case reflect.Uint32:
		return ens.TypeUint32
	case reflect.Uint64:
		return ens.TypeUint64
	case reflect.Float32:
		return ens.TypeFloat32
	case reflect.Float64:
		return ens.TypeFloat64
	case reflect.String:
		return ens.TypeFloat64
	case reflect.Struct:
	case reflect.Array:
	case reflect.Pointer:
	case reflect.Slice:
	}
	return t
}
