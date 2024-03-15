//go:generate go run internal/gen.go

package ens

import (
	"fmt"

	"github.com/things-go/ens/matcher"
	"gorm.io/plugin/soft_delete"
)

type FieldDescriptor struct {
	Name     string // field name
	Comment  string // field comment
	Nullable bool   // Nullable reports whether the column may be null.
	Column   ColumnDef
	// for go
	Type           *GoType  // go type information.
	Optional       bool     // nullable struct field.
	Tags           []string // Tags struct tag
	RapierDataType string   // rapier data type
}

func (field *FieldDescriptor) goType(typ any) {
	field.Type = NewGoType(field.Type.Type, typ)
}

func (field *FieldDescriptor) build(opt *Option) {
	field.RapierDataType = field.Type.Type.IntoRapierDataType()
	if field.Name == "deleted_at" && field.Type.IsInteger() {
		field.Optional = false
		field.goType(soft_delete.DeletedAt(0))
	}
	if opt == nil {
		return
	}

	if opt.EnableInt {
		switch field.Type.Type {
		case TypeInt8, TypeInt16, TypeInt32:
			field.goType(int(0))
			field.RapierDataType = TypeInt.IntoRapierDataType()
		case TypeUint8, TypeUint16, TypeUint32:
			field.goType(uint(0))
			field.RapierDataType = TypeUint.IntoRapierDataType()
		}
	}
	if opt.EnableIntegerInt {
		switch field.Type.Type {
		case TypeInt32:
			field.goType(int(0))
			field.RapierDataType = TypeInt.IntoRapierDataType()
		case TypeUint32:
			field.goType(uint(0))
			field.RapierDataType = TypeUint.IntoRapierDataType()
		}
	}
	if opt.EnableBoolInt && field.Type.IsBool() {
		field.goType(int(0))
		field.RapierDataType = TypeInt.IntoRapierDataType()
	}
	if field.Nullable && opt.DisableNullToPoint {
		gt, ok := sqlNullValueGoType[field.Type.Type]
		if ok {
			field.Type = gt.Clone()
			field.Optional = false
		}
	}

	for tag, kind := range opt.Tags {
		if tag == "json" {
			if vv := matcher.JsonTag(field.Comment); vv != "" {
				field.Tags = append(field.Tags, fmt.Sprintf(`%s:"%s"`, tag, vv))
				continue
			}
		}
		vv := TagName(kind, field.Name)
		if vv == "" {
			continue
		}
		if tag == "json" && matcher.HasAffixJSONTag(field.Comment) {
			field.Tags = append(field.Tags, fmt.Sprintf(`%s:"%s,omitempty,string"`, tag, vv))
		} else {
			field.Tags = append(field.Tags, fmt.Sprintf(`%s:"%s,omitempty"`, tag, vv))
		}
	}
}
