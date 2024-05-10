package ens

import (
	"fmt"

	"github.com/things-go/ens/matcher"
	"github.com/things-go/ens/utils"
	"gorm.io/plugin/soft_delete"
)

type FieldDescriptor struct {
	Name     string // field name
	Comment  string // field comment
	Nullable bool   // Nullable reports whether the column may be null.
	Column   ColumnDef
	// for go
	Type     *GoType  // go type information.
	Optional bool     // nullable struct field.
	Tags     []string // Tags struct tag
}

func (field *FieldDescriptor) GoType(typ any) {
	field.Type = NewGoType(field.Type.Type, typ)
}

func (field *FieldDescriptor) build(opt *Option) {
	if field.Name == "deleted_at" && field.Type.IsInteger() {
		field.Optional = false
		field.GoType(soft_delete.DeletedAt(0))
	}
	if opt == nil {
		return
	}

	if opt.EnableInt {
		switch field.Type.Type {
		case TypeInt8, TypeInt16, TypeInt32:
			field.GoType(int(0))
		case TypeUint8, TypeUint16, TypeUint32:
			field.GoType(uint(0))
		}
	}
	if opt.EnableBoolInt && field.Type.IsBool() {
		field.GoType(int(0))
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
		vv := utils.StyleName(kind, field.Name)
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
