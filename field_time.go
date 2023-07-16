package ens

import (
	"reflect"
	"time"
)

var _ Fielder = (*timeBuilder)(nil)
var timeType = reflect.TypeOf(time.Time{})

// Time returns a new Field with type timestamp.
func Time(name string) *timeBuilder {
	return &timeBuilder{
		&FieldDescriptor{
			Name: name,
			Type: TimeType(),
		},
	}
}

// timeBuilder is the builder for time fields.
type timeBuilder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *timeBuilder) Comment(c string) *timeBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *timeBuilder) Nullable() *timeBuilder {
	b.inner.Nullable = true
	return b
}

// Definition set the sql definition of the field.
func (b *timeBuilder) Definition(s string) *timeBuilder {
	b.inner.Definition = s
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.String("dir").
//		GoType(http.Dir("dir"))
func (b *timeBuilder) GoType(typ any) *timeBuilder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *timeBuilder) Optional() *timeBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of Tags to the field object.
//
//	field.String("dir").
//		Tags("yaml:"xxx"")
func (b *timeBuilder) Tags(tags ...string) *timeBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Field interface by returning its descriptor.
func (b *timeBuilder) Build(opt *Option) *FieldDescriptor {
	// b.inner.checkGoType(timeType)
	b.inner.build(opt)
	return b.inner
}
