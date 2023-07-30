package ens

var _ Fielder = (*jsonBuilder)(nil)

// JSON returns a new Field with type json that is serialized to the given object.
func JSON(name string) *jsonBuilder {
	return &jsonBuilder{
		&FieldDescriptor{
			Name: name,
			Type: JSONRawMessageType(),
		},
	}
}

// jsonBuilder is the builder for json fields.
type jsonBuilder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *jsonBuilder) Comment(c string) *jsonBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *jsonBuilder) Nullable() *jsonBuilder {
	b.inner.Nullable = true
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.String("dir").
//		GoType(http.Dir("dir"))
func (b *jsonBuilder) GoType(typ any) *jsonBuilder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *jsonBuilder) Optional() *jsonBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of Tags to the field object.
//
//	field.String("dir").
//		Tags("yaml:"xxx"")
func (b *jsonBuilder) Tags(tags ...string) *jsonBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Field interface by returning its descriptor.
func (b *jsonBuilder) Build(opt *Option) *FieldDescriptor {
	// b.inner.checkGoType(stringType)
	b.inner.build(opt)
	return b.inner
}
