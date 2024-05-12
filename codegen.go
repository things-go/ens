package ens

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/things-go/ens/matcher"
	"github.com/things-go/ens/utils"
	"golang.org/x/tools/imports"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var mustEscapeNames = []string{"TableName"}

type CodeGen struct {
	buf               bytes.Buffer
	Entities          []*EntityDescriptor
	ByName            string
	Version           string
	PackageName       string
	DisableDocComment bool
	Option
}

// Bytes returns the CodeBuf's buffer.
func (g *CodeGen) Bytes() []byte {
	return g.buf.Bytes()
}

// FormatSource return formats and adjusts imports contents of the CodeGen's buffer.
func (g *CodeGen) FormatSource() ([]byte, error) {
	data := g.buf.Bytes()
	if len(data) == 0 {
		return data, nil
	}
	// return format.Source(data)
	return imports.Process("", data, nil)
}

// Write appends the contents of p to the buffer,
func (g *CodeGen) Write(b []byte) (n int, err error) {
	return g.buf.Write(b)
}

// Print formats using the default formats for its operands and writes to the generated output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Print(a ...any) (n int, err error) {
	return fmt.Fprint(&g.buf, a...)
}

// Printf formats according to a format specifier for its operands and writes to the generated output.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(&g.buf, format, a...)
}

// Fprintln formats using the default formats to the generated output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Println(a ...any) (n int, err error) {
	return fmt.Fprintln(&g.buf, a...)
}

func (g *CodeGen) Gen() *CodeGen {
	if !g.DisableDocComment {
		g.Printf("// Code generated by %s. DO NOT EDIT.\n", g.ByName)
		g.Printf("// version: %s\n", g.Version)
		g.Println()
	}
	g.Printf("package %s\n", g.PackageName)
	g.Println()

	//* import
	imports := make(map[string]struct{})
	for _, st := range g.Entities {
		for _, field := range st.Fields {
			if field.Type.PkgPath != "" {
				imports[field.Type.PkgPath] = struct{}{}
			}
		}
	}
	if len(imports) > 0 {
		g.Println("import (")
		for k := range imports {
			g.Printf("\"%s\"\n", k)
		}
		g.Println(")")
	}

	//* struct
	for _, et := range g.Entities {
		structName := utils.CamelCase(et.Name)
		tableName := et.Name
		newEt := TransferEntityField(et, &g.Option)
		g.Printf("// %s %s\n", structName, strings.ReplaceAll(strings.TrimSpace(et.Comment), "\n", "\n// "))
		g.Printf("type %s struct {\n", structName)
		for _, field := range newEt.Fields {
			g.Println(g.genModelStructField(field))
		}
		g.Println("}")
		g.Println()
		g.Println("// TableName implement schema.Tabler interface")
		g.Printf("func (*%s) TableName() string {\n", structName)
		g.Printf("return \"%s\"\n", tableName)
		g.Println("}")
		g.Println()
	}
	return g
}

func (g *CodeGen) genModelStructField(field *FieldDescriptor) string {
	b := strings.Builder{}
	b.Grow(256)
	ident := field.Type.Ident
	if field.GoPointer && !field.Type.NoPointer {
		ident = "*" + field.Type.Ident
	}
	// field
	b.WriteString(field.GoName)
	b.WriteString(" ")
	b.WriteString(ident)
	if len(field.Tags) > 0 {
		b.WriteString(" `")
		b.WriteString(strings.Join(field.Tags, " "))
		b.WriteString("`")
	}
	if field.Comment != "" {
		b.WriteString(" // ")
		b.WriteString(field.Comment)
	}
	return b.String()
}

// BUG: 占CPU
func TransferEntityField(et *EntityDescriptor, opt *Option) *EntityDescriptor {
	newEt := *et
	newEt.Fields = make([]*FieldDescriptor, 0, len(et.Fields))

	if opt == nil {
		opt = defaultOption()
	}
	escapeNames := make(map[string]struct{})
	for _, v := range mustEscapeNames {
		escapeNames[v] = struct{}{}
	}
	for _, k := range opt.EscapeName {
		escapeNames[k] = struct{}{}
	}
	existFieldName := make(map[string]struct{}, len(et.Fields))
	for _, field := range et.Fields {
		existFieldName[field.GoName] = struct{}{}
	}
	for _, field := range et.Fields {
		newField := TransferField(field, opt)
		goName := newField.GoName
		for {
			_, ok := escapeNames[goName]
			if !ok { // need to escape
				break
			}
			goName = "X" + goName
			// 和当前字段存在的重复, 再追加一个
			_, ok = existFieldName[goName]
			if ok {
				goName = "X" + goName
			}
		}
		if newField.GoName != goName {
			newField.GoName = goName
			escapeNames[goName] = struct{}{} // 添加为必须转义
		}
		newEt.Fields = append(newEt.Fields, newField)
	}
	return &newEt
}

// 根规则转义一些数据
func TransferField(oldField *FieldDescriptor, opt *Option) *FieldDescriptor {
	newField := *oldField
	if newField.ColumnName == "deleted_at" {
		if newField.Type.IsInteger() {
			newField.GoPointer = false
			newField.GoType(soft_delete.DeletedAt(0))
		} else if newField.Type.IsInteger() {
			newField.GoPointer = false
			newField.GoType(gorm.DeletedAt{})
		}
	}
	if opt == nil {
		opt = defaultOption()
	}
	if opt.EnableInt {
		switch newField.Type.Type {
		case TypeInt8, TypeInt16, TypeInt32:
			newField.GoType(int(0))
		case TypeUint8, TypeUint16, TypeUint32:
			newField.GoType(uint(0))
		}
	}
	if opt.EnableBoolInt && newField.Type.IsBool() {
		newField.GoType(int(0))
	}
	if newField.Nullable && opt.DisableNullToPoint {
		gt, ok := sqlNullValueGoType[newField.Type.Type]
		if ok {
			newField.Type = gt.Clone()
			newField.GoPointer = false
		}
	}
	for tag, kind := range opt.Tags {
		if tag == "json" {
			if vv := matcher.JsonTag(newField.Comment); vv != "" {
				newField.Tags = append(newField.Tags, fmt.Sprintf(`%s:"%s"`, tag, vv))
				continue
			}
		}
		vv := utils.StyleName(kind, newField.ColumnName)
		if vv == "" {
			continue
		}
		if tag == "json" && matcher.HasAffixJSONTag(newField.Comment) {
			newField.Tags = append(newField.Tags, fmt.Sprintf(`%s:"%s,omitempty,string"`, tag, vv))
		} else {
			newField.Tags = append(newField.Tags, fmt.Sprintf(`%s:"%s,omitempty"`, tag, vv))
		}
	}
	return &newField
}
