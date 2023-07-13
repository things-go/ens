package ens

type MixinEntity interface {
	Fields() []Fielder
	Indexes() []Indexer
	Metadata() EntityMetadata
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
