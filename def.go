package ens

import "ariga.io/atlas/sql/schema"

type TableDef interface {
	Table() *schema.Table
	Definition() string
	PrimaryKey() IndexDef
}

type ColumnDef interface {
	Column() *schema.Column
	Definition() string
	GormTag(*schema.Table) string
}

type IndexDef interface {
	Index() *schema.Index
	Definition() string
}

type ForeignKeyDef interface {
	ForeignKey() *schema.ForeignKey
	Definition() string
}

type Indexer interface {
	Build() *IndexDescriptor
}

type Fielder interface {
	Build(*Option) *FieldDescriptor
}

type ForeignKeyer interface {
	Build() *ForeignKeyDescriptor
}

type MixinEntity interface {
	Metadata() (string, string)
	Table() TableDef
	Fields() []Fielder
	Indexes() []Indexer
	ForeignKeys() []ForeignKeyer

	Build(*Option) *EntityDescriptor
}
