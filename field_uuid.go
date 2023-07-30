package ens

import (
	"database/sql/driver"
)

var _ Fielder = (*uuidBuilder)(nil)

// UUID returns a new Field with type UUID.
func UUID(name string, typ driver.Valuer) *uuidBuilder {
	return &uuidBuilder{
		&FieldDescriptor{
			Name: name,
			Type: NewGoType(TypeUUID, typ),
		}}
}

// uuidBuilder is the builder for uuid fields.
type uuidBuilder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *uuidBuilder) Comment(c string) *uuidBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *uuidBuilder) Nullable() *uuidBuilder {
	b.inner.Nullable = true
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.String("dir").
//		GoType(http.Dir("dir"))
func (b *uuidBuilder) GoType(typ any) *uuidBuilder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *uuidBuilder) Optional() *uuidBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of Tags to the field object.
//
//	field.String("dir").
//		Tags("yaml:"xxx"")
func (b *uuidBuilder) Tags(tags ...string) *uuidBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Field interface by returning its descriptor.
func (b *uuidBuilder) Build(opt *Option) *FieldDescriptor {
	// b.inner.checkGoType(stringType)
	b.inner.build(opt)
	return b.inner
}
