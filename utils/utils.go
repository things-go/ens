package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func TrimComment(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ";", ",")
	s = strings.ReplaceAll(s, "`", "'")
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

func GetPkgName(path string) string {
	pkgName := filepath.Base(path)
	if pkgName == "" || pkgName == "." {
		dir, _ := os.Getwd()
		workdir := strings.ReplaceAll(dir, "\\", "/")
		pkgName = filepath.Base(workdir)
	}
	return pkgName
}
