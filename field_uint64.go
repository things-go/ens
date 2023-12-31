// Code generated by internal/integer.tpl, DO NOT EDIT.

package ens

import (
	"reflect"
)

var _ Fielder = (*uint64Builder)(nil)
var uint64Type = reflect.TypeOf(uint64(0))

// Uint64 returns a new Field with type uint64.
func Uint64(name string) *uint64Builder {
	return &uint64Builder{
		&FieldDescriptor{
			Name: name,
			Type: Uint64Type(),
		},
	}
}

// uint64Builder is the builder for uint64 field.
type uint64Builder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *uint64Builder) Comment(c string) *uint64Builder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *uint64Builder) Nullable() *uint64Builder {
	b.inner.Nullable = true
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Uint64("uint64").
//		GoType(pkg.Uint64(0))
func (b *uint64Builder) GoType(typ any) *uint64Builder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *uint64Builder) Optional() *uint64Builder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of tags to the field tag.
//
//	field.Uint64("uint64").
//		Tags("yaml:"xxx"")
func (b *uint64Builder) Tags(tags ...string) *uint64Builder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Fielder interface by returning its descriptor.
func (b *uint64Builder) Build(opt *Option) *FieldDescriptor {
	//	b.inner.checkGoType(uint64Type)
	b.inner.build(opt)
	return b.inner
}
