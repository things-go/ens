//go:generate go run internal/gen.go

package ens

import (
	"fmt"

	"github.com/things-go/ens/matcher"
)

type FieldDescriptor struct {
	Name       string // field name
	Comment    string // comment
	Nullable   bool   // Nullable reports whether the column may be null.
	Definition string // field sql definition
	// for go
	Type           *GoType  //  go type information.
	Optional       bool     // nullable struct field.
	Tags           []string // Tags struct tag
	AssistDataType string   // assist data type
}

func (f *FieldDescriptor) goType(typ any) {
	f.Type = NewGoType(f.Type.Type, typ)
}

func (field *FieldDescriptor) build(opt *Option) {
	field.AssistDataType = field.Type.Type.IntoAssistDataType()
	if field.Name == "deleted_at" && field.Type.IsInteger() {
		field.Optional = false
		field.goType(zeroSoftDelete)
	}
	if opt == nil {
		return
	}

	if opt.EnableInt {
		switch field.Type.Type {
		case TypeInt8, TypeInt16, TypeInt32:
			field.goType(zeroInt)
			field.AssistDataType = TypeInt.IntoAssistDataType()
		case TypeUint8, TypeUint16, TypeUint32:
			field.goType(zeroUint)
			field.AssistDataType = TypeUint.IntoAssistDataType()
		}
	}
	if opt.EnableIntegerInt {
		switch field.Type.Type {
		case TypeInt32:
			field.goType(zeroInt)
			field.AssistDataType = TypeInt.IntoAssistDataType()
		case TypeUint32:
			field.goType(zeroUint)
			field.AssistDataType = TypeUint.IntoAssistDataType()
		}
	}
	if opt.EnableBoolInt && field.Type.IsBool() {
		field.goType(zeroInt)
		field.AssistDataType = TypeInt.IntoAssistDataType()
	}
	if field.Nullable && opt.DisableNullToPoint {
		zeroValue, ok := zeroSqlNullValue[field.Type.Type]
		if ok {
			field.goType(zeroValue)
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

// Field returns a new Field with the type.
func Field(t *GoType, name string) *fieldBuilder {
	return &fieldBuilder{
		inner: &FieldDescriptor{
			Name: name,
			Type: t,
		},
	}
}

var _ Fielder = (*fieldBuilder)(nil)

// fieldBuilder is the builder for field.
type fieldBuilder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *fieldBuilder) Comment(c string) *fieldBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *fieldBuilder) Nullable() *fieldBuilder {
	b.inner.Nullable = true
	return b
}

// Definition set the sql definition of the field.
func (b *fieldBuilder) Definition(s string) *fieldBuilder {
	b.inner.Definition = s
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Bool("deleted").
//		GoType(sql.NullBool{})
func (b *fieldBuilder) GoType(typ any) *fieldBuilder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *fieldBuilder) Optional() *fieldBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of tags to the field tag.
func (b *fieldBuilder) Tags(tags ...string) *fieldBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Fielder interface by returning its descriptor.
func (b *fieldBuilder) Build(opt *Option) *FieldDescriptor {
	b.inner.build(opt)
	return b.inner
}
