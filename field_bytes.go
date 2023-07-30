package ens

import (
	"reflect"
)

var _ Fielder = (*bytesBuilder)(nil)
var bytesType = reflect.TypeOf([]byte(nil))

// Bytes returns a new Field with type bytes/buffer.
// In MySQL and SQLite, it is the "BLOB" type, and it does not support for Gremlin.
func Bytes(name string) *bytesBuilder {
	return &bytesBuilder{
		&FieldDescriptor{
			Name: name,
			Type: BytesType(),
		},
	}
}

// bytesBuilder is the builder for bytes fields.
type bytesBuilder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *bytesBuilder) Comment(c string) *bytesBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *bytesBuilder) Nullable() *bytesBuilder {
	b.inner.Nullable = true
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Bytes("ip").
//		GoType(net.IP("127.0.0.1"))
func (b *bytesBuilder) GoType(typ any) *bytesBuilder {
	// b.desc.goType(typ)
	return b
}

// Optional indicates that this field is optional on create.
// Unlike edges, fields are required by default.
func (b *bytesBuilder) Optional() *bytesBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of Tags to the field object.
//
//	field.String("dir").
//		Tags("yaml:"xxx"")
func (b *bytesBuilder) Tags(tags ...string) *bytesBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Field interface by returning its descriptor.
func (b *bytesBuilder) Build(opt *Option) *FieldDescriptor {
	// b.inner.checkGoType(bytesType)
	b.inner.build(opt)
	return b.inner
}
