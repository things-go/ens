// Code generated by internal/integer.tpl, DO NOT EDIT.

package ens

import (
	"reflect"
)

var _ Fielder = (*int8Builder)(nil)
var int8Type = reflect.TypeOf(int8(0))

// Int8 returns a new Field with type int8.
func Int8(name string) *int8Builder {
	return &int8Builder{
		&FieldDescriptor{
			Name: name,
			Type: Int8Type(),
		},
	}
}

// int8Builder is the builder for int8 field.
type int8Builder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *int8Builder) Comment(c string) *int8Builder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *int8Builder) Nullable() *int8Builder {
	b.inner.Nullable = true
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Int8("int8").
//		GoType(pkg.Int8(0))
func (b *int8Builder) GoType(typ any) *int8Builder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *int8Builder) Optional() *int8Builder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of tags to the field tag.
//
//	field.Int8("int8").
//		Tags("yaml:"xxx"")
func (b *int8Builder) Tags(tags ...string) *int8Builder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Fielder interface by returning its descriptor.
func (b *int8Builder) Build(opt *Option) *FieldDescriptor {
	//	b.inner.checkGoType(int8Type)
	b.inner.build(opt)
	return b.inner
}
