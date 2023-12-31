package ens

import (
	"reflect"
)

var _ Fielder = (*boolBuilder)(nil)
var boolType = reflect.TypeOf(false)

// Bool returns a new Field with type bool.
func Bool(name string) *boolBuilder {
	return &boolBuilder{
		&FieldDescriptor{
			Name: name,
			Type: BoolType(),
		},
	}
}

// boolBuilder is the builder for boolean fields.
type boolBuilder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *boolBuilder) Comment(c string) *boolBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *boolBuilder) Nullable() *boolBuilder {
	b.inner.Nullable = true
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Bool("deleted").
//		GoType(sql.NullBool{})
func (b *boolBuilder) GoType(typ any) *boolBuilder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *boolBuilder) Optional() *boolBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of Tags to the field object.
//
//	field.String("dir").
//		Tags("yaml:"xxx"")
func (b *boolBuilder) Tags(tags ...string) *boolBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Field interface by returning its descriptor.
func (b *boolBuilder) Build(opt *Option) *FieldDescriptor {
	// b.inner.checkGoType(boolType)
	b.inner.build(opt)
	return b.inner
}
