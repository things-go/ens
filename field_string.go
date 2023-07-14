package ens

import (
	"reflect"
)

var _ Fielder = (*stringBuilder)(nil)
var stringType = reflect.TypeOf("")

func StringType() *GoType {
	return NewGoType(TypeString, "")
}

func DecimalType() *GoType {
	return NewGoType(TypeDecimal, "")
}

// String returns a new Field with type string.
// limitation on the size 255.
func String(name string) *stringBuilder {
	return &stringBuilder{
		&FieldDescriptor{
			Name: name,
			Type: StringType(),
		},
	}
}

// Decimal returns a new Field with type decimal.
// limitation on the size 255.
func Decimal(name string) *stringBuilder {
	return &stringBuilder{
		&FieldDescriptor{
			Name: name,
			Type: DecimalType(),
		},
	}
}

// stringBuilder is the builder for string fields.
type stringBuilder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *stringBuilder) Comment(c string) *stringBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *stringBuilder) Nullable() *stringBuilder {
	b.inner.Nullable = true
	return b
}

// Definition set the sql definition of the field.
func (b *stringBuilder) Definition(s string) *stringBuilder {
	b.inner.Definition = s
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.String("dir").
//		GoType(http.Dir("dir"))
func (b *stringBuilder) GoType(typ any) *stringBuilder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *stringBuilder) Optional() *stringBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of Tags to the field object.
//
//	field.String("dir").
//		Tags("yaml:"xxx"")
func (b *stringBuilder) Tags(tags ...string) *stringBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Field interface by returning its descriptor.
func (b *stringBuilder) Build(opt *Option) *FieldDescriptor {
	// b.inner.checkGoType(stringType)
	b.inner.build(opt)
	return b.inner
}
