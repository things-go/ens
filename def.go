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

type Schemaer interface {
	Build(opt *Option) *Schema
}

type Indexer interface {
	Build() *IndexDescriptor
}

type Fielder interface {
	Build(opt *Option) *FieldDescriptor
}

type MixinEntity interface {
	Fields() []Fielder
	Indexes() []Indexer
	Metadata() (string, string)
	Table() TableDef
}
