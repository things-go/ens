package ens

import (
	"reflect"
	"strings"
)

// PkgQualifier returns the package name from a Go
// identifier with a package qualifier.
// eg. time.Time -> time
func PkgQualifier(ident string) string {
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

// PkgName returns the package name from a Go
// path with a package qualifier.
// github.com/things/ens -> ens
func PkgName(path string) string {
	if path == "" {
		return ""
	}
	i := strings.LastIndexByte(path, '/')
	if i == -1 {
		return path
	}
	return path[i+1:]
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
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}
