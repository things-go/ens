package ens

// IndexDescriptor
type IndexDescriptor struct {
	Name       string   // index name
	Fields     []string // field columns
	Definition string   // index sql definition
}

type indexBuilder struct {
	inner *IndexDescriptor
}

func Index(name string) *indexBuilder {
	return &indexBuilder{
		inner: &IndexDescriptor{
			Name:       name,
			Fields:     []string{},
			Definition: "",
		},
	}
}
func (b *indexBuilder) Fields(fields ...string) *indexBuilder {
	b.inner.Fields = append(b.inner.Fields, fields...)
	return b
}
func (b *indexBuilder) Definition(s string) *indexBuilder {
	b.inner.Definition = s
	return b
}
func (b *indexBuilder) Build() *IndexDescriptor {
	return b.inner
}
