package ens

type FieldDescriptor struct {
	ColumnName string // 列名, snake case
	Comment    string // 注释
	Nullable   bool   // Nullable reports whether the column may be null.
	Column     ColumnDef
	// for go
	Type      *GoType  // go type information.
	GoName    string   // Go name
	GoPointer bool     // go field is pointer.
	Tags      []string // Tags struct tag
}

func (field *FieldDescriptor) GoType(typ any) {
	field.Type = NewGoType(field.Type.Type, typ)
}
