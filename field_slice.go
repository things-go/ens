package ens

import (
	"reflect"
)

// Strings returns a new JSON Field with type []string.
func Strings(name string) *sliceBuilder[string] {
	return sb[string](name)
}

// Ints returns a new JSON Field with type []int.
func Ints(name string) *sliceBuilder[int] {
	return sb[int](name)
}

// Floats returns a new JSON Field with type []float.
func Floats(name string) *sliceBuilder[float64] {
	return sb[float64](name)
}

type sliceType interface {
	int | string | float64
}

// sliceBuilder is the builder for string slice fields.
type sliceBuilder[T sliceType] struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *sliceBuilder[T]) Comment(c string) *sliceBuilder[T] {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *sliceBuilder[T]) Nullable() *sliceBuilder[T] {
	b.inner.Nullable = true
	return b
}

// GoType overrides the default Go type with a custom one.
func (b *sliceBuilder[T]) GoType(typ any) *sliceBuilder[T] {
	return b
}

// Optional indicates that this field is optional on create.
// Unlike edges, fields are required by default.
func (b *sliceBuilder[T]) Optional() *sliceBuilder[T] {
	b.inner.Optional = true
	return b
}

// Tags adds a list of Tags to the field object.
//
//	field.String("dir").
//		Tags("yaml:"xxx"")
func (b *sliceBuilder[T]) Tags(tags ...string) *sliceBuilder[T] {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Field interface by returning its descriptor.
func (b *sliceBuilder[T]) Build(opt *Option) *FieldDescriptor {
	// b.inner.checkGoType(stringType)
	b.inner.build(opt)
	return b.inner
}

// sb is a generic helper method to share code between Strings, Ints and Floats builder.
func sb[T sliceType](name string) *sliceBuilder[T] {
	var typ []T
	b := &FieldDescriptor{
		Name: name,
		Type: &GoType{
			Type: TypeJSON,
		},
	}
	t := reflect.TypeOf(typ)
	if t == nil {
		return &sliceBuilder[T]{b}
	}
	b.Type.Ident = t.String()
	b.Type.PkgPath = t.PkgPath()
	// b.desc.goType(typ)
	// b.desc.checkGoType(t)
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Map:
		b.Type.Nullable = true
		b.Type.PkgPath = pkgPath(t)
	}
	return &sliceBuilder[T]{b}
}
