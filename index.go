package ens

import "github.com/things-go/ens/internal/sqlx"

// IndexDescriptor
type IndexDescriptor struct {
	Name   string   // index name
	Fields []string // field columns
	Index  IndexDef
}

type indexBuilder struct {
	inner *IndexDescriptor
}

// Index returns a new Index with the name.
func Index(name string) *indexBuilder {
	return &indexBuilder{
		inner: &IndexDescriptor{
			Name:   name,
			Fields: []string{},
		},
	}
}

// IndexFromDef returns a new Index with the IndexDef.
// auto set name, fields, index.
func IndexFromDef(e IndexDef) *indexBuilder {
	index := e.Index()
	return &indexBuilder{
		inner: &IndexDescriptor{
			Name:   index.Name,
			Fields: sqlx.IndexPartColumnNames(index.Parts),
			Index:  e,
		},
	}
}
func (b *indexBuilder) Fields(fields ...string) *indexBuilder {
	b.inner.Fields = append(b.inner.Fields, fields...)
	return b
}
func (b *indexBuilder) Index(i IndexDef) *indexBuilder {
	b.inner.Index = i
	return b
}
func (b *indexBuilder) Build() *IndexDescriptor {
	return b.inner
}
