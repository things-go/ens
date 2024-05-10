package ens

import (
	"database/sql/driver"

	"github.com/things-go/ens/internal/insql"
)

func Bool(name string) *FieldBuilder                    { return Field(BoolType(), name) }
func Int(name string) *FieldBuilder                     { return Field(IntType(), name) }
func Int8(name string) *FieldBuilder                    { return Field(Int8Type(), name) }
func Int16(name string) *FieldBuilder                   { return Field(Int16Type(), name) }
func Int32(name string) *FieldBuilder                   { return Field(Int32Type(), name) }
func Int64(name string) *FieldBuilder                   { return Field(Int64Type(), name) }
func Uint(name string) *FieldBuilder                    { return Field(UintType(), name) }
func Uint8(name string) *FieldBuilder                   { return Field(Uint8Type(), name) }
func Uint16(name string) *FieldBuilder                  { return Field(Uint16Type(), name) }
func Uint32(name string) *FieldBuilder                  { return Field(Uint32Type(), name) }
func Uint64(name string) *FieldBuilder                  { return Field(Uint64Type(), name) }
func Float32(name string) *FieldBuilder                 { return Field(Float32Type(), name) }
func Float64(name string) *FieldBuilder                 { return Field(Float64Type(), name) }
func String(name string) *FieldBuilder                  { return Field(StringType(), name) }
func Bytes(name string) *FieldBuilder                   { return Field(BytesType(), name) }
func JSON(name string) *FieldBuilder                    { return Field(JSONRawMessageType(), name) }
func Enum(name string) *FieldBuilder                    { return Field(EnumType(), name) }
func Time(name string) *FieldBuilder                    { return Field(TimeType(), name) }
func UUID(name string, typ driver.Valuer) *FieldBuilder { return Field(NewGoType(TypeUUID, typ), name) }

var _ Fielder = (*FieldBuilder)(nil)

// FieldBuilder is the builder for field.
type FieldBuilder struct {
	inner *FieldDescriptor
}

// Field returns a new Field with the type.
func Field(t *GoType, name string) *FieldBuilder {
	return &FieldBuilder{
		inner: &FieldDescriptor{
			Name: name,
			Type: t,
		},
	}
}

// FieldFromDef returns a new Field with the type and ColumnDef.
// auto set name, comment, nullable, column and optional.
func FieldFromDef(t *GoType, def ColumnDef) *FieldBuilder {
	col := def.Column()
	return &FieldBuilder{
		inner: &FieldDescriptor{
			Name:     col.Name,
			Comment:  insql.MustComment(col.Attrs),
			Nullable: col.Type.Null,
			Column:   def,
			Type:     t,
			Optional: col.Type.Null,
		},
	}
}

// Column the column expression of the field.
func (b *FieldBuilder) Column(e ColumnDef) *FieldBuilder {
	b.inner.Column = e
	return b
}

// Comment sets the comment of the field.
func (b *FieldBuilder) Comment(c string) *FieldBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *FieldBuilder) Nullable() *FieldBuilder {
	b.inner.Nullable = true
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Bool("deleted").
//		GoType(sql.NullBool{})
//	field.Bytes("ip").
//		GoType(net.IP("127.0.0.1"))
//	field.String("dir").
//		GoType(http.Dir("dir"))
func (b *FieldBuilder) GoType(typ any) *FieldBuilder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *FieldBuilder) Optional() *FieldBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of tags to the field tag.
func (b *FieldBuilder) Tags(tags ...string) *FieldBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Fielder interface by returning its descriptor.
func (b *FieldBuilder) Build(opt *Option) *FieldDescriptor {
	b.inner.build(opt)
	return b.inner
}
