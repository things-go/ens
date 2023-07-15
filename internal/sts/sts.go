package sts

import (
	"fmt"
	"reflect"
)

func IsNill(v reflect.Value) bool {
	// if vv is not ptr, return false(v is not nil)
	// if vv is ptr, return v.IsNil()
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func IsStructType(v reflect.Type) bool {
	if v.Kind() == reflect.Struct {
		return true
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Kind() == reflect.Struct
}

func Parse(v any) error {
	st := reflect.TypeOf(v)
	sv := reflect.ValueOf(v)

	if IsNill(sv) {
		return nil
	}
	if !IsStructType(st) {
		return fmt.Errorf("%s is not a struct", st.String())
	}
	fmt.Println("---> ", st.Name())
	fmt.Println("---> ", st.String())
	fmt.Println("---> ", st.PkgPath())
	fmt.Println("---> ", st.PkgPath())
	return nil
}
