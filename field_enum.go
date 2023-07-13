package ens

func EnumType() *GoType {
	return NewGoType(TypeEnum, "")
}

// Enum returns a new Field with type enum.
func Enum(name string) *enumBuilder {
	return &enumBuilder{
		&FieldDescriptor{
			Name: name,
			Type: EnumType(),
		},
	}
}

// enumBuilder is the builder for enum fields.
type enumBuilder struct {
	inner *FieldDescriptor
}

// SchemaType sets the column type of the field.
func (b *enumBuilder) SchemaType(ct string) *enumBuilder {
	b.inner.SchemaType = ct
	return b
}

// Comment sets the comment of the field.
func (b *enumBuilder) Comment(c string) *enumBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *enumBuilder) Nullable() *enumBuilder {
	b.inner.Nullable = true
	return b
}

// Definition set the sql definition of the field.
func (b *enumBuilder) Definition(s string) *enumBuilder {
	b.inner.Definition = s
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Bool("deleted").
//		GoType(sql.NullBool{})
func (b *enumBuilder) GoType(typ any) *enumBuilder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *enumBuilder) Optional() *enumBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of Tags to the field object.
//
//	field.String("dir").
//		Tags("yaml:"xxx"")
func (b *enumBuilder) Tags(tags ...string) *enumBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Field interface by returning its descriptor.
func (b *enumBuilder) Build(opt *Option) *FieldDescriptor {
	// b.inner.checkGoType(stringType)
	b.inner.build(opt)
	return b.inner
}
