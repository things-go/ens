package ens

import (
	"reflect"
	"strings"
)

// pkgName returns the package name from a Go
// identifier with a package qualifier.
func pkgName(ident string) string {
	i := strings.LastIndexByte(ident, '.')
	if i == -1 {
		return ""
	}
	s := ident[:i]
	if i := strings.LastIndexAny(s, "]*"); i != -1 {
		s = s[i+1:]
	}
	return s
}

func pkgPath(t reflect.Type) string {
	pkg := t.PkgPath()
	if pkg != "" {
		return pkg
	}
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Map:
		return pkgPath(t.Elem())
	}
	return pkg
}

func indirect(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
