package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func TrimFieldComment(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ";", ",")
	s = strings.ReplaceAll(s, "`", "'")
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

// PkgName returns the package name from a filepath
// with a package qualifier.
// ./model -> model
func GetPkgName(path string) string {
	pkgName := filepath.Base(path)
	if pkgName == "" || pkgName == "." {
		dir, _ := os.Getwd()
		workdir := strings.ReplaceAll(dir, "\\", "/")
		pkgName = filepath.Base(workdir)
	}
	return pkgName
}

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
